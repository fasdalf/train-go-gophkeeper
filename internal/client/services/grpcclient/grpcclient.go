// Package grpcclient - grpc client
package grpcclient

import (
	"context"
	"fmt"
	pb "github.com/fasdalf/train-go-gophkeeper/internal/common/proto/gophkeeperserver"
	"github.com/fasdalf/train-go-gophkeeper/internal/server/jwt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type Service struct {
	client pb.GophKeeperClient
	token  string
}

func NewService(addr string) *Service {
	conn, _ := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	client := pb.NewGophKeeperClient(conn)
	return &Service{client: client}
}

func (s *Service) SignUp(login, pass string) (string, error) {
	resp, err := s.client.UserSignUp(context.Background(), &pb.UserSignUpRequest{
		Login:    login,
		Password: pass,
	})
	if err != nil {
		return "", fmt.Errorf("could not sign up user: %w", err)
	}
	if resp == nil || resp.Token == "" {
		return "", fmt.Errorf("JWT is empty")
	}
	return resp.Token, nil
}

func (s *Service) SignIn(login, pass string) (string, error) {
	resp, err := s.client.UserSignIn(context.Background(), &pb.UserSignInRequest{
		Login:    login,
		Password: pass,
	})
	if err != nil {
		return "", fmt.Errorf("could not sign in user: %w", err)
	}
	if resp == nil || resp.Token == "" {
		return "", fmt.Errorf("JWT is empty")
	}
	return resp.Token, nil
}

func (s *Service) SetToken(token string) {
	s.token = token
}

func (s *Service) contextWithToken() context.Context {
	mdMap := map[string]string{
		jwt.AuthHeader: jwt.AuthPrefix + s.token,
	}
	return metadata.NewOutgoingContext(context.Background(), metadata.New(mdMap))
}

func (s *Service) ListSecrets(fromTs int64) (*pb.ListSecretsResponse, error) {
	res, err := s.client.ListSecrets(s.contextWithToken(), &pb.ListSecretsRequest{FromTs: fromTs})
	if err != nil || res == nil {
		err = fmt.Errorf("could not list secrets: %w", err)
	}
	return res, err
}

func (s *Service) CreateSecret(data *[]byte) (id, ts int64, err error) {
	res, err := s.client.CreateSecret(s.contextWithToken(), &pb.CreateSecretRequest{Data: *data})
	if err != nil || res == nil {
		err = fmt.Errorf("could not create secret: %w", err)
	}
	return res.GetId(), res.GetTs(), err
}

func (s *Service) UpdateSecret(id, oldTs int64, data *[]byte) (ts int64, err error) {
	res, err := s.client.UpdateSecret(s.contextWithToken(), &pb.UpdateSecretRequest{
		Id:    id,
		OldTs: oldTs,
		Data:  *data,
	})
	if err != nil || res == nil || res.Status != pb.TimeStampStatus_SUCCESS {
		if res != nil {
			err = fmt.Errorf("timestamp status is %s", res.Status.String())
		}
		err = fmt.Errorf("could not create secret: %w", err)
	}
	return res.GetNewTs(), err
}

func (s *Service) DeleteSecret(id, oldTs int64) (err error) {
	res, err := s.client.DeleteSecret(s.contextWithToken(), &pb.DeleteSecretRequest{
		Id:    id,
		OldTs: oldTs,
	})
	if err != nil || res == nil || res.Status != pb.TimeStampStatus_SUCCESS {
		if res != nil {
			err = fmt.Errorf("timestamp status is %s", res.Status.String())
		}
		err = fmt.Errorf("could not create secret: %w", err)
	}
	return err
}
