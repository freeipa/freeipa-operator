#!/bin/bash

##
# Generate checkov report for the kustomize directories.
##

source "devel/include/verbose.inc"


if command -v podman &>/dev/null; then CRI=podman
elif command -v docker &>/dev/null; then CRI=docker
else die "No supported container runtime was found"
fi


function analize-directory
{
    local directory="$1"
    [ -e "${directory}" ] || die "'${directory}' directory does not exist"

    echo "Scaning '${directory}' directory"
    verbose \
    ${CRI} run --rm -t \
               -v "${PWD}:${PWD}" \
               -w "${PWD}" \
               bridgecrew/checkov \
               -d "${directory}"
}

reto=0

analize-directory "./config/crd" || reto=$(( reto + 1 ))
analize-directory "./config/certmanager" || reto=$(( reto + 1 ))
analize-directory "./config/default" || reto=$(( reto + 1 ))
analize-directory "./config/manager" || reto=$(( reto + 1 ))
analize-directory "./config/prometheus" || reto=$(( reto + 1 ))
analize-directory "./config/rbac" || reto=$(( reto + 1 ))
analize-directory "./config/webhook" || reto=$(( reto + 1 ))

if [ "${reto}" -gt 0 ]; then yield "Found hints on ${reto} directories"
else yield "No hints found in directories"
fi

exit ${reto}
