package handler

import (
	"context"

	model "notionboy/api/pb/model"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Server) CreateProduct(ctx context.Context, req *model.CreateProductRequest) (*model.Product, error) {
	acc := getAccFromContext(ctx)
	if acc == nil || !acc.IsAdmin {
		return nil, status.Errorf(codes.PermissionDenied, "permission denied")
	}
	return s.ProductService.CreateProduct(ctx, acc, req)
}

func (s *Server) GetProduct(ctx context.Context, req *model.GetProductRequest) (*model.Product, error) {
	return s.ProductService.GetProduct(ctx, nil, req)
}

func (s *Server) ListProducts(ctx context.Context, req *model.ListProductsRequest) (*model.ListProductsResponse, error) {
	return s.ProductService.ListProducts(ctx, nil, req)
}

func (s *Server) DeleteProduct(ctx context.Context, req *model.DeleteProductRequest) (*emptypb.Empty, error) {
	acc := getAccFromContext(ctx)
	if acc == nil || !acc.IsAdmin {
		return nil, status.Errorf(codes.PermissionDenied, "permission denied")
	}

	return s.ProductService.DeleteProduct(ctx, acc, req)
}

func (s *Server) UpdateProduct(ctx context.Context, req *model.UpdateProductRequest) (*model.Product, error) {
	acc := getAccFromContext(ctx)
	if acc == nil || !acc.IsAdmin {
		return nil, status.Errorf(codes.PermissionDenied, "permission denied")
	}
	return s.ProductService.UpdateProduct(ctx, acc, req)
}
