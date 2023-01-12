SHELL:=/bin/bash

ifdef test_run
	TEST_ARGS := -run $(test_run)
endif


run:
	go run . server

migration:
	go run . migrate

mock_product_repository:
	mockgen -destination=internal/model/mock/mock_product_repository.go -package=mock github.com/irvankadhafi/erajaya-product-service/internal/model ProductRepository

mock_product_usecase:
	mockgen -destination=internal/model/mock/mock_product_usecase.go -package=mock github.com/irvankadhafi/erajaya-product-service/internal/model ProductUsecase

mockgen: mock_product_repository \
	mock_product_usecase

test: check-gotest mockgen
	SVC_DISABLE_CACHING=true richgo test ./... $(TEST_ARGS) -v --cover

check-gotest:
ifeq (, $(shell which richgo))
	$(warning "richgo is not installed, falling back to plain go test")
	$(eval TEST_BIN=go test)
else
	$(eval TEST_BIN=richgo test)
endif

ifdef test_run
	$(eval TEST_ARGS := -run $(test_run))
endif
	$(eval test_command=$(TEST_BIN) ./... $(TEST_ARGS) -v --cover)

docker:
	@ docker-compose up -d --build