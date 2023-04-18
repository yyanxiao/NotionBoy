package handler

import (
	"context"

	"notionboy/api/pb"

	status "google.golang.org/genproto/googleapis/rpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

func (s *Server) Status(ctx context.Context, req *emptypb.Empty) (*pb.CheckStatusResponse, error) {
	return &pb.CheckStatusResponse{
		Status: &status.Status{
			Code:    0,
			Message: "OK",
		},
	}, nil
}
