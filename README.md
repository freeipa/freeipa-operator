# freeipa-operator

Experimental freeipa-operator for Freeipa.

## Quick Start

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
   make install
   make run ENABLE_WEBHOOKS=false
   ```

1. Or run inside the cluster by (first build and push the image):

   ```shell
   kubectl login https://my-cluster:6443
   make container-build IMG=quay.io/freeipa/freeipa-operator:dev-test
   podman login quay.io
   make container-push IMG=quay.io/freeipa/freeipa-operator:dev-test
   make deploy IMG=quay.io/freeipa/freeipa-operator:dev-test
   ```

1. And clean-up the cluster by:

   ```shell
   kubectl delete -f config/samples/freeipa_v1alpha1_freeipa.yaml
   kubectl delete deployments,service -l control-plane=controller-manager
   kubectl delete role,rolebinding --all
   kustomize build config/default | kubectl delete -f -
   ```

See also: [Operator SDK 1.0.0 - Quick Start](https://sdk.operatorframework.io/docs/building-operators/golang/quickstart/).
