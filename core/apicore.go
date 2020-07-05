package core

import (
	"context"
	"errors"
	"fmt"
	l "log"
	"net/http"
	"time"

	"github.com/felixge/httpsnoop"
	"github.com/google/uuid"
	muxcontext "github.com/gorilla/context"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/pontiyaraja/AB/ablog"
	"github.com/unrolled/render"
)

const (
	//TransactionID for generating a random number for a request
	TransactionID = "transctionId"

	//RequestID of the parent request
	RequestID = "requestId"

	//Email address of the client
	UserID = "userID"
)

func ValidateToken(apikey, apisecret, requestID string) (context.Context, error) {
	if len(apikey) == 0 || len(apisecret) == 0 {
		return nil, errors.New("api key or api key secret is missing")
	}
	//validate credentials
	return nil, nil
}

// commonAuthRules applies common authentication rules across middleware
func commonAuthRules(apikey, apisecret, requestID string) (context.Context, error) {
	if len(apikey) == 0 || apikey == "undefined" {
		return nil, errors.New("missing api key")
	}
	if len(apisecret) == 0 || apisecret == "undefined" {
		return nil, errors.New("missing api key secret")
	}
	//ctx, err := ValidateToken(apikey, apisecret)
	return ValidateToken(apikey, apisecret, requestID)
}

// requestWrapper generates API request ID for incoming request, and appends in request ID
func requestWrapper(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := uuid.New().String()
		r.Header.Set(RequestID, requestID)
		// ctx, err := commonAuthRules(
		// 	r.Header.Get(APIKey),
		// 	r.Header.Get(APIKeySecret), requestID)
		logMap := make(ablog.LogDataMap)
		//logMap[UserID] = ctx.Value(UserID)
		logMap[RequestID] = requestID
		// if err != nil {
		// 	ablog.Error(requestID, err, logMap)
		// 	WriteHTTPErrorResponse(w, "apicore", requestID, http.StatusUnauthorized, err)
		// 	return
		// }
		// r = r.Clone(ctx)
		next.ServeHTTP(w, r)
	})
}

//logRequest logs each HTTP incoming Requests
func logRequest(handler http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		m := httpsnoop.CaptureMetrics(handler, w, r)
		ablog.HTTPLog(constructHTTPLog(r, m, time.Since(start)))
	})
}

//constructHTTPLogMessage
func constructHTTPLog(r *http.Request, m httpsnoop.Metrics, duration time.Duration) string {
	return fmt.Sprintf("|%s|%s|%s|%s|%d|%d|%s|%s|",
		"requestId="+r.Header.Get(RequestID),
		r.RemoteAddr,
		r.Method,
		r.URL,
		m.Code,
		m.Written,
		r.UserAgent(),
		duration,
	)
}

// NewRouter provides a mux Router.
// Handles all incoming request who matches registered routes against the request.
func newRouter(subroute string) *mux.Router {

	muxRouter := mux.NewRouter().StrictSlash(true)
	subRouter := muxRouter.PathPrefix(subroute).Subrouter()
	for _, route := range routes {
		subRouter.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}

	return muxRouter
}

func useMiddleware(h http.HandlerFunc, middleware ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	for _, m := range middleware {
		h = m(h)
	}
	return h
}

var routes = make(Routes, 0)

func AddNoAuthRoutes(methodName string, methodType string, mRoute string, handlerFunc http.HandlerFunc) {
	r := route{
		Name:        methodName,
		Method:      methodType,
		Pattern:     mRoute,
		HandlerFunc: useMiddleware(handlerFunc, logRequest)}
	routes = append(routes, r)

}

func validateReq(h http.HandlerFunc, middleware ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	for _, m := range middleware {
		h = m(h)
	}
	return h
}

//Start - http servers
func Start(port, subroute string) {
	allowedOrigins := handlers.AllowedOrigins([]string{"*"}) // Allowing all origin as of now

	allowedHeaders := handlers.AllowedHeaders([]string{
		"X-Requested-With",
		"X-CSRF-Token",
		"X-Auth-Token",
		"Content-Type",
		"contentType",
		"Origin",
		"Authorization",
		"Accept",
		"Accept-Encoding",
		"timezone",
		"locale"})

	allowedMethods := handlers.AllowedMethods([]string{
		"POST",
		"GET",
		"DELETE",
		"PUT",
		"PATCH",
		"OPTIONS"})

	allowCredential := handlers.AllowCredentials()

	l.Fatal(http.ListenAndServe(":"+port, handlers.CORS(
		allowedHeaders,
		allowedMethods,
		allowedOrigins,
		allowCredential)(
		muxcontext.ClearHandler(
			newRouter(subroute),
		),
	),
	),
	)
}

func WriteHTTPResponse(w http.ResponseWriter, statusCode int, reqiestID, msg string, data interface{}) {
	renderer := render.New()
	res := Response{}
	res.Meta.RequestID = reqiestID
	res.Meta.Message = msg
	res.Meta.Code = statusCode
	res.Data = data
	renderer.JSON(w, statusCode, res)
}

func WriteHTTPDataResponse(w http.ResponseWriter, statusCode int, msg string, data []byte) {
	renderer := render.New()
	renderer.Data(w, statusCode, data)
}

func WriteHTTPErrorResponse(w http.ResponseWriter, reqID, msg string, errorCode int, err error) {
	renderer := render.New()
	res := Response{}
	res.Meta.Code = errorCode
	res.Meta.Message = msg
	res.Meta.RequestID = reqID
	res.Data = err.Error()
	renderer.JSON(w, errorCode, res)
}

// AddRoute is to create routes with ACL enforcer
func AddRoute(methodName, methodType, mRoute string, handlerFunc http.HandlerFunc) {
	r := route{
		Name:        methodName,
		Method:      methodType,
		Pattern:     mRoute,
		HandlerFunc: useMiddleware(handlerFunc, requestWrapper, logRequest),
	}
	routes = append(routes, r)
}

// AddRouteWithAuth is to create routes with ACL enforcer
func AddRouteWithAuth(methodName, methodType, mRoute string, handlerFunc http.HandlerFunc, authentication func(http.HandlerFunc) http.HandlerFunc) {
	r := route{
		Name:        methodName,
		Method:      methodType,
		Pattern:     mRoute,
		HandlerFunc: useMiddleware(handlerFunc, authentication, logRequest),
	}
	routes = append(routes, r)
}
