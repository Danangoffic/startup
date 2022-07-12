package product

type ServiceProduct interface {
	RegisterProduct(input RegisterProductInput) (Product, error)
}

type serviceproduct struct {
	repository RepositoryProduct
}

func NewProductService(repository RepositoryProduct) *serviceproduct {
	return &serviceproduct{repository}
}

func (s *serviceproduct) RegisterProduct(input RegisterProductInput) (Product, error) {
	product := Product{}
	product.Name = input.Name
	product.Desc = input.Desc
	product.SKU = input.SKU
	product.Category = input.Category

	NewProduct, err := s.repository.Save(product)
	if err != nil {
		return NewProduct, err
	}
	return NewProduct, nil
}
