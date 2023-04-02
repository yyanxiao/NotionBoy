package order

import (
	"notionboy/api/pb/model"
	"notionboy/db/ent"
	"notionboy/internal/service/product"
	"time"
)

type OrderDTO struct {
	*ent.Order
	Product *product.ProductDTO
}

func NewOrderDTO(o *ent.Order, p *ent.Product) *OrderDTO {
	if p == nil {
		return &OrderDTO{o, nil}
	}
	return &OrderDTO{o, product.NewProductDTO(p)}
}

func (o *OrderDTO) ToProto() *model.Order {
	var p *model.Product
	if o.Product != nil {
		p = o.Product.ToProto()
	}

	return &model.Order{
		Id:        o.UUID.String(),
		CreatedAt: o.CreatedAt.Format(time.RFC3339),
		UpdatedAt: o.UpdatedAt.Format(time.RFC3339),
		Uuid:      o.UUID.String(),
		UserId:    o.UserID.String(),
		Price:     float32(o.Price),
		Status:    o.Status.String(),
		Note:      o.Note,
		Product:   p,
	}
}
