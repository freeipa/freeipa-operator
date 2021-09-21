# freeipa-operator

Experimental freeipa-operator for Freeipa.

## Quick Start

> It requires golang 1.16; if your system is providing a lower
> version, consider to install [gvm](https://github.com/moovweb/gvm#installing).
> for using different golang versions.

1. Clone the repository by:

   ```shell
   git clone https://github.com/freeipa/freeipa-operator.git
   cd freeipa-operator
   ```

1. Install the necessary tools by:

   ```shell
   ./devel/install-local-tools.sh
   ```

1. Build by:

   ```shell
   make
   ```

1. Launch tests by:

   ```shell
   make test
   ```

1. Run locally outside the cluster by:

   ```shell
   make run
   ```

1. Or run inside the cluster by (first build and push the image):

   ```shell
   kubectl login https://my-cluster:6443
   make container-build IMG=quay.io/freeipa/freeipa-operator:dev-test
   podman login quay.io
   make container-push IMG=quay.io/freeipa/freeipa-operator:dev-test
   make deploy-cluster IMG=quay.io/freeipa/freeipa-operator:dev-test
   ```

   Now it support deploy in kind bu just typing:

   ```shell
   make deploy-kind IMG=quay.io/freeipa/freeipa-operator:dev-test
   ```

1. Now create a new namespace by: `kubectl create namespace my-freeipa`

1. And create a new idm resource by: `make recreate-sample-idm`

1. Look at your objects by: `kubectl get all,idm`

1. And clean-up the cluster by:

   ```shell
   kubectl delete -f config/samples/freeipa_v1alpha1_freeipa.yaml
   kustomize build config/default | kubectl delete -f -
   ```

See also: [Operator SDK 1.0.0 - Quick Start](https://sdk.operatorframework.io/docs/building-operators/golang/quickstart/).
