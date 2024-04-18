package main

import (
	"context"
	"fmt"
	"github.com/chialab/aws-ecr-get-login-password/ecr"
	"log"
)

func main() {
	if token, err := ecr.GetToken(context.Background()); err != nil {
		log.Fatal(err)
	} else {
		fmt.Println(*token.Password)
	}
}
