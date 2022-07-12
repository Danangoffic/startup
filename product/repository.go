package product

import "gorm.io/gorm"

type RepositoryProduct interface {
	Save(product Product) (Product, error)
}

type repositoryproduct struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repositoryproduct {
	return &repositoryproduct{db}
}

func (r *repositoryproduct) Save(product Product) (Product, error) {
	err := r.db.Create(&product).Error
	if err != nil {
		return product, err
	}
	return product, nil
}
