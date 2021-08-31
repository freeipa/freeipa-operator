# hostpath storage

This storage MUST NOT used in production. This is provided just for developing
propose for making quicker and easier to check the operator behavior.

## Preparing the node

This configuration has been thought to make developer life easier so
that this can be tested from a SNC (Single Node Cluster) as
CodeReadyContainers or a SNC deployed with kcli.

The usage of this configuration require a manual steps into the node
to prepare it.

- Retrieve the namespace information for freeipa for msc, we will need
  it for setting the right context to the hostPath.

  ```shell
  oc describe namespace/freeipa
  ```

- Login into the node:

  ```shell
  oc get nodes
  oc debug nodes/NODE
  ```

- Create the internal path:

  ```shell
  mkdir /opt/freeipa/data
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
oc create -f - | cat <<EOF
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
