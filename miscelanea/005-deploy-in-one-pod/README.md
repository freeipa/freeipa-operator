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

- Freeipa is deployed using ephimeral volume.
- The pod is using SYS_ADMIN capability.
- The pod is using NET_BIND capability.
- The pod is using host paths.
- Running systemd as root inside the container means the uid is mapped
  to the root user in the host (critical privilege scalation security vector),
  as it was disclosed by @ftweedal investigations
  ([link](https://frasertweedale.github.io/blog-redhat/posts/2020-11-05-openshift-user-namespace.html)).

Features:

- A single pod deployment.
- The web UI is exposed and accesible.
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

## Notes

- The certificate subject need to be different between deployments or after
  the first redeployment we will get a `SEC_ERROR_REUSED_ISSUER_AND_SERIAL`
  error code in Firefox; this is automated but it is something to keep in
  mind when moving to the operator.
- It is only exposed the web interface, more investigation is needed for the
  intracluster scenarios, but this prototype will allow an easier way of creating
  proof of concepts for that scenario.
