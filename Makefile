SHELL:=/bin/bash

mock_product_repository:
	mockgen -destination=internal/model/mock/mock_product_repository.go -package=mock github.com/irvankadhafi/erajaya-product-service/internal/model ProductRepository
mock_product_usecase:
	mockgen -destination=internal/model/mock/mock_product_usecase.go -package=mock github.com/irvankadhafi/erajaya-product-service/internal/model ProductUsecase