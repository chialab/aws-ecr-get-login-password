# Both images are bumped by Dependabot; the image tags follow the Docker version.
# The Docker tag is X.Y on purpose: the weekly rebuild picks up its patch releases and base image fixes.
FROM ghcr.io/chialab/aws-ecr-get-login-password:1.1.0 AS ecr

FROM docker:29.9-cli
# /usr/local/bin as in the image built by the former chialab/docker-buildx-ecr repository
COPY --from=ecr /usr/bin/aws-ecr-get-login-password /usr/local/bin/aws-ecr-get-login-password
