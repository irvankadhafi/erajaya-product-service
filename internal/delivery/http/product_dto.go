package http

import "github.com/irvankadhafi/erajaya-product-service/internal/model"

type productResponse struct {
	*model.Product
	Price     string `json:"price"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type createProductRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       string `json:"price"`
	Quantity    int    `json:"quantity"`
}
