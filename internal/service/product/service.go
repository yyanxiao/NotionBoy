package product

import (
	"context"
	"notionboy/api/pb/model"
	"notionboy/db/ent"
	"notionboy/db/ent/product"
	"notionboy/internal/pkg/db"
	"notionboy/internal/pkg/db/dao"

	"github.com/google/uuid"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *ProductServiceImpl) CreateProduct(ctx context.Context, acc *ent.Account, req *model.CreateProductRequest) (*model.Product, error) {
	client := db.GetClient()
	p, err := client.Product.Create().
		SetUUID(uuid.New()).
		SetName(req.Name).
		SetPrice(float64(req.Price)).
		SetToken(req.Token).
		SetStorage(req.Storage).
		SetDescription(req.Description).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return NewProductDTO(p).ToProto(), nil
}

func (s *ProductServiceImpl) GetProduct(ctx context.Context, acc *ent.Account, req *model.GetProductRequest) (*model.Product, error) {
	uid, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(400, "invalid id, %s", err.Error())
	}
	p, err := dao.GetProduct(ctx, uid)
	if err != nil {
		return nil, err
	}
	return NewProductDTO(p).ToProto(), nil
}

func (s *ProductServiceImpl) ListProducts(ctx context.Context, acc *ent.Account, req *model.ListProductsRequest) (*model.ListProductsResponse, error) {
	products, err := dao.ListProducts(ctx, 10, 0)
	if err != nil {
		return nil, err
	}
	productsDTO := make([]*model.Product, 0)
	for _, p := range products {
		productsDTO = append(productsDTO, NewProductDTO(p).ToProto())
	}
	return &model.ListProductsResponse{
		Products: productsDTO,
	}, nil
}

func (s *ProductServiceImpl) DeleteProduct(ctx context.Context, acc *ent.Account, req *model.DeleteProductRequest) (*emptypb.Empty, error) {
	uid, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(400, "invalid id, %s", err.Error())
	}
	err = dao.DeleteProduct(ctx, uid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *ProductServiceImpl) UpdateProduct(ctx context.Context, acc *ent.Account, req *model.UpdateProductRequest) (*model.Product, error) {
	uid, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(400, "invalid id, %s", err.Error())
	}

	query := db.GetClient().Product.Update().
		Where(product.UUIDEQ(uid))

	if req.Name != "" {
		query = query.SetName(req.Name)
	}

	if req.Price != 0 {
		query = query.SetPrice(float64(req.Price))
	}

	if req.Token != 0 {
		query = query.SetToken(req.Token)
	}

	if req.Storage != 0 {
		query = query.SetStorage(req.Storage)
	}

	if req.Description != "" {
		query = query.SetDescription(req.Description)
	}

	if err := query.Exec(ctx); err != nil {
		return nil, err
	}

	p, err := dao.GetProduct(ctx, uid)
	if err != nil {
		return nil, err
	}
	return NewProductDTO(p).ToProto(), nil
}
