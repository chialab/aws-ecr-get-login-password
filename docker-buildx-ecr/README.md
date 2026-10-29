# docker-buildx-ecr

[![Docker image](https://github.com/chialab/aws-ecr-get-login-password/actions/workflows/docker-buildx-ecr.yml/badge.svg)](https://github.com/chialab/aws-ecr-get-login-password/actions/workflows/docker-buildx-ecr.yml)

**The official Docker CLI image, plus `aws-ecr-get-login-password`.**

`ghcr.io/chialab/docker-buildx-ecr` is [`docker:<version>-cli`](https://hub.docker.com/_/docker) with [`aws-ecr-get-login-password`](../README.md) in `/usr/local/bin`.
It is meant for CI/CD pipelines that build images and push them to Amazon ECR, but works with any registry.

Images are published for `linux/amd64`, `linux/arm64` and `linux/arm/v7`.
Tags follow the Docker version: `29.9.2`, `29.9`, `29`, and `latest` (also published as `main`, for compatibility).
The image is rebuilt weekly, to pick up Docker patch releases and base image fixes.

> [!NOTE]
> This image used to live in [chialab/docker-buildx-ecr](https://github.com/chialab/docker-buildx-ecr), now archived. The image name and tags are unchanged.

## Usage

### Log in to ECR

The identity running the command needs the `ecr:GetAuthorizationToken` permission; see [Credentials and region](../README.md#credentials-and-region).

```sh
aws-ecr-get-login-password | docker login --username AWS --password-stdin "${ECR_REGISTRY}"
```

### Build a multi-architecture image

```sh
docker run --privileged --rm tonistiigi/binfmt --install all
docker context create multi-arch
docker --context multi-arch buildx create --use
docker --context multi-arch buildx build --platform linux/amd64,linux/arm64 --tag "${IMAGE_NAME}:${IMAGE_TAG}" --push .
```

## Verifying the image

Images come with [build provenance attestations](https://docs.github.com/en/actions/security-for-github-actions/using-artifact-attestations), produced by this repository's workflow:

```sh
gh attestation verify oci://ghcr.io/chialab/docker-buildx-ecr:latest --repo chialab/aws-ecr-get-login-password
```
