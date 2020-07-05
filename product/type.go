package product

type Product struct {
	ProductID             int64  `json:"productID" bson:"productID"`
	ProductName           string `json:"productName" bson:"productName"`
	ProductTag            string `json:"productTag" bson:"productTag"`
	IsLightingDealEnabled bool   `json:"isLightingDealEnabled" bson:"isLightingDealEnabled"`
	Quantity              int    `json:"quantity" bson:"quantity"`
}
