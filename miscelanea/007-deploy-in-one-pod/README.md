# Deploy in one Pod

This proof of concept try to deploy a unique Pod with all the services
running, but not keeping the expected security standards. This can be
used for future scenarios which needs a deployment to review the
configurations and how to make work scenarios with multiple clouds
collaborating.

As I have said, this is quite far from the goals we want to get, starting
because this is no secure by default; it must be seen just as it is,
a proof of concept, that is all.

Unwished features:

- Freeipa is deployed using epimeral volume.
- The pod is using SYS_ADMIN capability.
- The pod is using NET_BIND capability.
- The pod is using host paths.
- Running systemd as root inside the container means the uid is mapped
  to the root user in the host (critical privilege scalation security vector),
  as it was disclosed by @ftweedal investigations (at least in OpenShift 4.4
  that is the current permanent cluster).

Features:

- A single pod deployment.
- The web UI is accessible.
- The services are accessible inside the cluster.

## Quick Start

Preconditions:

- Review the image to be build, that it is referenced properly, and the
  registry you have access and you are logged in.
- Review you are loggid in the cluster and a namespace has been created.

Afterward just run the below:

```shell
make container-build container-push app-deploy
```
