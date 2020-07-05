package product

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pontiyaraja/AB/core"
)

func getProduct(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	productID := vars["ID"]
	resp, err := getProductData(productID)
	if err != nil {
		core.WriteHTTPErrorResponse(w, "123", "failed to get data", http.StatusNoContent, err)
		return
	}
	core.WriteHTTPResponse(w, 200, "123", "success", resp)
}
