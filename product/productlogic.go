package product

import (
	"encoding/json"
	"strings"

	"github.com/pontiyaraja/AB/product/core"
)

func getProductData(productID string) (*Product, error) {
	data, err := core.HGet("Product", productID)
	if err != nil {
		return nil, err
	}
	var prod Product
	if err := json.NewDecoder(strings.NewReader(data)).Decode(&prod); err != nil {
		return nil, err
	}
	return &prod, err
}
