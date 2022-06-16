Get AWS ECR login password
==========================

This is a simple utility written in Go to build a portable binary whose sole purpose
is to mimick `aws ecr get-login-password`. This is useful in contexts like
Alpine Linux based container images, where glibc cannot be used and AWS CLI v2 doesn't work.

Usage
-----

```console
$ aws-ecr-get-login-password | docker login -u AWS --password-stdin 000000000000.dkr.ecr.eu-south-1.amazonaws.com
```

> :warning: Remember to replace `000000000000.dkr.ecr.eu-south-1.amazonaws.com` with your actual ECR repository!

Credentials
-----------

The binary is as simple as it can be. It does not accept any option, but it infers
configuration from the environment using most of the environment variables AWS CLI
uses as well. A non-exhaustive list of environment variables that can be set to
change the way this binary retrieves credentials are:

- `AWS_DEFAULT_REGION`
- `AWS_PROFILE`
- `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY` and optionally `AWS_SESSION_TOKEN`
- `AWS_IAM_ROLE`
- `AWS_WEB_IDENTITY_TOKEN_FILE`
