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
  chcon -t container_file_t -l s0:c25,c10 /opt/freeipa/dat
  ```

## Running the controller

You need to run the controller with an extra argument for using this feature.
This is the `--default-storage {ephimeral,hostpath}` argument, or setting
the `DEFAULT_STORAGE` environment variable for using `ephimeral` or
`hostpath` values. This will allow to deploy a pod using one of those
storages, leting you to quickly deploy in a SNC for testing or developing
purposes.

- Using arguments:

  - From the workstation:

    ```shell
    make run OPERATOR_ARGS="--default-storage hostpath"
    ```

  - Deploying into the cluster:

    ```shell
    make deploy-cluster OPERATOR_ARGS="--default-storage hostpath"
    ```
