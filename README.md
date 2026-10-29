# aws-ecr-get-login-password

[![CI](https://github.com/chialab/aws-ecr-get-login-password/actions/workflows/ci.yml/badge.svg)](https://github.com/chialab/aws-ecr-get-login-password/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/chialab/aws-ecr-get-login-password/graph/badge.svg)](https://codecov.io/gh/chialab/aws-ecr-get-login-password)
[![Release](https://img.shields.io/github/v/release/chialab/aws-ecr-get-login-password?sort=semver)](https://github.com/chialab/aws-ecr-get-login-password/releases/latest)
[![Go version](https://img.shields.io/github/go-mod/go-version/chialab/aws-ecr-get-login-password)](go.mod)
[![License: MIT](https://img.shields.io/github/license/chialab/aws-ecr-get-login-password)](LICENSE)

**Log in to Amazon ECR where the AWS CLI can't run.**

`aws-ecr-get-login-password` does one thing: it prints an ECR login password to stdout, exactly like `aws ecr get-login-password`.

- 🪶 **Tiny**: a single static binary, instead of the whole AWS CLI and its Python or glibc runtime.
- 🐧 **Runs anywhere**: Alpine/musl, `scratch` and distroless images, on Linux, macOS and Windows, with embedded root CAs for minimal images.
- 🔑 **Same credentials as the AWS CLI**: environment variables, profiles, web identity, container and instance roles.

```sh
aws-ecr-get-login-password | docker login --username AWS --password-stdin 000000000000.dkr.ecr.eu-south-1.amazonaws.com
```

## Installation

### Docker

The image is built `FROM scratch` and ships just the binary, for `linux/amd64`, `linux/arm64` and `linux/arm/v7`.
Copy it into your own image: Docker picks the right architecture automatically.

```dockerfile
COPY --from=ghcr.io/chialab/aws-ecr-get-login-password:latest /usr/bin/aws-ecr-get-login-password /usr/bin/aws-ecr-get-login-password
```

### Linux packages

`.deb`, `.rpm` and `.apk` packages for `amd64`, `arm64` and `armv7` are attached to each release, and install the binary in `/usr/bin`.
For example, on amd64:

```sh
# Debian, Ubuntu
curl -fsSLO https://github.com/chialab/aws-ecr-get-login-password/releases/latest/download/aws-ecr-get-login-password-linux-amd64.deb
sudo dpkg -i aws-ecr-get-login-password-linux-amd64.deb

# Fedora, RHEL
sudo rpm -i https://github.com/chialab/aws-ecr-get-login-password/releases/latest/download/aws-ecr-get-login-password-linux-amd64.rpm

# Alpine (packages aren't signed, see "Verifying artifacts" below)
curl -fsSLO https://github.com/chialab/aws-ecr-get-login-password/releases/latest/download/aws-ecr-get-login-password-linux-amd64.apk
apk add --allow-untrusted aws-ecr-get-login-password-linux-amd64.apk
```

### Pre-built binaries

Download the binary for your platform from the [latest release](https://github.com/chialab/aws-ecr-get-login-password/releases/latest).
For example, on Linux amd64:

```sh
curl -fsSL -o /usr/local/bin/aws-ecr-get-login-password \
  https://github.com/chialab/aws-ecr-get-login-password/releases/latest/download/aws-ecr-get-login-password-linux-amd64
chmod +x /usr/local/bin/aws-ecr-get-login-password
```

Each binary is also available as a `.tar.gz` archive (`.zip` on Windows), with the license and this README:

```sh
curl -fsSL https://github.com/chialab/aws-ecr-get-login-password/releases/latest/download/aws-ecr-get-login-password-linux-amd64.tar.gz \
  | tar -xz -C /usr/local/bin aws-ecr-get-login-password
```

Checksums are published alongside each release in `checksums.txt`.

### Verifying artifacts

Binaries, archives, packages and images come with [build provenance attestations](https://docs.github.com/en/actions/security-for-github-actions/using-artifact-attestations), proving they were built by this repository's release workflow:

```sh
gh attestation verify aws-ecr-get-login-password-linux-amd64 --repo chialab/aws-ecr-get-login-password
gh attestation verify oci://ghcr.io/chialab/aws-ecr-get-login-password:latest --repo chialab/aws-ecr-get-login-password
```

### From source

```sh
go install github.com/chialab/aws-ecr-get-login-password@latest
```

## Usage

```sh
aws-ecr-get-login-password | docker login --username AWS --password-stdin 000000000000.dkr.ecr.eu-south-1.amazonaws.com
```

> [!NOTE]
> Replace `000000000000.dkr.ecr.eu-south-1.amazonaws.com` with your own registry. The password works for any registry your identity has access to, and expires after 12 hours.

The output is the password alone, so it works with any tool reading it from stdin, e.g. `helm registry login --password-stdin` or `crane auth login --password-stdin`.

| Option | Description |
|---|---|
| `--version` | Print version and exit. |
| `-h`, `--help` | Print help and exit. |

### Credentials and region

There are no options for credentials or region on purpose: they are resolved by the [AWS SDK for Go](https://docs.aws.amazon.com/sdkref/latest/guide/standardized-credentials.html) as for the AWS CLI.
The most common environment variables are:

- `AWS_REGION` (or `AWS_DEFAULT_REGION`)
- `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY` and optionally `AWS_SESSION_TOKEN`
- `AWS_PROFILE`, with the shared `~/.aws/config` and `~/.aws/credentials` files
- `AWS_ROLE_ARN` and `AWS_WEB_IDENTITY_TOKEN_FILE`, e.g. with EKS IAM roles for service accounts

Without any of these, the ECS task role or the EC2 instance profile are used.
The identity needs the `ecr:GetAuthorizationToken` permission.

### Exit codes

| Situation | Exit code |
|---|---|
| Password printed | 0 |
| Password could not be obtained (e.g. missing credentials or permissions) | 1 |
| Invalid options or unexpected arguments | 2 |

Errors are printed to stderr, never to stdout.

## As a Go library

The `ecr` package returns the whole authorization data, including the expiration and the registry endpoint:

```go
token, err := ecr.GetToken(ctx, config.WithRegion("eu-south-1"))
if err != nil {
	return err
}

fmt.Println(*token.Username, *token.ProxyEndpoint, token.ExpiresAt)
```

## Security

The password grants access to your ECR registries: pipe it straight into the tool that needs it, and keep it out of logs and shell history.

To report a vulnerability, see [SECURITY.md](SECURITY.md).

## Contributing

Contributions are welcome! See [CONTRIBUTING.md](CONTRIBUTING.md) to get started.

## License

[MIT](LICENSE) © Chialab
