package main

import (
	"context"
	"log"
	"net"
	"testing"

	pb "grpc_user_service/generated"

	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

func dialer() func(context.Context, string) (net.Conn, error) {
	lis := bufconn.Listen(1024 * 1024)
	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, &userServiceServer{
		users: []*pb.User{
			{Id: 1, Fname: "Steve", City: "LA", Phone: 1234567890, Height: 5.8, Married: true},
		},
	})
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()
	return func(ctx context.Context, address string) (net.Conn, error) {
		return lis.Dial()
	}
}

func TestGetUser(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(dialer()), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()
	client := pb.NewUserServiceClient(conn)

	req := &pb.GetUserRequest{Id: 1}
	resp, err := client.GetUser(ctx, req)
	if err != nil {
		t.Fatalf("GetUser failed: %v", err)
	}
	if resp.User == nil || resp.User.Fname != "Steve" {
		t.Fatalf("GetUser returned unexpected result: %v", resp.User)
	}
}

func TestListUsers(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(dialer()), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()
	client := pb.NewUserServiceClient(conn)

	req := &pb.ListUsersRequest{Ids: []int32{1}}
	resp, err := client.ListUsers(ctx, req)
	if err != nil {
		t.Fatalf("ListUsers failed: %v", err)
	}
	if len(resp.Users) != 1 || resp.Users[0].Fname != "Steve" {
		t.Fatalf("ListUsers returned unexpected result: %v", resp.Users)
	}
}

func TestSearchUsers(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(dialer()), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()
	client := pb.NewUserServiceClient(conn)

	req := &pb.SearchUsersRequest{City: "LA", Married: true}
	resp, err := client.SearchUsers(ctx, req)
	if err != nil {
		t.Fatalf("SearchUsers failed: %v", err)
	}
	if len(resp.Users) != 1 || resp.Users[0].Fname != "Steve" {
		t.Fatalf("SearchUsers returned unexpected result: %v", resp.Users)
	}
}
