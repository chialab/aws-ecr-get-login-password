# Built by GoReleaser (dockers_v2): binaries are pre-built in the context under $TARGETPLATFORM/.
FROM scratch
ARG TARGETPLATFORM
COPY $TARGETPLATFORM/aws-ecr-get-login-password /usr/bin/aws-ecr-get-login-password
ENTRYPOINT ["/usr/bin/aws-ecr-get-login-password"]
