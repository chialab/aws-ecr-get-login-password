# Agent instructions

`aws-ecr-get-login-password` is a single-binary Go CLI: it prints an Amazon ECR login password to stdout, like `aws ecr get-login-password`, for systems where the AWS CLI can't run (e.g. Alpine/musl or `scratch` images).
User-facing docs live in [README.md](README.md); contribution rules in [CONTRIBUTING.md](CONTRIBUTING.md).

## Layout

- `main.go`: the CLI (`--version`, `-h`, exit codes); `run` takes args, writers and the token source, so it is testable without AWS.
- `ecr/token.go`: importable package. Public API: `GetToken`, `NewToken`, `AuthorizationToken`, `ErrInvalidToken`, `ErrUnexpectedResponse`. The `ecrTokenClient` interface exists to mock the ECR client in tests.
- `main_test.go`, `ecr/token_internal_test.go`: table-driven tests, all `t.Parallel()`, never calling AWS.
- `.goreleaser.yaml` + `Dockerfile`: release of plain binaries and tar.gz/zip archives (linux/darwin/windows, amd64/arm64, linux/arm/v7), Linux packages (deb/rpm/apk, attached to the release only) and the multi-arch `FROM scratch` image on GHCR.
- `.github/workflows/`: `ci.yml` (test matrix with `-race` and Codecov, golangci-lint, govulncheck, zizmor), `release.yml` (GoReleaser on `v*` tags, with attestations). CodeQL runs through the repository's default setup, not a workflow.

## Commands

Run all of these before declaring a change done:

```sh
go test -race ./...
golangci-lint fmt && golangci-lint run
go tool govulncheck ./...
go mod tidy -diff
zizmor --persona=pedantic .         # when touching .github/
goreleaser check                    # when touching .goreleaser.yaml
```

golangci-lint may be missing locally (e.g. in the dev container): CI runs it, so say so instead of installing it ad hoc.

## Behavior contract

These are deliberate decisions; don't change them without being asked:

- stdout carries the password and nothing else, followed by a newline: it is piped into `docker login --password-stdin`.
- Errors go to stderr as `aws-ecr-get-login-password: <error>`; exit codes are 0 on success, 1 if the token can't be obtained, 2 for invalid options or unexpected arguments.
- The only options are `--version` and `-h`/`--help`. Region, credentials and everything else come from the AWS SDK default chain (environment, shared config files, web identity, container and instance roles), as for the AWS CLI. Don't add flags that duplicate that.
- Release asset names (`aws-ecr-get-login-password-<os>-<arch>`, no version; plain binary, plus `.tar.gz`/`.zip` archives and packages) are a public interface: install scripts download `releases/latest/download/<name>`.
- The `ecr` package's public API is used by other projects: keep it backward compatible.
- Root CAs are embedded (`golang.org/x/crypto/x509roots/fallback`) and used only when the system has none, so the binary works in `scratch` images.

## Conventions

- **Linting is maximal**: `.golangci.yml` enables every linter except a few documented ones. Fix findings rather than disabling linters. For a genuine false positive use `//nolint:<linter> // <reason>`: the reason is mandatory and `nolintlint` checks it.
- **No new dependencies** unless asked: the `depguard` allowlist in `.golangci.yml` only permits the stdlib, `github.com/aws/aws-sdk-go-v2` and `golang.org/x/crypto/x509roots/fallback`.
- Errors are sentinel values or wrapped with `%w`, tested with `errors.Is`.
- Plain `err :=` assignments, not inline `if err := ...` or `if ... else` chains.
- GitHub Actions are referenced by **major tag, not SHA** (on purpose, see `.github/zizmor.yml`); keep least-privilege `permissions` and `persist-credentials: false`.
- Keep README and this file in sync when options or behavior change.
- User-facing text, code, comments and commits are in English.
