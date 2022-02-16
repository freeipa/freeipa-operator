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
   make
   ```

1. Launch tests by:

   ```sh
   make test
   ```

1. Run locally outside the cluster by:

   ```sh
   make run
   ```

1. Now create a new namespace by: `kubectl create namespace my-freeipa`

1. Or run inside the cluster by (first build and push the image):

   ```sh
   kubectl login https://my-cluster:6443
   make container-build IMG=quay.io/USER_ORG/freeipa-operator:dev-test
   podman login quay.io
   make container-push IMG=quay.io/USER_ORG/freeipa-operator:dev-test

   # We need cert-manager for generating the certificates for the webhooks
   make -f mk/cert-manager cert-manager-install
   # When the cert-manager operator is installed, run this:
   make -f mk/cert-manager.mk cert-manager-self-signed-issuer-create

   # Finally deploy the operator in the cluster with:
   make deploy-cluster IMG=quay.io/USER_ORG/freeipa-operator:dev-test
   ```

1. And create a new idm resource by:

   ```sh
   IDM_ADMIN_PASSWORD=myPassword124 \
   IDM_DM_PASSWORD=myPassword125 \
   SAMPLE=ephemeral-storage \
   make sample-create
   ```

   > You can check more samples at `config/samples` directory.

1. Look at your objects by: `kubectl get all,idm,pvc,secrets`

1. And clean-up the cluster by:

   ```sh
   SAMPLE=ephemeral-storage make sample-delete
   make undeploy-cluster
   ```

<!-- TODO When the read of ingresDomain is implemented, remove the
          block below. -->

> When using CodeReadyContainers, you will need to add the entry
> `192.168.130.11   NAMESPACE.apps.crc.testing` to your `/etc/hosts` file
> or it will not work as expected; in a real cluster it works
> properly because the ingressDomain use to match `*.apps.<basedomain>`.
>
> Now it is known that the ingressDomain information can be retrieved more
> accurate from a cluster resource and it will be corrected in a future PR.

See also: [Operator SDK 1.0.0 - Quick Start](https://sdk.operatorframework.io/docs/building-operators/golang/quickstart/).
