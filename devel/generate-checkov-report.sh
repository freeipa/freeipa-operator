#!/bin/bash

##
# Generate checkov report for the kustomize directories.
##

source "devel/include/verbose.inc"


if command -v podman &>/dev/null; then CRI=podman
elif command -v docker &>/dev/null; then CRI=docker
else die "No supported container runtime was found"
fi

TTY_OPTS=
tty &>/dev/null && TTY_OPTS="-ti"


function analize-directory
{
    local directory="$1"
    [ -e "${directory}" ] || die "'${directory}' directory does not exist"

    echo "Scaning '${directory}' directory"
    verbose \
    ${CRI} run --rm ${TTY_OPTS} \
               -v "${PWD}:${PWD}" \
               -w "${PWD}" \
               bridgecrew/checkov \
               --framework kubernetes \
               -d "${directory}"
}

reto=0

analize-directory "./config/" || reto=$(( reto + 1 ))

if [ "${reto}" -gt 0 ]; then yield "Found hints on ${reto} directories"
else yield "No hints found in directories"
fi

exit ${reto}
