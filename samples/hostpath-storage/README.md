# hostpath storage

This storage MUST NOT used in production. This is provided just for developing
propose for making quicker and easier to check the operator behavior.

## Preparing the node

This configuration has been thought to make developer life easier so
that this can be tested from a SNC (Single Node Cluster) as
CodeReadyContainers or a SNC deployed with kcli.

The usage of this configuration require a manual steps into the node
to prepare it.

- Login into the node:

  ```shell
  oc get nodes
  oc debug nodes/NODE
  ```

- Create the internal path:

  ```shell
  mkdir -p /opt/freeipa/data
  ```

- Change selinux context for the directory so that it can be accessed
  from the containers.

  ```shell
  chcon -t container_file_t /opt/freeipa/data
  ```

## Create the PersistentVolume

For using a hostPath persistent volume, you have to create it
before deploy the workload using the operator; for that you will
need to:

```shell
cat <<EOF | oc create -f -
---
apiVersion: v1
kind: PersistentVolume
metadata:
  name: freeipa
  labels:
    app: freeipa
spec:
  capacity:
    storage: 10Gi
  accessModes:
  - ReadWriteOnce
  - ReadWriteMany
  - ReadOnlyMany
  persistentVolumeReclaimPolicy: Recycle
  hostPath:
    path: /opt/freeipa/data
EOF
```

> `hostPath` attribute should match the local directory you created
> previously into the cluster node.

## Running the operator

We have two options here, running the controller locally or
running the controller from a cluster. For developing propose
the first option will fit better when debugging from the
workstation; By the way, we will need to be logged in a cluster
with enough privileges.

```shell
oc login -u username API_URL
```

### Locally

Just run:

```shell
# Install CRDS into the cluster by
make install-crds
# Create RBAC resources
kustomize build config/rbac | oc create -f -
# Run the controller
make run
```

### Cluster

Just run:

```shell
make deploy-cluster
```

## Creating the sample idm resource

Finally for using this sample we only have to run:

```shell
SAMPLE=hostpath-storage make sample-create
```
