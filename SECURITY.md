# Security Policy

## Supported versions

Only the latest release receives security fixes.

## Reporting a vulnerability

Please **do not** open a public issue for security problems.
Instead, use GitHub's [private vulnerability reporting](https://github.com/chialab/aws-ecr-get-login-password/security/advisories/new), or write to [security-notices@chialab.it](mailto:security-notices@chialab.it).

Include as much detail as you can: affected version, steps to reproduce, and impact.
We'll acknowledge your report as soon as possible and keep you posted on the fix.

## Keep your credentials secret

The output of `aws-ecr-get-login-password` is a password that grants access to your ECR registries for 12 hours.
Pipe it straight into `docker login --password-stdin` (or similar) instead of storing it, and never paste it, or your AWS credentials, in issues or logs.
