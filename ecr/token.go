package ecr

import (
	"context"
	"encoding/base64"
	"errors"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/aws/aws-sdk-go-v2/service/ecr/types"
	"strings"
)

// InvalidTokenError is returned when token is not in base64-encoded `username:password` format.
var InvalidTokenError = errors.New("unexpected token format")

// UnexpectedResponseError is returned when ECR client returns zero or more than one types.AuthorizationData in its response.
var UnexpectedResponseError = errors.New("unexpected number of authorization data returned")

// Interface required to obtain an ECR authentication token.
type ecrTokenClient interface {
	GetAuthorizationToken(ctx context.Context, params *ecr.GetAuthorizationTokenInput, optFns ...func(*ecr.Options)) (*ecr.GetAuthorizationTokenOutput, error)
}

// getClient returns an ecrTokenClient instantiating a new AWS ECR client using AWS default configuration with options.
func getClient(context context.Context, optsFns ...func(*config.LoadOptions) error) (ecrTokenClient, error) {
	if cfg, err := config.LoadDefaultConfig(context, optsFns...); err != nil {
		return nil, err
	} else {
		return ecr.NewFromConfig(cfg), nil
	}
}

// Obtain a types.AuthorizationData from an ecrTokenClient.
func getAuthData(context context.Context, client ecrTokenClient) (*types.AuthorizationData, error) {
	if token, err := client.GetAuthorizationToken(context, &ecr.GetAuthorizationTokenInput{}); err != nil {
		return nil, err
	} else if len(token.AuthorizationData) != 1 {
		return nil, UnexpectedResponseError
	} else {
		return &token.AuthorizationData[0], nil
	}
}

// AuthorizationToken extends types.AuthorizationData by adding Username and Password explicitly.
type AuthorizationToken struct {
	types.AuthorizationData
	Username *string
	Password *string
}

// NewToken builds an AuthorizationToken from types.AuthorizationData.
func NewToken(authData *types.AuthorizationData) (*AuthorizationToken, error) {
	token := *authData.AuthorizationToken
	if data, err := base64.StdEncoding.DecodeString(token); err != nil {
		return nil, err
	} else if parts := strings.SplitN(string(data), ":", 2); len(parts) != 2 {
		return nil, InvalidTokenError
	} else {
		return &AuthorizationToken{*authData, &parts[0], &parts[1]}, nil
	}
}

// GetToken obtains an AuthorizationToken from AWS ECR client using the passed context and options
func GetToken(context context.Context, optsFns ...func(*config.LoadOptions) error) (*AuthorizationToken, error) {
	if client, err := getClient(context, optsFns...); err != nil {
		return nil, err
	} else if authData, err := getAuthData(context, client); err != nil {
		return nil, err
	} else {
		return NewToken(authData)
	}
}
