#!/bin/bash
[ "$DOCKER_IMAGE" == "" ] && echo "ERROR: 'DOCKER_IMAGE' is not defined" && exit 1
[ ! -e data ] || sudo rm -rf data
mkdir data
container_id="$( podman run --rm \
           -d \
           --tty \
           --sysctl net.ipv6.conf.all.disable_ipv6=0 \
           --hostname ipa.example.test \
           --volume "$PWD/data:/data:z" \
           -e INIT_WRAPPER=1 \
           -e container=podman \
           -e container_uuid="$( < /etc/machine-id )" \
           -e PASSWORD=Administrator \
           -e DEBUG_TRACE=2 \
           --entrypoint "/usr/sbin/init" \
           quay.io/avisied0/freeipa-openshift-container:dev-test exit-on-finished \
           -U \
           --realm EXAMPLE.TEST \
           --no-ntp --no-sshd --no-ssh --verbose "$@"
)"

trap 'echo Stopping container ${container_id}; podman stop ${container_id}' SIGINT
podman logs -f "${container_id}"
