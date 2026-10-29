// Command aws-ecr-get-login-password prints an Amazon ECR login password to stdout,
// like "aws ecr get-login-password", for systems where the AWS CLI can't run.
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/chialab/aws-ecr-get-login-password/ecr"
	// Embedded root CAs, used only when the system has none (e.g. scratch images).
	_ "golang.org/x/crypto/x509roots/fallback"
)

const (
	exitFailure = 1
	exitUsage   = 2
)

var version = "dev"

var errUnexpectedArgs = errors.New("unexpected arguments")

const usage = `Usage: aws-ecr-get-login-password [options]

Print an Amazon ECR login password to stdout, like "aws ecr get-login-password".
Credentials and region are read from the environment, as the AWS CLI does.

Options:
      --version  print version and exit
  -h, --help     print this help and exit
`

// tokenSource obtains the token to print: ecr.GetToken, or a fake in tests.
type tokenSource func(ctx context.Context) (*ecr.AuthorizationToken, error)

func main() {
	getToken := func(ctx context.Context) (*ecr.AuthorizationToken, error) { return ecr.GetToken(ctx) }

	os.Exit(run(context.Background(), os.Args[1:], os.Stdout, os.Stderr, getToken))
}

// warn reports to stderr; there is nowhere left to report a failed write.
func warn(stderr io.Writer, err error) {
	_, _ = fmt.Fprintf(stderr, "aws-ecr-get-login-password: %v\n", err)
}

func run(ctx context.Context, args []string, stdout, stderr io.Writer, getToken tokenSource) int {
	var showVersion bool

	flags := flag.NewFlagSet("aws-ecr-get-login-password", flag.ContinueOnError)
	flags.SetOutput(stderr)
	flags.Usage = func() { _, _ = fmt.Fprint(flags.Output(), usage) }
	flags.BoolVar(&showVersion, "version", false, "")

	err := flags.Parse(args)
	if errors.Is(err, flag.ErrHelp) {
		return 0
	}

	if err != nil {
		return exitUsage
	}

	if flags.NArg() > 0 {
		warn(stderr, fmt.Errorf("%w: %q", errUnexpectedArgs, flags.Args()))

		return exitUsage
	}

	if showVersion {
		_, _ = fmt.Fprintln(stdout, version)

		return 0
	}

	token, err := getToken(ctx)
	if err != nil {
		warn(stderr, err)

		return exitFailure
	}

	_, err = fmt.Fprintln(stdout, *token.Password)
	if err != nil {
		warn(stderr, err)

		return exitFailure
	}

	return 0
}
