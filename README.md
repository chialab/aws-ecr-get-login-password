Get AWS ECR login password
==========================

This is a simple utility written in Go to build a portable binary whose sole purpose
is to mimick `aws ecr get-login-password`. This is useful in contexts like
Alpine Linux based container images, where glibc cannot be used and AWS CLI v2 doesn't work.