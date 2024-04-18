package ecr

import (
	"context"
	"encoding/base64"
	"errors"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/aws/aws-sdk-go-v2/service/ecr/types"
	"testing"
	"time"
)

func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}
func ind[T any](val T) *T { return &val }

type mockECRClient func(ctx context.Context, params *ecr.GetAuthorizationTokenInput, optFns ...func(*ecr.Options)) (*ecr.GetAuthorizationTokenOutput, error)

func (m mockECRClient) GetAuthorizationToken(ctx context.Context, params *ecr.GetAuthorizationTokenInput, optFns ...func(*ecr.Options)) (*ecr.GetAuthorizationTokenOutput, error) {
	return m(ctx, params, optFns...)
}

func TestGetAuthData(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		expected := types.AuthorizationData{AuthorizationToken: ind("Zm9vOmJhcg==")}
		client := mockECRClient(func(context.Context, *ecr.GetAuthorizationTokenInput, ...func(*ecr.Options)) (*ecr.GetAuthorizationTokenOutput, error) {
			return &ecr.GetAuthorizationTokenOutput{AuthorizationData: []types.AuthorizationData{expected}}, nil
		})

		authData, err := getAuthData(context.TODO(), client)
		if err != nil {
			t.Fatalf("expected nil error, got %v", err)
		}
		if authData == nil {
			t.Fatalf("expected authorization data, got nil")
		}

		if expected.AuthorizationToken != authData.AuthorizationToken {
			t.Errorf("expected authorization token %v, got %v", expected.AuthorizationToken, authData.AuthorizationToken)
		}
	})

	testError := errors.New("my error")
	fail := map[string]struct {
		client   ecrTokenClient
		expected error
	}{
		"Error": {
			mockECRClient(func(context.Context, *ecr.GetAuthorizationTokenInput, ...func(*ecr.Options)) (*ecr.GetAuthorizationTokenOutput, error) {
				return nil, testError
			}),
			testError,
		},
		"Invalid(EmptyAuthorizationDataArray)": {
			mockECRClient(func(context.Context, *ecr.GetAuthorizationTokenInput, ...func(*ecr.Options)) (*ecr.GetAuthorizationTokenOutput, error) {
				return &ecr.GetAuthorizationTokenOutput{AuthorizationData: []types.AuthorizationData{}}, nil
			}),
			UnexpectedResponseError,
		},
		"Invalid(AuthorizationDataArrayTooLong)": {
			mockECRClient(func(context.Context, *ecr.GetAuthorizationTokenInput, ...func(*ecr.Options)) (*ecr.GetAuthorizationTokenOutput, error) {
				return &ecr.GetAuthorizationTokenOutput{AuthorizationData: []types.AuthorizationData{{AuthorizationToken: ind("Zm9vOmJhcg==")}, {AuthorizationToken: ind("Zm9vOmJhcjpiYXo=")}}}, nil
			}),
			UnexpectedResponseError,
		},
	}
	for name, tt := range fail {
		t.Run(name, func(t *testing.T) {
			client, expected := tt.client, tt.expected
			authData, err := getAuthData(context.TODO(), client)
			if authData != nil {
				t.Fatalf("expected nil authorization data, got %v", authData)
			}
			if !errors.Is(err, expected) {
				t.Errorf("expected error %v, got %v", expected, err)
			}
		})
	}
}

func TestNewToken(t *testing.T) {
	success := map[string]struct {
		input    *types.AuthorizationData
		expected [2]string
	}{
		"Success(foo:bar)": {
			&types.AuthorizationData{AuthorizationToken: ind("Zm9vOmJhcg=="), ExpiresAt: ind(must(time.Parse(time.RFC3339, "2024-04-18T13:18:00Z")))},
			[2]string{"foo", "bar"},
		},
		"Success(foo:bar:baz)": {
			&types.AuthorizationData{AuthorizationToken: ind("Zm9vOmJhcjpiYXo="), ExpiresAt: ind(must(time.Parse(time.RFC3339, "2024-04-18T13:18:00Z")))},
			[2]string{"foo", "bar:baz"},
		},
	}

	for name, tt := range success {
		t.Run(name, func(t *testing.T) {
			input, username, password := tt.input, tt.expected[0], tt.expected[1]
			token, err := NewToken(input)
			if err != nil {
				t.Fatalf("expected nil error, got %v", err)
			}
			if token == nil {
				t.Fatalf("expected token, got nil")
			}

			if input.AuthorizationToken != token.AuthorizationToken {
				t.Errorf("expected authorization token %v, got %v", input.AuthorizationToken, token.AuthorizationToken)
			}
			if input.ExpiresAt != token.ExpiresAt {
				t.Errorf("expected expiration %s, got %s", input.ExpiresAt, token.ExpiresAt)
			}
			if input.ProxyEndpoint != token.ProxyEndpoint {
				t.Errorf("expected proxy endpoint %v, got %v", input.ProxyEndpoint, token.ProxyEndpoint)
			}
			if input.AuthorizationToken != token.AuthorizationToken {
				t.Errorf("expected authorization token %v, got %v", input.AuthorizationToken, token.AuthorizationToken)
			}
			if username != *token.Username {
				t.Errorf("expected username %s, got %s", username, *token.Username)
			}
			if password != *token.Password {
				t.Errorf("expected password %s, got %s", password, *token.Password)
			}
		})
	}

	fail := map[string]struct {
		input    *types.AuthorizationData
		expected error
	}{
		"Invalid(CorruptInput)": {
			&types.AuthorizationData{AuthorizationToken: ind("__NOT_A_BASE64_STRING!"), ExpiresAt: ind(must(time.Parse(time.RFC3339, "2024-04-18T13:18:00Z")))},
			base64.CorruptInputError(0),
		},
		"Invalid(`foo`)": {
			&types.AuthorizationData{AuthorizationToken: ind("Zm9v"), ExpiresAt: ind(must(time.Parse(time.RFC3339, "2024-04-18T13:18:00Z")))},
			InvalidTokenError,
		},
	}
	for name, tt := range fail {
		t.Run(name, func(t *testing.T) {
			input, expected := tt.input, tt.expected
			token, err := NewToken(input)
			if token != nil {
				t.Fatalf("expected nil token, got %v", token)
			}
			if !errors.Is(err, expected) {
				t.Errorf("expected error %v, got %v", expected, err)
			}
		})
	}
}
