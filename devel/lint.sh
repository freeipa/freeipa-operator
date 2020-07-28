#!/bin/bash

##
# Helper script for linting files in the repository.
##

# shellcheck source=./devel/include/verbose.inc
source "./devel/include/verbose.inc"


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

FORCE=""




##
# Lint shellscripts
##
function lint-shellscript
{
    local filepath
    filepath="$1"
    shift 1
    podman run --rm -it \
               --volume "$PWD:/data:z" \
               --workdir "/data" \
               --entrypoint shellcheck \
               docker.io/nlknguyen/alpine-shellcheck:latest \
               -x "${filepath}"
}


##
# Lint a Dockerfile file
##
function lint-dockerfile
{
    local filepath
    filepath="$1"
    shift 1
    podman run --rm -it \
               --volume "$PWD:/data:z" \
               --workdir "/data" \
               --entrypoint /bin/hadolint \
               hadolint/hadolint:latest \
               "${filepath}"
}


##
# Lint a YAML file
##
function lint-yaml
{
    local filepath
    filepath="$1"
    shift 1
    podman run --rm -it \
               --volume "$PWD:/data:z" \
               --workdir "/data" \
               --entrypoint yamllint \
               docker.io/cytopia/yamllint:latest \
               "${filepath}"
}


##
# Lint a GO file or directory
##
function lint-go
{
    local filepath
    filepath="$1"
    shift 1
    podman run --rm -it \
               --volume "$PWD:/data:z" \
               --workdir "/data" \
               --entrypoint golint \
               docker.io/cytopia/golint:latest \
               "${filepath}"
}


##
# Lint a Markdown docoument
##
function lint-markdown
{
    local filepath
    filepath="$1"
    shift 1
    podman run --rm -it \
               --volume "$PWD:/data:z" \
               --workdir "/data" \
               docker.io/markdownlint/markdownlint \
               "${filepath}"
}


##
# Lint a Kubernete object.
##
function lint-kubeobject
{
    local filepath
    filepath="$1"
    shift 1
    podman run --rm -it \
               --volume "$PWD:/data:z" \
               --workdir "/data" \
               --entrypoint /bin/bash \
               docker.io/openshift/origin-cli:v4.0 \
               -c "oc login -u '${OC_USERNAME}' -p '${OC_PASSWORD}' --insecure-skip-tls-verify=true '${OC_API}' && oc apply --dry-run --validate -f '${filepath}'"
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


function prepare-lists
{
    local files
    local file

    # array=()
    # while IFS='' read -r line; do array+=("$line"); done < <(mycommand)
    files=()
    while IFS='' read -r line; do files+=("$line"); done < <( find . -name '*.sh' )
    files=("${files[@]}" \
           "src/idmocp/build/bin/entrypoint" \
           "src/idmocp/build/bin/user_setup" \
    )
    SHELL_FILES=("${files[@]}")
    debug-msg "SHELL_FILES=${SHELL_FILES[*]}"


    files=()
    while IFS='' read -r line; do files+=("$line"); done < <( find . -name '*.yaml' -o -name '*.yml' )
    YAML_FILES=("${files[@]}")
    debug-msg "YAML_FILES=${YAML_FILES[*]}"


    files=()
    while IFS='' read -r line; do files+=("$line"); done < <( find . -name 'Dockerfile' -o -name 'Dockerfile.*' )
    DOCKERFILE_FILES=("${files[@]}")
    debug-msg "DOCKERFILE_FILES=${DOCKERFILE_FILES[*]}"


    files=()
    while IFS='' read -r line; do files+=("$line"); done < <( find . -name '*.go' )
    GO_FILES=("${files[@]}")
    debug-msg "GO_FILES=${GO_FILES[*]}"


    files=()
    while IFS='' read -r line; do files+=("$line"); done < <( find . -name '*.md' )
    MARKDOWN_FILES=("${files[@]}")
    debug-msg "MARKDOWN_FILES=${MARKDOWN_FILES[*]}"
}


function cmd-lint-all
{
    local reto
    local shellscript_err_count
    local yaml_err_count
    local dockerfile_err_count
    local go_err_count
    local markdown_err_count
    local err_count
    prepare-lists

    # Lint shell scripts
    echo ">> Linting shell script files"
    shellscript_err_count=0
    for file in "${SHELL_FILES[@]}"
    do
        lint-shellscript "$file"
        reto=$?
        [ $reto -ne 0 ] && shellscript_err_count=$(( shellscript_err_count + 1 ))
    done

    # Lint YAML files
    echo ">> Linting YAML files"
    yaml_err_count=0
    for file in "${YAML_FILES[@]}"
    do
        lint-yaml "$file"
        reto=$?
        [ $reto -ne 0 ] && yaml_err_count=$(( yaml_err_count + 1 ))
    done

    # Lint Dockerfile files
    echo ">> Linting Dockerfile files"
    dockerfile_err_count=0
    for file in "${DOCKERFILE_FILES[@]}"
    do
        lint-dockerfile "$file"
        reto=$?
        [ $reto -ne 0 ] && dockerfile_err_count=$(( dockerfile_err_count + 1 ))
    done

    # Lint Go Files
    echo ">> Linting GO files"
    go_err_count=0
    lint-go .
    reto=$?
    [ $reto -ne 0 ] && go_err_count=$(( go_err_count + 1 ))

    # Lint Markdown Files
    echo ">> Linting Markdown files"
    markdown_err_count=0
    for file in "${MARKDOWN_FILES[@]}"
    do
        lint-markdown "$file"
        reto=$?
        [ $reto -ne 0 ] && markdown_err_count=$(( markdown_err_count + 1 ))
    done

    err_count=$(( shellscript_err_count + yaml_err_count + dockerfile_err_count + go_err_count + markdown_err_count ))
    return $err_count
}


function lint-forced
{
    local force
    local filepath
    force="$1"
    filepath="$2"
    [ "$force" == "" ] && die "\$1 was expected to be some lint forced type"
    [ "$filepath" == "" ] && die "\$2 must be specified"

    case "$force" in
        "--force-shellcheck" )
            lint-shellscript "${filepath}"
            ;;
        "--force-dockerfile" )
            lint-dockerfile "${filepath}"
            ;;
        "--force-yaml" )
            lint-yaml "${filepath}"
            ;;
        "--force-go" )
            lint-go "${filepath}"
            ;;
        "--force-kubeobject" )
            lint-kubeobject "${filepath}"
            ;;
        "--force-markdown" )
            lint-markdown "${filepath}"
            ;;
        * )
            ;;
    esac
}

function cmd-lint-file-list
{
    local current_file
    local filename
    local err_count

    err_count=0
    while [ $# -gt 0 ]
    do
        current_file="$1"
        shift 1
        filename="$( basename "${current_file}" )"

        if [ "${filename%%.sh}" != "${filename}" ]
        then
            lint-shellscript "${current_file}"
            reto=$?
            [ $reto -ne 0 ] && err_count=$(( err_count + 1 ))
        elif [ "${filename}" == "Dockerfile" ] \
             || [ "${filename##Dockerfile.}" != "${filename}" ]
        then
            lint-dockerfile "${current_file}"
            reto=$?
            [ $reto -ne 0 ] && err_count=$(( err_count + 1 ))
        else
            if [ "$FORCE" != "" ]
            then
                lint-forced "$FORCE" "${current_file}"
                reto=$?
                [ $reto -ne 0 ] && err_count=$(( err_count + 1 ))
            else
                warning-msg "The '${filename}' not linted"
                err_count=$(( err_count + 1 ))
            fi
        fi
    done

    return $err_count
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

    cmd-run "$@"
}


function cmd-run
{
    # Check help
    [ "$1" == "help" ] && cmd-help

    # Run the corresponding subcommand
    if [ $# -eq 0 ]
    then
        echo ">> Linting all files in the repository"
        cmd-lint-all
        exit $?
    else
        echo ">> Linting a file list"
        cmd-lint-file-list "$@"
        exit $?
    fi
}

# Check repository root path
if [ ! -e "${PWD}/.git" ] \
   || [ ! -e "${PWD}/devel" ]
then
    die "This script must be launched from the repository root path"
fi


check-args-and-run "$@"
