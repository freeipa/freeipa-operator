#!/bin/bash

##
# Helper script for linting files in the repository.
##

# shellcheck disable=SC1091
source "./devel/include/verbose.inc"

# Global lint filter bypass. It can be passed from the caller to bypass
# lint.ignore filtering, forcing to lint the specified file, or all the
# files found in the repository.
#   LINT_FILTER_BYPASS=1 ./devel/lint.sh ./devel/lint.sh
[ "${LINT_FILTER_BYPASS}" == "" ] && LINT_FILTER_BYPASS=0

declare -a SHELL_FILES
SHELL_FILES=()

declare -a YAML_FILES
YAML_FILES=()

declare -a DOCKERFILE_FILES
DOCKERFILE_FILES=()

declare -a GO_FILES
GO_FILES=()

declare -a MARKDOWN_FILES
MARKDOWN_FILES=()

declare -a UNKNOWN_FILES
UNKNOWN_FILES=()

FORCE=""

if command -v podman &>/dev/null; then
    oci="podman"
elif command -v docker &>/dev/null; then
    oci="docker"
else
    die "Container runtime not found"
fi

##
# Lint shellscripts
##
function lint-shellscript
{
    [ $# -eq 0 ] && return 0
    $oci run --rm -it \
             --volume "$PWD:/data:z" \
             --workdir "/data" \
             --entrypoint shellcheck \
             docker.io/nlknguyen/alpine-shellcheck:latest \
             -x "$@"
}


##
# Lint a Dockerfile files
##
function lint-dockerfile
{
    [ $# -eq 0 ] && return 0
    $oci run --rm -it \
             --volume "$PWD:/data:z" \
             --workdir "/data" \
             --entrypoint /bin/hadolint \
             hadolint/hadolint:latest \
             "$@" \
    || return 1
    return 0
}


##
# Lint a YAML files
##
function lint-yaml
{
    [ $# -eq 0 ] && return 0
    $oci   run --rm -it \
               --volume "$PWD:/data:z" \
               --workdir "/data" \
               --entrypoint yamllint \
               docker.io/cytopia/yamllint:latest \
               "$@" \
    || return 1
    return 0
}


##
# Lint GO files
##
function lint-go
{
    local reto
    [ $# -eq 0 ] && return 0
    [ "$GOPATH" == "" ] && die "GOPATH is not defined or empty"
    reto=0
    for file in "${@}"
    do
        $oci   run --rm -it \
                   -e "GOPATH=/go" \
                   --volume "$GOPATH:/go:z" \
                   --volume "$PWD:/data:z" \
                   --workdir "/data" \
                   --entrypoint golint \
                   docker.io/cytopia/golint:latest \
                   "${file}" \
        || reto=1
    done
    return $reto
}


##
# Lint a Markdown docouments
##
function lint-markdown
{
    local reto
    [ $# -eq 0 ] && return 0
    reto=0
    for file in "$@"
    do
        $oci   run --rm -it \
                   --volume "$PWD:/data:z" \
                   --workdir "/data" \
                   docker.io/markdownlint/markdownlint \
                   "${file}" \
        || reto=1
    done
    return $reto
}


##
# Lint a Kubernete objects
##
function lint-kubeobject
{
    [ $# -eq 0 ] && return 0
    $oci   run --rm -it \
               --volume "$PWD:/data:z" \
               --workdir "/data" \
               --entrypoint /bin/bash \
               docker.io/openshift/origin-cli:v4.0 \
               -c "oc login -u '${OC_USERNAME}' -p '${OC_PASSWORD}' --insecure-skip-tls-verify=true '${OC_API}' && oc apply --dry-run --validate -f '${filepath}'" \
    || return 1
    return 0
}


##
# Handle lint forced files for unknown files.
##
function lint-forced
{
    local force
    local linter_func
    force="$1"
    shift 1
    [ $# -eq 0 ] && return 0
    [ "$force" == "" ] && die "\$1 was expected to be some lint forced type"

    case "$force" in
        "--force-shellcheck" )
            linter_func="lint-shellscript"
            ;;
        "--force-dockerfile" )
            linter_func="lint-dockerfile"
            ;;
        "--force-yaml" )
            linter_func="lint-yaml"
            ;;
        "--force-go" )
            linter_func="lint-go"
            ;;
        "--force-kubeobject" )
            linter_func="lint-kubeobject"
            ;;
        "--force-markdown" )
            linter_func="lint-markdown"
            ;;
        * )
            return 0
            ;;
    esac

    "${linter_func}" "$@" || return 1
    return 0
}


function cmd-help
{
    cat <<EOF
Usage: ./devel/lint.sh [{options}] [files]

Options could be:

  --force-shellcheck  Force lint shellecheck when the type can not be
                      auto-discovered.
  --force-dockerfile  For lint Dockerfile when the type can not be
                      auto-discovered.
  --force-yaml        For lint yaml files when the type can not be
                      auto-discovered.
  --force-go          For lint go files when the type can not be
                      auto-discovered.
  --force-kubeobject  Force lint a Kubernete object when the type can
                      not be auto-discovered.
  --force-markdown    Force lint a Markdown document.

* By default lint all the files in the repository.
EOF
    exit 0
}


function is-in-lintignore
{
    local filepath="$1"
    filepath="${filepath#./}"
    [ "${LINT_FILTER_BYPASS}" -eq 0 ] && grep -q "^${filepath}\$" devel/lint.ignore && return 0
    return 1
}


function prepare-lists
{
    local DIRNAME
    for filepath in "$@"
    do
        [ ! -e "$filepath" ] && continue
        is-in-lintignore "${filepath}" && continue
        filename="$( basename "${filepath}" )"
        case "${filename}" in
            *.sh )
                SHELL_FILES+=("${filepath}")
                ;;
            "Dockerfile" | Dockerfile.* )
                DOCKERFILE_FILES+=("${filepath}")
                ;;
            *.md )
                MARKDOWN_FILES+=("${filepath}")
                ;;
            *.go )
                GO_FILES+=("${filepath}")
                ;;
            *.yaml | *.yml )
                DIRNAME="$( dirname "$filepath" )"
                DIRNAME="${DIRNAME#./}"
                if [ "${DIRNAME#config/crd/}" == "${DIRNAME}" ]
                then
                    YAML_FILES+=("${filepath}")
                fi
                ;;
            * )
                UNKNOWN_FILES+=("${filepath}")
                ;;
        esac
    done
}


function can-run-check-k8s
{
    command -v oc &>/dev/null || return 1
    oc whoami &>/dev/null && return 0
    [ "${OC_USERNAME}" != "" ] && [ "${OC_PASSWORD}" != "" ] && [ "${OC_API_URL}" != "" ] && return 0
    return 1
}


function cmd-lint-all
{
    local reto

    if [ $# -gt 0 ]
    then
        prepare-lists "$@"
    else
        # shellcheck disable=SC2046
        prepare-lists $( find . -name 'Dockerfile' -o -name 'Dockerfile.*' \
                                -o -name '*.md' \
                                -o -name '*.go' \
                                -o -name '*.sh'; \
                          find ./config/certmanager -name '*.yaml'; \
                          find ./config/crd -maxdepth 1 -name '*.yaml'; \
                          find ./config/default -name '*.yaml'; \
                          find ./config/manager -name '*.yaml'; \
                          find ./config/prometheus -name '*.yaml'; \
                          find ./config/rbac -name '*.yaml'; \
                          find ./config/samples -name '*.yaml'; \
                          find ./config/webhook -name '*.yaml'; \
                        )
    fi

    err_count=0

    # Linting shell script files
    lint-shellscript "${SHELL_FILES[@]}"
    reto=$?
    [ $reto -ne 0 ] && err_count=$(( err_count + 1 ))

    # Linting YAML files
    lint-yaml "${YAML_FILES[@]}"
    reto=$?
    [ $reto -ne 0 ] && err_count=$(( err_count + 1 ))

    # Linting Dockerfile files
    lint-dockerfile "${DOCKERFILE_FILES[@]}"
    reto=$?
    [ $reto -ne 0 ] && err_count=$(( err_count + 1 ))

    # Linting GO files
    lint-go "${GO_FILES[@]}"
    reto=$?
    [ $reto -ne 0 ] && err_count=$(( err_count + 1 ))

    # Linting Markdown files
    lint-markdown "${MARKDOWN_FILES[@]}"
    reto=$?
    [ $reto -ne 0 ] && err_count=$(( err_count + 1 ))

    [ ${#UNKNOWN_FILES[@]} -gt 0 ] && {
        if [ "${FORCE}" == "" ]
        then
            echo "The following files are not linted: ${UNKNOWN_FILES[*]}"
        else
            lint-forced "${FORCE}" "${UNKNOWN_FILES[@]}"
        fi
    }

    # Validate kubernetes objects if possible
    if can-run-check-k8s; then
        ./devel/check-k8s.sh || err_count=$(( err_count + 1 ))
    else
        warning-msg "It can not run validation on k8s kustomize objects"
    fi

    return $err_count
}


function cmd-run
{
    # Check help
    [ "$1" == "help" ] && cmd-help

    # Run the corresponding subcommand
    cmd-lint-all "$@"
    exit $?
}


function check-args-and-run
{
    local argument
    while [ "${1##--}" != "${1}" ]
    do
        argument="$1"
        shift 1
        case "${argument}" in
            "--force-shellcheck" \
            | "--force-dockerfile" \
            | "--force-yaml" \
            | "--force-go" \
            | "--force-markdown" \
            | "--force-kubeobject" )
                [ "$FORCE" != "" ] && die "Can not be forced two different linters"
                FORCE="${argument}"
                ;;
        esac
    done
    [ "${1##--}" != "${1}" ] && shift 1

    cmd-run "$@"
}


# Check repository root path
if [ ! -e "${PWD}/.git" ] \
   || [ ! -e "${PWD}/devel" ]
then
    die "This script must be launched from the repository root path"
fi


check-args-and-run "$@"
