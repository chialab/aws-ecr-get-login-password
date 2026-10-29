package ecr

import (
	"context"
	"encoding/base64"
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/aws/aws-sdk-go-v2/service/ecr/types"
)

var errTest = errors.New("my error")

type mockECRClient struct {
	output *ecr.GetAuthorizationTokenOutput
	err    error
}

func (m mockECRClient) GetAuthorizationToken(
	context.Context,
	*ecr.GetAuthorizationTokenInput,
	...func(*ecr.Options),
) (*ecr.GetAuthorizationTokenOutput, error) {
	return m.output, m.err
}

func authData(tokens ...string) *ecr.GetAuthorizationTokenOutput {
	output := &ecr.GetAuthorizationTokenOutput{AuthorizationData: []types.AuthorizationData{}}
	for _, token := range tokens {
		output.AuthorizationData = append(output.AuthorizationData, types.AuthorizationData{AuthorizationToken: &token})
	}

	return output
}

func TestGetAuthData(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		client   mockECRClient
		expected string
		err      error
	}{
		"Success":                              {mockECRClient{authData("Zm9vOmJhcg=="), nil}, "Zm9vOmJhcg==", nil},
		"Error":                                {mockECRClient{nil, errTest}, "", errTest},
		"Invalid(EmptyAuthorizationDataArray)": {mockECRClient{authData(), nil}, "", ErrUnexpectedResponse},
		"Invalid(AuthorizationDataArrayTooLong)": {
			mockECRClient{authData("Zm9vOmJhcg==", "Zm9vOmJhcjpiYXo="), nil}, "", ErrUnexpectedResponse,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			data, err := getAuthData(t.Context(), tt.client)
			if !errors.Is(err, tt.err) {
				t.Fatalf("expected error %v, got %v", tt.err, err)
			}

			if tt.err != nil {
				if data != nil {
					t.Errorf("expected nil authorization data, got %v", data)
				}

				return
			}

			if *data.AuthorizationToken != tt.expected {
				t.Errorf("expected authorization token %s, got %s", tt.expected, *data.AuthorizationToken)
			}
		})
	}
}

func TestNewToken(t *testing.T) {
	t.Parallel()

	expiresAt := time.Date(2024, 4, 18, 13, 18, 0, 0, time.UTC)

	tests := map[string]struct {
		token    string
		username string
		password string
		err      error
	}{
		"Success(foo:bar)":      {"Zm9vOmJhcg==", "foo", "bar", nil},
		"Success(foo:bar:baz)":  {"Zm9vOmJhcjpiYXo=", "foo", "bar:baz", nil},
		"Invalid(CorruptInput)": {"__NOT_A_BASE64_STRING!", "", "", base64.CorruptInputError(0)},
		"Invalid(`foo`)":        {"Zm9v", "", "", ErrInvalidToken},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			input := &types.AuthorizationData{AuthorizationToken: &tt.token, ExpiresAt: &expiresAt}

			token, err := NewToken(input)
			if !errors.Is(err, tt.err) {
				t.Fatalf("expected error %v, got %v", tt.err, err)
			}

			if tt.err != nil {
				if token != nil {
					t.Errorf("expected nil token, got %v", token)
				}

				return
			}

			if token.AuthorizationToken != input.AuthorizationToken || token.ExpiresAt != input.ExpiresAt {
				t.Errorf("expected authorization data %v, got %v", *input, token.AuthorizationData)
			}

			if *token.Username != tt.username || *token.Password != tt.password {
				t.Errorf("expected %s:%s, got %s:%s", tt.username, tt.password, *token.Username, *token.Password)
			}
		})
	}
}
