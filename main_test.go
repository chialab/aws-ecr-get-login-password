package main

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/chialab/aws-ecr-get-login-password/ecr"
)

var errTest = errors.New("my error")

func TestRun(t *testing.T) {
	t.Parallel()

	password := "s3cr3t"
	succeed := func(context.Context) (*ecr.AuthorizationToken, error) {
		return &ecr.AuthorizationToken{Password: &password}, nil
	}
	fail := func(context.Context) (*ecr.AuthorizationToken, error) { return nil, errTest }

	tests := map[string]struct {
		args     []string
		getToken tokenSource
		code     int
		stdout   string
		stderr   string
	}{
		"Password": {nil, succeed, 0, "s3cr3t\n", ""},
		"Error":    {nil, fail, exitFailure, "", "aws-ecr-get-login-password: my error\n"},
		"Version":  {[]string{"--version"}, fail, 0, "dev\n", ""},
		"Help":     {[]string{"-h"}, fail, 0, "", "Usage: aws-ecr-get-login-password"},
		"UnknownFlag": {
			[]string{"--region=eu-west-1"},
			fail, exitUsage, "",
			"aws-ecr-get-login-password: flag provided but not defined: -region\nUsage:",
		},
		"UnexpectedArgs": {
			[]string{"foo"}, fail, exitUsage, "", "aws-ecr-get-login-password: unexpected arguments: [\"foo\"]\nUsage:",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var stdout, stderr bytes.Buffer

			code := run(t.Context(), tt.args, &stdout, &stderr, tt.getToken)
			if code != tt.code {
				t.Errorf("expected exit code %d, got %d", tt.code, code)
			}

			if stdout.String() != tt.stdout {
				t.Errorf("expected stdout %q, got %q", tt.stdout, stdout.String())
			}

			if !strings.Contains(stderr.String(), tt.stderr) || (tt.stderr == "") != (stderr.Len() == 0) {
				t.Errorf("expected stderr containing %q, got %q", tt.stderr, stderr.String())
			}
		})
	}
}
