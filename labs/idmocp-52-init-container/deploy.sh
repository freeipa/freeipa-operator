#!/bin/bash

PRJ_NAME="avisiedo-freeipa-openshift-container"

function yield
{
    echo "$*" >&2
}

function error-msg
{
    yield "ERROR:$*"
}

function die
{
    local err=$?
    [ $err -eq 0 ] && err=127
    error-msg "$@"
    exit $err
}


# -----------------------------------------------------------------------------


echo ">> Checking log in OpenShift"
oc whoami &>/dev/null || die "Not logged in the cluster"
echo ">> Setting or creating the project"
oc project "${PRJ_NAME}" &>/dev/null || oc new-project "${PRJ_NAME}" || die "It could not be created the project"
echo ">> Creating Service Account"
oc get sa/freeipa &>/dev/null || oc create serviceaccount freeipa || die "Creating Service Account 'freeipa'"
oc adm policy add-scc-to-user anyuid -z freeipa
echo ">> Creating Persistent Volume"
oc get pv/freeipa &>/dev/null || oc apply -f pv.yaml || die "Creating Persistent Volume"
echo ">> Creating Persistent Volume Claim"
oc get pvc/freeipa &>/dev/null || oc apply -f pvc.yaml || die "Creating Persistent Volume Claim"
echo ">> Creating Secrets"
oc get secret/freeipa-password &>/dev/null || oc apply -f secret.yaml || die "Creating Secrets"
echo ">> Creating Pod"
oc get pod/freeipa-master &>/dev/null || oc apply -f pod.yaml || dier "Ceating Pod"

