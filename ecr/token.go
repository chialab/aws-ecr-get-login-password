// Package ecr obtains Amazon ECR authorization tokens and splits them into username and password.
package ecr

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/aws/aws-sdk-go-v2/service/ecr/types"
)

// ErrInvalidToken is returned when token is not in base64-encoded `username:password` format.
var ErrInvalidToken = errors.New("unexpected token format")

// ErrUnexpectedResponse is returned when ECR client returns zero or more than one types.AuthorizationData.
var ErrUnexpectedResponse = errors.New("unexpected number of authorization data returned")

// Interface required to obtain an ECR authentication token.
type ecrTokenClient interface {
	GetAuthorizationToken(
		ctx context.Context,
		params *ecr.GetAuthorizationTokenInput,
		optFns ...func(*ecr.Options),
	) (*ecr.GetAuthorizationTokenOutput, error)
}

// getClient instantiates a new AWS ECR client using AWS default configuration with options.
func getClient(ctx context.Context, optsFns ...func(*config.LoadOptions) error) (*ecr.Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx, optsFns...)
	if err != nil {
		return nil, fmt.Errorf("loading AWS configuration: %w", err)
	}

	return ecr.NewFromConfig(cfg), nil
}

// Obtain a types.AuthorizationData from an ecrTokenClient.
func getAuthData(ctx context.Context, client ecrTokenClient) (*types.AuthorizationData, error) {
	token, err := client.GetAuthorizationToken(ctx, &ecr.GetAuthorizationTokenInput{})
	if err != nil {
		return nil, fmt.Errorf("getting ECR authorization token: %w", err)
	}

	if len(token.AuthorizationData) != 1 {
		return nil, ErrUnexpectedResponse
	}

	return &token.AuthorizationData[0], nil
}

// AuthorizationToken extends types.AuthorizationData by adding Username and Password explicitly.
type AuthorizationToken struct {
	types.AuthorizationData

	Username *string
	Password *string
}

// NewToken builds an AuthorizationToken from types.AuthorizationData.
func NewToken(authData *types.AuthorizationData) (*AuthorizationToken, error) {
	data, err := base64.StdEncoding.DecodeString(*authData.AuthorizationToken)
	if err != nil {
		return nil, fmt.Errorf("decoding token: %w", err)
	}

	username, password, found := strings.Cut(string(data), ":")
	if !found {
		return nil, ErrInvalidToken
	}

	return &AuthorizationToken{AuthorizationData: *authData, Username: &username, Password: &password}, nil
}

// GetToken obtains an AuthorizationToken from AWS ECR client using the passed context and options.
func GetToken(ctx context.Context, optsFns ...func(*config.LoadOptions) error) (*AuthorizationToken, error) {
	client, err := getClient(ctx, optsFns...)
	if err != nil {
		return nil, err
	}

	authData, err := getAuthData(ctx, client)
	if err != nil {
		return nil, err
	}

	return NewToken(authData)
}
