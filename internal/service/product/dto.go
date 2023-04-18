package product

import (
	"time"

	"notionboy/api/pb/model"
	"notionboy/db/ent"
)

type ProductDTO struct {
	*ent.Product
}

func NewProductDTO(p *ent.Product) *ProductDTO {
	return &ProductDTO{p}
}

func (p *ProductDTO) ToProto() *model.Product {
	return &model.Product{
		Id:          p.UUID.String(),
		CreatedAt:   p.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   p.UpdatedAt.Format(time.RFC3339),
		Uuid:        p.UUID.String(),
		Name:        p.Name,
		Price:       float32(p.Price),
		Token:       p.Token,
		Storage:     p.Storage,
		Description: p.Description,
	}
}
