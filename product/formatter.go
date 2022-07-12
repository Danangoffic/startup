package product

type ProductFormatter struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func FormatProduct(product Product) ProductFormatter {
	formatter := ProductFormatter{
		ID:   product.Id,
		Name: product.Name,
	}
	return formatter
}
