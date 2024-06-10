package main

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "grpc_user_service/generated"

	"google.golang.org/grpc"
)

type userServiceServer struct {
	pb.UnimplementedUserServiceServer
	users []*pb.User
}

func (s *userServiceServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	for _, user := range s.users {
		if user.Id == req.Id {
			return &pb.GetUserResponse{User: user}, nil
		}
	}
	return nil, fmt.Errorf("user not found")
}

func (s *userServiceServer) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	var users []*pb.User
	for _, id := range req.Ids {
		for _, user := range s.users {
			if user.Id == id {
				users = append(users, user)
			}
		}
	}
	return &pb.ListUsersResponse{Users: users}, nil
}

func (s *userServiceServer) SearchUsers(ctx context.Context, req *pb.SearchUsersRequest) (*pb.SearchUsersResponse, error) {
	var users []*pb.User
	for _, user := range s.users {
		if (req.City == "" || user.City == req.City) &&
			(req.Phone == 0 || user.Phone == req.Phone) &&
			(!req.Married || user.Married == req.Married) {
			users = append(users, user)
		}
	}
	return &pb.SearchUsersResponse{Users: users}, nil
}

func main() {
	users := []*pb.User{
		{Id: 1, Fname: "Steve", City: "LA", Phone: 1234567890, Height: 5.8, Married: true},
	}

	server := grpc.NewServer()
	pb.RegisterUserServiceServer(server, &userServiceServer{users: users})

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("server listening at %v", lis.Addr())
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
