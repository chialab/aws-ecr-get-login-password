# Contributing

Thanks for taking the time to contribute! 🎉
Bug reports, ideas, docs fixes and pull requests are all welcome, big or small.

By participating you agree to abide by our [Code of Conduct](CODE_OF_CONDUCT.md).

## Reporting bugs and requesting features

Before opening an issue, please search the [existing ones](https://github.com/chialab/aws-ecr-get-login-password/issues): someone may have already reported it.
When you open a new one, the issue templates will ask for what we need to help you quickly.

**Never paste passwords or AWS credentials in an issue**, and feel free to redact your account ID.

Found a security issue? Please don't open a public issue: follow [SECURITY.md](SECURITY.md) instead.

## Development setup

You need the Go version declared in [`go.mod`](go.mod).
Alternatively, open the repository in the provided [dev container](.devcontainer/devcontainer.json) (VS Code, GitHub Codespaces, …) and you're ready to go.

Linters used by CI are optional locally, since CI runs them anyway.
To run them yourself, install them with the method you prefer:

- [golangci-lint](https://golangci-lint.run/docs/welcome/install/local/)
- [zizmor](https://docs.zizmor.sh/installation/)

## Before opening a pull request

Please make sure these pass:

```sh
go test -race ./...
golangci-lint fmt && golangci-lint run
go tool govulncheck ./...
zizmor --persona=pedantic .   # only if you touched .github/
```

A few guidelines:

- Keep pull requests focused: one change per PR is easier to review.
- Add or update tests for any behavior change. Tests are table-driven, run in parallel, and never call AWS.
- Linters run with (almost) everything enabled. If a finding is a false positive, silence it with a `//nolint:<linter> // <reason>` comment: the reason is mandatory.
- New dependencies must be allowed explicitly in the `depguard` rules of [`.golangci.yml`](.golangci.yml): please discuss them in an issue first, as the goal is a small, self-contained binary.
- Update the [README](README.md) when you change options or behavior.

## Releases

Maintainers release by pushing a signed `vX.Y.Z` tag. GoReleaser then publishes the binaries and Linux packages to GitHub Releases and the multi-arch image to GHCR.

Thanks again for helping make `aws-ecr-get-login-password` better! 💚
