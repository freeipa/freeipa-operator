# freeipa-operator

Experimental freeipa-operator for Freeipa.

## Quick Start

> It requires golang 1.16; if your system is providing a lower
> version, consider to install [gvm](https://github.com/moovweb/gvm#installing).
> for using different golang versions.

1. Clone the repository by:

   ```sh
   git clone https://github.com/freeipa/freeipa-operator.git
   cd freeipa-operator
   ```

1. Install the necessary tools by:

   ```sh
   ./devel/install-local-tools.sh
   ```

1. Build by:

   ```sh
   make build
   ```

1. Launch tests by:

   ```sh
   make test
   operator-sdk scorecard bundle
   ```

1. Now create a new namespace by: `kubectl create namespace my-freeipa`

1. Run locally outside the cluster by (webhooks are disabled):

   ```sh
   make run
   ```

1. Or run inside the cluster by (first build and push the image):

   ```sh
   kubectl login https://my-cluster:6443
   export IMAGE_TAG_BASE=quay.io/USER_ORG/freeipa-operator
   podman login quay.io
   make docker-build
   make docker-push

   # We need cert-manager for generating the certificates for the webhooks
   make cert-manager-install
   # When the cert-manager operator is installed, run this:
   make cert-manager-self-signed-issuer-create

   # Finally deploy the operator in the cluster with:
   make deploy
   ```

1. And create a new idm resource by:

   ```sh
   cat > private.mk <<EOF
   IDM_ADMIN_PASSWORD=myPassword124
   IDM_DM_PASSWORD=DMmyPassword124
   SAMPLE=config/samples/ephemeral-storage
   EOF
   make sample-create
   ```

   > You can check more samples at `config/samples` directory.

1. Look at your objects by: `kubectl get all,idm,pvc,secrets`

1. And clean-up the cluster by:

   ```sh
   make undeploy
   ```

## Executing tests

- For the unit tests run:

  ```sh
  make test
  ```

- For the integration tests with scorecard run:

  ```sh
  # Generate bundle directory
  make bundle
  # Running scorecard tests generated in the bundle directory by
  make scorecard-bundle
  ```

## Deploying with OLM

**Pre-requisites**:

- A proper `private.mk` file setup. (see `private.mk.example`).
- A namespace selected (eg. `oc new-project ipa`).
- The freeipa SecurityContextConstraint created (`oc create -f config/rbac/scc.yaml`).

**Steps**:

1. Create a namespace:

   ```sh
   oc new-project ipa
   ```

1. Build and publish container images:

   ```sh
   make docker-build docker-push \
        bundle bundle-build bundle-push \
        catalog-build catalog-push
   ```

1. Install operator with OLM in the current namespace by:

   ```sh
   make bundle-install
   ```

1. Create a sample idm resource:

   ```sh
   oc create -f config/samples/persistent-storage.yaml
   ```

1. Delete the custom resource created:

   ```sh
   oc delete -f config/samples/persistent-storage.yaml
   ```

   > TODO You will need to delete the PVC by hand if a new
   > IDM resource have to be created with different options.

1. Cleanup the operator from the cluster:

   ```sh
   make bundle-uninstall
   ```

1. Remove the namespace:

   ```sh
   oc delete namespace ipa
   ```

----

See also: [Operator SDK 1.0.0 - Quick Start](https://sdk.operatorframework.io/docs/building-operators/golang/quickstart/).
