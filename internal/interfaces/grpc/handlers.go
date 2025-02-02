package grpc

// import (
// 	"backend/internal/shared/logging"
// 	pb "backend/shared/proto"
// )

// type AssetServiceHandler struct {
// 	pb.UnimplementedAssetServiceServer
// 	assetService application.AssetService
// 	logger       *logging.Logger
// }

// func NewAssetServiceHandler(assetService application.AssetService, logger *logging.Logger) *AssetServiceHandler {
// 	return &AssetServiceHandler{
// 		assetService: assetService,
// 		logger:       logger,
// 	}
// }

// // Implement your gRPC service methods here
