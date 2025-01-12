package handler

import (
	"context"
	"gc2-yugo/pb"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func BorrowBook(c echo.Context) error {
	// Extract book ID from the URL
	bookID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid book ID format")
	}

	// Retrieve token from the request header
	token := c.Request().Header.Get("Authorization")
	if token == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "Missing token")
	}

	// Create gRPC connection
	grpcConn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to connect to gRPC server")
	}
	defer grpcConn.Close()

	client := pb.NewBookRentalServiceClient(grpcConn)

	// Add token to metadata
	md := metadata.Pairs("authorization", token)
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	// Prepare BorrowBookRequest
	req := &pb.BorrowBookRequest{
		BookId: bookID.Hex(),
	}

	// Call BorrowBook gRPC service
	resp, err := client.BorrowBook(ctx, req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// Return success response
	return c.JSON(http.StatusOK, resp)
}
