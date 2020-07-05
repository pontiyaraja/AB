package product

import (
	"net/http"

	"github.com/pontiyaraja/AB/product/core"
)

func Start() {
	core.AddNoAuthRoutes(
		"product",
		http.MethodGet,
		"/product/{ID}",
		getProduct,
	)

}
