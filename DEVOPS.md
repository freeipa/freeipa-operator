# DevOps documentation

## Setting up the pipeline

The current github pipeline require some settings to let it run successfully.
You will need to set up the following secrets:

- `DOCKER_AUTH`: The file content with the credentials to login into the
  container image registry. This is used to create the
  `$HOME/.docker/config.json` file.
- `ENV`: This contains the specific environment variables to be used in
  your repository from github-actions. Below example of content (update
  to fit your forked repository).

  ```sh
  # Base name to use for tagging the images which reference
  # the image registry and the scope name
  IMAGE_TAG_BASE=quay.io/YOURUSER/freeipa-operator
  ```

The main repo deliveries will be stored at:

- [quay.io/freeipa/freeipa-operator](https://quay.io/repository/freeipa/freeipa-operator).
  For the controller image.
- [quay.io/freeipa/freeipa-operator-scorecard](https://quay.io/repository/freeipa/freeipa-operator-scorecard).
  For the custom scorecard that implement the functional tests for freeipa-operator
  and it is referenced by the freeipa-operator-bundle.
- [quay.io/freeipa/freeipa-operator-bundle](https://quay.io/repository/freeipa/freeipa-operator-bundle).
  For the bundle that container freeipa-operator.
- [quay.io/freeipa/freeipa-operator-catalog](https://quay.io/repository/freeipa/freeipa-operator-catalog).
  For the catalog source that reference freeipa-operator bundle.

- A lint.ignore mechanism is available. Just editing the file
  `devel/lint.ignore`, and adding the files to be ignored. The mechanism
  can be bypassed by setting `LINT_FILTER_BYPASS=1`.

- For validating K8S objects created by kustomize when launching `make lint`,
  or `./devel/lint.sh` script, we need to be logged in the
  cluster or set the variables **OC_USERNAME**, **OC_PASSWORD** and
  **OC_API_URL**.

## Checking image size

The pipeline launch dive tool to verify the layer size of the image generated
for the operator. This tool can be launched from the workstation with just:

```shell
make container-build container-save dive
```

The settings for dive tool are located at .dive-ci.yml which are the settings
used in the pipeline. The helper script `./devel/dive.sh` use the same
settings; the command `make dive` is calling to the helper script under the
hood.

## Checking kustomize manifests

A helper script has been provided for running checkov tool for the kustomize
manifests generated, so that the security can be analyzed.

For using it, we only have to run from the repository root the below:

```shell
./devel/generate-checkov-report.sh
```

The script return 0 if nothing failed scanning the manifest security, else
a number greater than 0 indicating how many kustomize directories failed.

The report is printed out in the standard output and error.
