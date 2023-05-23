package app

import (
	"context"
	"fmt"
	"net"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"

	"awesomeProject/internal/database"
	proto "awesomeProject/pb"
)

type OrderServerImpl struct {
	proto.UnimplementedOrderServiceServer
}

// UpdateOrderToDone is implement here.
func (s *OrderServerImpl) UpdateOrderToDone(ctx context.Context, in *proto.OrderRequest) (*proto.OrderResponse, error) {
	if err := database.UpdateOrderDoneByID(in.Id); err != nil {
		return nil, err
	}

	return &proto.OrderResponse{}, nil
}

// Start GRPC Server
func startGRPCServer() {
	var address = fmt.Sprintf("%s:%d", OrderServiceConfig.GRPC.Address, OrderServiceConfig.GRPC.Port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot start GRPC server")
	}

	s := grpc.NewServer()
	proto.RegisterOrderServiceServer(s, &OrderServerImpl{})

	log.Info().Msg("GRPC server listening on " + address)
	if err := s.Serve(listener); err != nil {
		log.Fatal().Err(err).Msg("Cannot start GRPC server")
	}
}
