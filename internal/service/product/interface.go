package product

import (
	"context"
	"notionboy/api/pb/model"
	"notionboy/db/ent"

	"google.golang.org/protobuf/types/known/emptypb"
)

type ProductServiceImpl struct{}

type ProductService interface {
	CreateProduct(ctx context.Context, acc *ent.Account, req *model.CreateProductRequest) (*model.Product, error)
	GetProduct(ctx context.Context, acc *ent.Account, req *model.GetProductRequest) (*model.Product, error)
	ListProducts(ctx context.Context, acc *ent.Account, req *model.ListProductsRequest) (*model.ListProductsResponse, error)
	DeleteProduct(ctx context.Context, acc *ent.Account, req *model.DeleteProductRequest) (*emptypb.Empty, error)
	UpdateProduct(ctx context.Context, acc *ent.Account, req *model.UpdateProductRequest) (*model.Product, error)
}

func NewProductService() ProductService {
	return &ProductServiceImpl{}
}
