package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
)

func main() {
	if token, err := getToken(); err != nil {
		log.Fatal(err)
	} else {
		fmt.Println(token)
	}
}

// Retrieve token for authentication against ECR registries.
func getToken() (string, error) {
	os.Setenv("AWS_SDK_LOAD_CONFIG", "1")
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return "", err
	}

	svc := ecr.NewFromConfig(cfg)
	token, err := svc.GetAuthorizationToken(context.TODO(), &ecr.GetAuthorizationTokenInput{})
	if err != nil {
		return "", err
	}

	authData := token.AuthorizationData[0].AuthorizationToken
	data, err := base64.StdEncoding.DecodeString(*authData)
	if err != nil {
		return "", err
	}

	parts := strings.SplitN(string(data), ":", 2)

	return parts[1], nil
}
