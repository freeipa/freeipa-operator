# DevOps documentation

## Setting up the pipeline

The current github pipeline require some settings to let it run successfully.
You will need to set up the following secrets:

- DOCKER_AUTH: The file content with the credentials to login into the
  container image registry. This is used to create the
  `$HOME/.docker/config.json` file.
- IMG_BASE: The base name for your image. This could be something like:
  - `quay.io/freeipa/freeipa-operator`.
  - `docker.io/freeipa/freeipa-operator`.
  - `quay.io/avisied0/my-freeipa-operator`.
  This provide flexibility, and allow that forked repositories could made
  deliveries on their own image registries, or different repository.

The deliveries will be stored at:
[quay.io/freeipa/freeipa-operator](https://quay.io/repository/freeipa/freeipa-operator).

- A lint.ignore mechanism is available. Just editing the file
  `devel/lint.ignore`, and adding the files to be ignored. The mechanism
  can be bypassed by setting `LINT_FILTER_BYPASS=1`.
