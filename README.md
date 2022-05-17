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
   ./bin/operator-sdk scorecard bundle
   ```

1. Now create a new namespace by: `oc create namespace my-freeipa`

1. As clust-admin user logged in the cluster run:

   ```sh
   make install
   ```

   > This will generate the CRD and install it into the cluster.
   > The CRD need to be installed into the cluster even if we
   > run the controller from our local workstation.

1. Run locally outside the cluster by (webhooks are disabled):

   ```sh
   make run
   ```

1. Or run inside the cluster by (first build and push the image):

   ```sh
   oc login https://my-cluster:6443
   export IMAGE_TAG_BASE=quay.io/USER_ORG/freeipa-operator
   podman login quay.io
   make docker-build
   make docker-push

   # We need cert-manager for generating the certificates for the webhooks
   make cert-manager-install
   # When the cert-manager operator is installed, run this:
   make cert-manager-self-signed-issuer-create

   # Create the scc object
   oc create -f config/rbac/scc.yaml

   # Finally deploy the operator in the cluster with:
   make deploy
   ```

1. Create `private.mk` file and update IMG_BASE variable value.

   ```sh
   cp -vf private.mk.example private.mk
   ```

   > Update `private.mk` where required

1. And create a new idm resource by:

   ```sh
   make sample-create
   ```

   > The deployment spend about 5 minutes to finish, after that
   > you will see something like the below when running:
   > `oc logs --tail=35 pod/idm-sample-main-0`

   ```raw
   [  OK  ] Finished Identity, Policy, Audit.
   ==============================================================================
   Setup complete

   Next steps:
      1. You must make sure these network ports are open:
         TCP Ports:
           * 80, 443: HTTP/HTTPS
           * 389, 636: LDAP/LDAPS
           * 88, 464: kerberos
         UDP Ports:
           * 88, 464: kerberos

      2. You can now obtain a kerberos ticket using the command: 'kinit admin'
         This ticket will allow you to use the IPA tools (e.g., ipa user-add)
         and the web user interface.
      3. Kerberos requires time synchronization between clients
         and servers for correct operation. You should consider enabling chronyd.

   Be sure to back up the CA certificates stored in /root/cacert.p12
   These files are required to create replicas. The password for these
   files is the Directory Manager password
   The ipa-server-install command was successful
   FreeIPA server does not run DNS server, skipping update-self-ip-address.
   Created symlink /etc/systemd/system/container-ipa.target.wants/ipa-server-
   update-self-ip-address.service → /usr/lib/systemd/system/ipa-server-update-
   self-ip-address.service.
   Created symlink /etc/systemd/system/container-ipa.target.wants/ipa-server-
   upgrade.service → /usr/lib/systemd/system/ipa-server-upgrade.service.
   Removed /etc/systemd/system/container-ipa.target.wants/ipa-server-configure-
   first.service.
   [  OK  ] Finished Configure IPA server upon the first start.
   FreeIPA server configured.
   ```

1. Now you should be able to reach out the web interface by:

   ```sh
   xdg-open "https://$(oc get route idm-sample -o jsonpath='{.spec.host}')"
   ```

1. Look at your objects by: `kubectl get all,idm,pvc,secrets`

1. And clean-up the cluster by:

   ```sh
   make undeploy
   oc delete -f config/rbac/scc.yaml
   ```

## Executing tests

- For the unit tests run:

  ```sh
  make test
  ```

- For the integration tests with scorecard run:

  ```sh
  # Generate bundle directory
  # bundle.Dockerfile is generated on this step
  # More information about the LABELS inside here:
  # https://github.com/operator-framework/operator-registry/blob/master/docs/design/operator-bundle.md#bundle-annotations
  # https://olm.operatorframework.io/docs/tasks/creating-operator-bundle/#contents-of-annotationsyaml-and-the-dockerfile
  make bundle
  # Running scorecard tests generated in the bundle directory by
  make scorecard-bundle
  ```

## Deploying with OLM

**Pre-requisites**:

- A proper `private.mk` file setup. (see `private.mk.example`).
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
