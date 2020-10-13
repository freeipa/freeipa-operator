# Init Containers

An init container is kubernete pattern included in structural patterns.
This pattern help us to decouple tasks related with initialization of the Pod
as it could be initialize a data volume, upgrade the data model in the data
volume, prepare and collect information feed configuration, synchronize with
other services running in the cluster before go on, shrink the privileged
actions allowing the main container run as no provileged or a reduced set
of privileges.

For all the above, it deserves to be investigated and create different proof
of concepts that could provide us a better understanding so we can use this
pattern as a tool to set up properly the container environment before launch
the Freeipa container.

The content of this file is break down as follow:

- Proof of Concepts: Several different proof of concepts which provide us
  a better understunding about how to play with the different objects.
  - **Hello world**: A minimal structure for using initContainer pattern.
  - **Read uid/gid ranges**: This read from the namespace and inject them into
    a configMap object. This will need a ServiceAccount with permissions to
    write to the ConfigMap.
  - **Read Pod IP**: In a similar way to the above, the POD IP is read from
    by using the downward API and stored into the configMap object for
    being used later.
  - **Initialize machine-id**: Investigating about systemd container interface
    it was seen a way to initialize the machine-id for systemd; this
    proof of concept just inplement the mechanism initialising the value
    into the configMap so that it can be used later to inject the value as
    a environment variable.
  - **Print out environment**: this is helpful to know the about the initial state
    for /etc/hosts, /etc/resolv.conf, capabilities and the filesystem mounted.
    - **DNS Config**: Play with the DNS config used to know more about the Pod
      configuration.
    - **Capabilities**: Initial study about the capabilities which are needed
      for initialise the data volume
  - **Initialize the volume data**: Run ipa-install-server for initialising the
    data volume.
    - **PENDING - Map uid/gid ranges**: This could be useful to prepare the UID/GID
      map to be used with all the different services running in FreeIpa.
    - **Add host alias**: Intent to remove some unnecesary name resolution
      externally.
- **PENDING - Application definition**: Use all the above for creating the
  application objects. This will be addressed in a different investigation.
  Just reference for this "Chapter 11: Stateful Service".
- **Shrink capabilities**: This will need a deeper study and investigation as
  it could increase the scope of this task more than I want. A different
  PR will be addressed for this task.

The logs are printed out in the console, so they can be read using the
command below: `oc logs pod/my-pod -c container-name`.

For making life easier, a Makefile is provided to automate some actions and
speed up playing with the different PoCs created as part of this investigation.

## Proof of Concepts

### Hello world

Minimal example running hello world container in an initContainer.

Play with it:

```shell
export APP=poc-00
make app-delete app-deploy get-info
```

See: [poc-00.yaml](poc-00.yaml)

### Read the namespace UID/GID ranges and inject in a ConfigMap

Here we can see a PoC which show how to use the API inside an init container to
retrieve the necesary information and modify a configmap to store there the
information which will be injected in the main container for using it.

In this case the information is just displayed, but once the information is
stored in the ConfigMap, it can be injected in a container and use it as it is
needed.

Play with it:

```shell
export APP=poc-01
make app-deploy
make get-info
make app-delete
```

See: [poc-01.yaml](poc-01.yaml).

### Read POD IP and inject in a ConfigMap

This proof of concept just show how to inject the POD IP inside the containers
and initContainers, so it can be used wherever it is needed.

Play with it:

```shell
export APP=poc-03
make app-deploy
make get-info
make app-delete
```

See: [poc-03.yaml](poc-03.yaml).

### Initialize /etc/machine-id

This proof of concept initialize the machine-id in the ConfigMap so that it can
be injected into the containers as an environment variable when running the
main container.

Play with it:

```shell
export APP=poc-04
make app-deploy
make get-info
make app-delete
```

See: [poc-04.yaml](poc-04.yaml).

### List features for different scenarios

- Filesystem mounted for "privileged: true" and not privileged containers.
- capabilities for different scc with no privileged attribute, and
  privileged: true: anyuid, privileged, restricted.
- State of /etc/hosts when using and not using `hostAliases`.
- State of /etc/resolv.conf when using and not using `dnsConfig`.

You need to build the image for this pod by:

```shell
podman build -t quay.io/avisied0/freeipa-openshift-container:print-features \
             -f Dockerfile.print-features .
podman push quay.io/avisied0/freeipa-openshift-container:print-features
```

And modify the container image used in the pod, by yours, if you change the
image name.

Finally, play with it by:

```shell
export APP=poc-05-f
make app-deploy
make get-info
make app-delete
```

See: [poc-05-f.yaml](poc-05-f.yaml).

#### /etc/resolv.conf

Default content:

```raw
search avisiedo-init-container.svc.cluster.local svc.cluster.local \
cluster.local permanent.idmocp.lab.eng.rdu2.redhat.com
nameserver 172.30.0.10
options ndots:5
```

Which match the IP assigned by the service. We can check it by
`oc get svc --namespace=openshift-dns`:

```raw
NAME          TYPE        CLUSTER-IP    EXTERNAL-IP   PORT(S)                  AGE
dns-default   ClusterIP   172.30.0.10   <none>        53/UDP,53/TCP,9153/TCP   18d
```

When we use this dnsConfig:

```yaml
  dnsConfig:
    nameservers:
      - 127.0.0.1
```

We get the following `/etc/resolv.conf`:

```raw
search avisiedo-init-container.svc.cluster.local svc.cluster.local \
cluster.local permanent.idmocp.lab.eng.rdu2.redhat.com
nameserver 172.30.0.10
nameserver 127.0.0.1
options ndots:5
```

#### Capabilities

In short, we need the below capabilities: CAP_SYS_ADMIN, CAP_MKNOD and
CAP_SYS_RESOURCE to run systemd, and by extension the ipa-server-install
in the init container.

Default set of capabilities for anyuid:

- CAP_CHOWN
- CAP_DAC_OVERRIDE
- CAP_FOWNER
- CAP_FSETID
- CAP_KILL
- CAP_SETGID
- CAP_SETUID
- CAP_SETPCAP
- CAP_NET_BIND_SERVICE
- CAP_NET_RAW
- CAP_SYS_CHROOT

Additional capabilities for running systemd for launching ipa-server-install:

- CAP_SYS_ADMIN: This is the big one to try to remove.
- CAP_SYS_RESOURCE: Needed because some pctl system calls are made by systemd.
- CAP_MKNOD: Needed because systemd could try to create some nodes, for terminals.
  Maybe this one could be removed.

Some article about CAP_SYS_ADMIN is interesting here, that could help to
remove CAP_SYS_ADMIN capability:

- [LXC containers without CAP_SYS_ADMIN under Debian Jessie](https://blog.iwakd.de/lxc-cap_sys_admin-jessie).

At the moment this PR was created, it was not working, and a new proof of
concept was addressed to focus only into the set of needed capabilities
and the SCC profile to be assigned to the service account, which could
allow to run the initContainer without `privileged: true`. So that a new
PR will be addressed for this.

Another PR will be addressed for trying to remove CAP_SYS_ADMIN based on
the previous article.

### Initialize the volume data by freeipa-server-install

This proof of concept provides a pod with a sequence of initContainers
initializing the data volume. The volume used is ephimeral so it will be
populated with every delete and redeployment.

This set of objects have been configured to provide the maximum levels of
traces so that it can be used to detect changes to be made or improvements
on the initilisation process.

This PoC need container image to be built and published, you can do that
by:

```shell
export DOCKER_IMAGE=quay.io/username/freeipa-server:dev-test
podman login quay.io
make container-build container-push
```

for docker, do the below:

```shell
export DOCKER_IMAGE=docker.io/username/freeipa-server:dev-test
docker login quay.io
make container-build
make container-push
```

> - You need a docker hub account or a quay.io account before launch the above,
>   and you need to be logged in.
>
> - You need to change the image: attribute into the Pod definition to point to
>   **DOCKER_IMAGE**. This should be change into **init-volume** and **main**
>   containers.

Finally, play with it by:

```shell
export APP=poc-05-a
make app-deploy
make get-info
make app-delete
```

see: [poc-05-a.yaml](poc-05-a.yaml).

> It is needed to shrink the privileges assigned to the container; as said
> above, it will be addressed in a new PR.

#### PENDING Use the UID/GID base values to create the uid/gid maps

Extend the previous PoC to generate /etc/passwd and /etc/group files mapping
the user ids to match the namespace ranges, and update them in the ConfifMap.

This content is injected later in the main container to map properly the
different userid and groupid schema. This could be used when running the
freeipa-server-install process, so all the uid/gid exist in the system.

see: [poc-05-b.yaml](poc-05-b.yaml).

#### Adding host alias to the Pod

This proof of concept add an alias to the /etc/host in all the pod containers.
This could be useful when resolving the full qualified name inside the Pod.

When we apply this configuration:

```yaml
  hostAliases:
    - ip: "127.0.0.1"
      hostnames:
        - poc-05-f.apps.permanent.idmocp.idm.lab.bos.redhat.com
        - apps.permanent.idmocp.idm.lab.bos.redhat.com
```

We get:

```raw
# Kubernetes-managed hosts file.
127.0.0.1       localhost
::1     localhost ip6-localhost ip6-loopback
fe00::0 ip6-localnet
fe00::0 ip6-mcastprefix
fe00::1 ip6-allnodes
fe00::2 ip6-allrouters
10.143.1.182    poc-05-f-default-with-host-aliases

# Entries added by HostAliases.
127.0.0.1       poc-05-f.apps.permanent.idmocp.idm.lab.bos.redhat.com   apps.permanent.idmocp.idm.lab.bos.redhat.com
```

```shell
export APP=poc-05-f
make app-deploy
make get-info
make app-delete
```

see: [poc-05-f.yaml](poc-05-f.yaml).

## References

- Systemd:
  - [Understanding and Using Systemd](https://www.linux.com/training-tutorials/understanding-and-using-systemd/).
  - [Rethinking PID 1](http://0pointer.de/blog/projects/systemd.html).
  - [systemd - The Container Interface](https://systemd.io/CONTAINER_INTERFACE/).
  - [systemd - Environment variables](https://systemd.io/ENVIRONMENT/).
  - [MachineId](https://wiki.debian.org/MachineId).
  - [LXC containers without CAP_SYS_ADMIN under Debian Jessie](https://blog.iwakd.de/lxc-cap_sys_admin-jessie).

- Kubernetes:
  - [Configuring Permissions in Kubernetes with RBAC](https://medium.com/containerum/configuring-permissions-in-kubernetes-with-rbac-a456a9717d5d).
  - [Using RBAC Authorization](https://v1-16.docs.kubernetes.io/docs/reference/access-authn-authz/rbac/).
  - [Access Cluster API](https://kubernetes.io/docs/tasks/administer-cluster/access-cluster-api/#go-client).
  - [Using service accounts in applications](https://docs.openshift.com/container-platform/4.4/authentication/using-service-accounts-in-applications.html).
  - [DNS for Services and Pods](https://kubernetes.io/docs/concepts/services-networking/dns-pod-service/).
  - [Access Clusters Using the Kubernetes API - Directly accessing the REST API](https://kubernetes.io/docs/tasks/administer-cluster/access-cluster-api/#directly-accessing-the-rest-api-1).
  - [Define container environment variable using configmap data](https://kubernetes.io/docs/tasks/configure-pod-container/configure-pod-configmap/#define-container-environment-variables-using-configmap-data).
  - [Adding additional entries with hostaliases](https://kubernetes.io/docs/concepts/services-networking/add-entries-to-pod-etc-hosts-with-host-aliases/#adding-additional-entries-with-hostaliases).
  - [Services networking - DNS pod service](https://kubernetes.io/docs/concepts/services-networking/dns-pod-service/).
  - [Service - External name](https://kubernetes.io/docs/concepts/services-networking/service/#externalname).
  - [Service - External IPs](https://kubernetes.io/docs/concepts/services-networking/service/#external-ips).

- OpenShift:
  - [Developer CLI operations](https://docs.openshift.com/container-platform/4.4/cli_reference/openshift_cli/developer-cli-commands.html).
  - [A Guide to OpenShift and UIDs](https://www.openshift.com/blog/a-guide-to-openshift-and-uids).
