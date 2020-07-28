# How to build

> EXPERIMENTAL OpenShift operator for FreeIPA

1. Install `operator-sdk`.

1. Check out repository under `$GOPATH/src/`.  (`GOPATH` defaults to
   `$HOME/go`.)

1. `operator-sdk build quay.io/username/freeipa-operator:v0.0.1`

## Development tools

For making life easier, a script has been included to deploy a complete
development environment. Actually it is verified for Fedora 32, more Linux
distributions are welcome and hoping this get more and more completed.

This script try to install:

- A set of needed packages for creating ansible or golang operators.
- OpenShift Client, so that we can communicata with an OpenShift cluster.
- Operator-SDK which make life easier creating operators.
- (optional) Visual Studio Code.
- (optional) CodeReady Containers.

Below we see how to use the script:

```shell
./devel/install-local-tools.sh
```
