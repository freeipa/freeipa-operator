#!/bin/bash
IMG_NAME="freeipa-operator"
CONTAINER_IMAGE_FILE="${CONTAINER_IMAGE_FILE:-${IMG_NAME}.tar}"

function yield
{
    echo "$*" >&2
}

function die
{
    local err=$?
    [ $err -eq 0 ] && err=127
    yield "ERROR: $*"
    exit $err
}

function verbose
{
    yield "$@"
    "$@"
}

if command -v podman &>/dev/null; then oci="podman"
elif command -v docker &>/dev/null; then oci="docker"
else die "No podman nor docker were found"
fi

args="--ci"
if [ "${CONTAINER_IMAGE_FILE}" != "" ] && [ -e "${CONTAINER_IMAGE_FILE}" ]; then
    args+=" --source docker-archive ${CONTAINER_IMAGE_FILE}"
    # shellcheck disable=SC2086
    verbose $oci run --rm -it \
                   -v "$PWD:$PWD" \
                   -w "$PWD" \
                   wagoodman/dive:latest \
                   ${args} "$@"
else
    case "$oci" in
        "docker" )
            args+=" --source docker"
            ;;
        "podman" )
            args+=" --source podman"
            ;;
    esac
    # shellcheck disable=SC2086
    verbose dive ${args} "$@"
fi

