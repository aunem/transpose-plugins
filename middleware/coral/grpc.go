package main

import (
	"github.com/aunem/coral/sdk/go/auth"
	"google.golang.org/grpc"
)

func getCoralClient(addr string) (a auth.AuthServiceClient, err error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return a, err
	}
	c := auth.NewAuthServiceClient(conn)
	return c, err
}
