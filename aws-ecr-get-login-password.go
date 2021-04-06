package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
)

func main() {
	os.Setenv("AWS_SDK_LOAD_CONFIG", "1")
	sess := session.Must(session.NewSession())
	svc := ecr.New(sess)

	input := &ecr.GetAuthorizationTokenInput{}
	token, err := svc.GetAuthorizationToken(input)
	if err != nil {
		log.Fatal(err)
	}

	authData := token.AuthorizationData[0].AuthorizationToken
	data, err := base64.StdEncoding.DecodeString(*authData)
	if err != nil {
		log.Fatal(err)
	}

	parts := strings.SplitN(string(data), ":", 2)
	fmt.Println(parts[1])
}
