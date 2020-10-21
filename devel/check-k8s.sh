#!/bin/bash

# shellcheck disable=SC1091
source ./devel/include/verbose.inc

if [ "${OC_USERNAME}" != "" ] && [ "${OC_PASSWORD}" != "" ] && [ "${OC_API_URL}" != "" ]; then
    oc login -u "${OC_USERNAME}" -p "${OC_PASSWORD}" "${OC_API_URL}" \
    || die "Failed to log in the cluster"
else
    oc whoami &>/dev/null \
    || die "You should be logged in the cluster before run this script"
fi

yield -n ">> " && verbose kustomize build config/default/ | kubectl create --dry-run --validate -f - \
&& yield -n ">> " && verbose kustomize build config/crd/ | kubectl create --dry-run --validate -f - \
&& yield -n ">> " && verbose kustomize build config/manager/ | kubectl create --dry-run --validate -f - \
&& yield -n ">> " && verbose kustomize build config/prometheus/ | kubectl create --dry-run --validate -f - 
