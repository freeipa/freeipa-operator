#!/bin/bash

##
# This script install development local tools to make easier the getting
# started process for everyone.
#
# Usage:
#   ./devel/install-local-tools.sh
#
# For debuging something wrong in the script:
#   VERBOSE=5 ./devel/install-local-tools.sh
##


# shellcheck disable=SC1091
source "./devel/include/common.inc"


FLAG_INSTALL_VSCODE="${FLAG_INSTALL_VSCODE:-ASK}"
FLAG_RUN_VSCODE_AFTER_INSTALL="${FLAG_RUN_VSCODE_AFTER_INSTALL:-ASK}"
FLAG_INSTALL_CRC="${FLAG_INSTALL_CRC:-ASK}"


OPERATOR_SDK_VERSION="v1.0.0"
OPENSHIFT_CLIENT_ARCHIVE_URL="https://mirror.openshift.com/pub/openshift-v4/clients/ocp/latest-4.7/openshift-client-linux.tar.gz"


##
# Test if the system is a Debian distribution.
# @return Return true if it is debian else false.
##
function test-debian
{
    [ -e "/etc/debian_version" ] && return 0
    return 1
} # test-debian


##
# Test if the system is a Centos distribution.
# @return Return true if it is Centos else false.
##
function test-centos
{
    [ -e "/etc/centos-release" ] && return 0
    return 1
} # test-centos


##
# Test if the system is a Fedora distribution.
# @return Return true if it is Fedora else false.
##
function test-fedora
{
    [ -e "/etc/fedora-release" ] && return 0
    return 1
} # test-fedora


##
# Print the name of the linux distribution detected in the standard output.
##
function print-distribution
{
    if test-fedora; then printf "fedora\\n"
    elif test-centos; then printf "centos\\n"
    elif test-debian; then printf "debian\\n"
    else die "Can not identify the linux distribution"
    fi
} # print-distribution


##
# Print the linux distribution and the main version in the standard output.
##
function print-distribution-version
{
    local distribution
    local version
    distribution="$( print-distribution )"
    version=""

    case "${distribution}" in
        "fedora" )
            version="$( < /etc/fedora-release )"
            version="${version##Fedora release }"
            version="${version%% (*)*}"
            ;;
        "centos" )
            version="$( < /etc/centos-release )"
            version="${version##CentOS Linux release }"
            version="${version%%.* (Core)*}"
            ;;
        "debian" )
            version="$( < /etc/debian_version )"
            version="${version%%.*}"
            version="${version%%/*}"  # Ubuntu 18.04 return "buster/sid"
            ;;
        * )
            die "'${distribution}' linux distribution unsupported"
            ;;
    esac

    printf "%s-%s\\n" "${distribution}" "${version}"
} # print-distribution-version


##
# Print the package list that match the linus distribution in the standard
# output.
##
function print-package-list
{
    local package_list_file
    local package_list
    package_list_file="./devel/packages/$( print-distribution-version ).txt"
    [ -e "${package_list_file}" ] || die "'${package_list_file}' does not exist"
    package_list="$( < "${package_list_file}" )"
    printf "%s\\n" "${package_list}"
} # print-package-list


##
# Install packages for different linux distributions.
##
function install-packages
{
    local distribution
    local package_list
    local version
    distribution="$( print-distribution )"
    package_list="$( print-package-list )"
    version="$( print-distribution-version )"
    version="${version##*-}"

    case "$distribution" in
        "fedora" )
            if [[ ${version} -ge 30 ]]
            then
                # shellcheck disable=SC2086
                get-root dnf install -y ${package_list} \
                || die "Installing dnf packages: '${package_list}'"
            else
                # shellcheck disable=SC2086
                get-root yum install -y ${package_list} \
                || die "Installing yum packages: '${package_list}'"
            fi
            ;;
        "debian" )
            get-root apt-get update \
            || die "Updating debian packages"
            # shellcheck disable=SC2086
            get-root apt-get install -y ${package_list} \
            || die "Installing debian packages"
            ;;
        "centos" )
            if [[ ${version} -ge 8 ]]
            then
                # shellcheck disable=SC2086
                get-root dnf install -y ${package_list}
            else
                # shellcheck disable=SC2086
                get-root yum install -y ${package_list}
            fi
            ;;
        * )
            die "Fedora version '$FEDORA_VERSION' is not supported"
    esac
} # install-packages

# FIXME Clean-up this function
##
# Install kustomize tool from source code. The other ways failed in Fedora 32.
# See: https://kubernetes-sigs.github.io/kustomize/installation/source/#install-the-kustomize-cli-from-local-source-with-cloning-the-repo
##
function install-kustomize-from-source
{
    local reto=0
    local tmp_path="/tmp/kustomize-tmp"
    # Need go 1.13 or higher
    # unset GOPATH
    # see https://golang.org/doc/go1.13#modules
    # unset GO111MODULES

    ( [ ! -e "${tmp_path}" ] || rm -rf "${tmp_path}" ) \
    && mkdir "${tmp_path}"
    reto=$?
    [ $reto -eq 0 ] || die "Failing preparing for building from source"
    
    # shellcheck disable=SC2164
    pushd "${tmp_path}" &>/dev/null
    git clone --depth 1 --branch kustomize/v3.2.3 git@github.com:kubernetes-sigs/kustomize.git \
    && cd kustomize \
    && git checkout -b kustomize/v3.2.3 \
    && GO111MODULES='' go install .
    reto=$?
    # shellcheck disable=SC2164
    popd &>/dev/null
    rm -rf "${tmp_path}"

    return $reto
} # install-kustomize-from-source


##
# Install dive tool
##
function install-dive-tool
{
    local system
    local dive_version="0.10.0"
    system="$( uname -s )"
    case "${system}" in
        "Linux" )
            curl --silent -L "https://github.com/wagoodman/dive/releases/download/v${dive_version}/dive_${dive_version}_linux_amd64.tar.gz" \
            | tar xz -C "${GOPATH}/bin/"
            ;;
        "Darwin" )
            curl --silent -L "https://github.com/wagoodman/dive/releases/download/v${dive_version}/dive_${dive_version}_darwin_amd64.tar.gz" \
            | tar xz -C "${GOPATH}/bin/"
            ;;
        * )
            die "'${system}' is not supported"
    esac
} # install-dive-tool


##
# Install golang tools
##
function install-go-tools
{
    local GODOC_VERSION="v0.1.6"
    local DELVE_VERSION="v1.7.1"
    local GOLINT_VERSION="master"
    local KUSTOMIZE_VERSION="v3.2.3"

    [ "$GOPATH" == "" ] && {
        GOPATH="$HOME/go"
        [ -e "$GOPATH" ] || mkdir -p "$GOPATH"
        export PATH="$GOPATH/bin:$PATH"
    }

    # Install godoc
    go install golang.org/x/tools/cmd/godoc@${GODOC_VERSION} || die "Installing godoc"

    # Install debugger
    go install github.com/go-delve/delve/cmd/dlv@${DELVE_VERSION} || die "Installing dlv"

    # Install linter
    go install golang.org/x/lint/golint@${GOLINT_VERSION} || die "Installing golint"

    # Install kustomize
    go install sigs.k8s.io/kustomize/kustomize/v3@${KUSTOMIZE_VERSION} || die "Installing kustomize"

    # Install dive
    install-dive-tool
} # install-go-tools


##
# Install Visual Studio Code.
##
function install-vscode
{
    local distribution
    local package_list
    local version
    distribution="$( print-distribution )"
    version="$( print-distribution-version )"
    version="${version##*-}"

    case "$distribution" in
        "fedora" )
            get-root rpm --import "https://packages.microsoft.com/keys/microsoft.asc"
            cat <<EOF | get-root tee /etc/yum.repos.d/vscode.repo 1>/dev/null
[code]
name=Visual Studio Code
baseurl=https://packages.microsoft.com/yumrepos/vscode
enabled=1
gpgcheck=1
gpgkey=https://packages.microsoft.com/keys/microsoft.asc
EOF
            if [[ ${version} -ge 30 ]]
            then
                get-root dnf check-update
                get-root dnf install -y code \
                || die "Installing dnf packages: 'code'"
            else
                get-root yum update
                get-root yum install -y code \
                || die "Installing yum packages: 'code'"
            fi
            ;;
        "debian" )
            # https://linuxhint.com/install_visual_studio_code_debian_10/
            if [[ ${version} -ge 10 ]]
            then
                curl -L -o /tmp/code.deb "https://go.microsoft.com/fwlink/?LinkID=760868" \
                || die "Downloading VSCode package"
                get-root apt update \
                || die "Updating debian packages"
                get-root apt install -y /tmp/code.deb \
                || die "Installing VSCode package"
            else
                die "Visual Studio Code is supported starting at Debian 10 (Buster)"
            fi
            ;;
        "centos" )
            get-root rpm --import "https://packages.microsoft.com/keys/microsoft.asc"
            cat <<EOF | get-root tee /etc/yum.repos.d/vscode.repo 1>/dev/null
[code]
name=Visual Studio Code
baseurl=https://packages.microsoft.com/yumrepos/vscode
enabled=1
gpgcheck=1
gpgkey=https://packages.microsoft.com/keys/microsoft.asc
EOF
            if [[ ${version} -ge 8 ]]
            then
                get-root dnf check-update
                get-root dnf install -y code \
                || die "Installing dnf packages: 'code'"
            else
                die "Starting support at Centos 8"
            fi
            ;;
        * )
            die "Fedora version '$version' is not supported"
            ;;
    esac
} # install-vscode


BIN_FILE="$0"
[ "$( dirname "${BIN_FILE}" )" == "./devel" ] || die "Run this script from the repository base path as: ./devel/install-local-tools.sh"


# Install packages
echo ">> Installing packages"
verbose install-packages

# Install golang tools
echo ">> Installing golang tools"
verbose install-go-tools

# Install OpenShift Client
command -v oc 2>/dev/null 1>/dev/null || {
    echo ">> Installing OpenShift Client"
    TMP_FILE="/tmp/openshift-client.tar.gz"
    [ ! -e "${TMP_FILE}" ] || rm -f "${TMP_FILE}" || die "Removing '${TMP_FILE}'"
    verbose curl -o "${TMP_FILE}" --silent -L "$OPENSHIFT_CLIENT_ARCHIVE_URL" \
    && verbose get-root tar xzf "${TMP_FILE}" -C /usr/local/bin/ \
    && oc completion bash | get-root tee /usr/share/bash-completion/completions/oc 1>/dev/null \
    && rm -f "${TMP_FILE}"
} || die "Installing OpenShift Client"


# Install operator-sdk
command -v operator-sdk 2>/dev/null 1>/dev/null || {
    echo ">> Installing Operator-SDK ${OPERATOR_SDK_VERSION}"
    ( verbose curl -L --silent -o "operator-sdk" "https://github.com/operator-framework/operator-sdk/releases/download/${OPERATOR_SDK_VERSION}/operator-sdk-${OPERATOR_SDK_VERSION}-$( uname -m )-linux-gnu" \
      && verbose chmod a+x operator-sdk \
      && verbose get-root mv operator-sdk /usr/local/bin/operator-sdk \
      && operator-sdk completion bash \
         | get-root tee /usr/share/bash-completion/completions/operator-sdk 1>/dev/null
    ) || die "Installing Operator-SDK ${OPERATOR_SDK_VERSION}"
}


# Install Visual Studio Code (optional)
command -v code 2>/dev/null 1>/dev/null || {
    echo "Visual Studio Code is an IDE which could provide a complete solution for"
    echo "working with the different technologies used with operator-sdk"

    if [ "$FLAG_INSTALL_VSCODE" == "YES" ] \
       || ( [ "$FLAG_INSTALL_VSCODE" == "ASK" ] \
            && confirm "Do you want to install Visual Studio Code?"
          )
    then
        echo ">> Installing repository for 'code' package"
        install-vscode
        echo "Recommended plug-ins"
        echo "- Go"
        echo "- Ansible"
        echo "- YAML"
        echo "- OpenShift Extension Pack"
        echo "- Knative"
        echo "- Docker"
    fi

    command -v code 2>/dev/null 1>/dev/null \
    && echo ">> Visual Studio Code is installed" \
    &&  if [ "$FLAG_RUN_VSCODE_AFTER_INSTALL" == "YES" ] \
           || ( [ "$FLAG_RUN_VSCODE_AFTER_INSTALL" == "ASK" ] \
                && confirm "Do you want to open Visual Studio Code now?"
              )
        then
            code
        fi
}


# Install CodeReady Containers
command -v crc 2>/dev/null 1>/dev/null || {
    if  [ "$FLAG_INSTALL_CRC" == "YES" ] \
        || ( [ "$FLAG_INSTALL_CRC" == "ASK" ] \
             && confirm "Do you want to install CodeReady Container?"
           )
    then
        echo ">> Installing CodeReadyContainer"
        [ ! -e pull-secret ] && [ ! -e "$( xdg-user-dir DOWNLOAD )/pull-secret" ] && {
            echo "WARNING:File '$PWD/pull-secret' does not exist"
            echo ">> Please download 'pull-secret' file in your Download file"
            echo "Opening browser for getting it at: https://cloud.redhat.com/openshift/install/crc/installer-provisioned"
            open-url "https://cloud.redhat.com/openshift/install/crc/installer-provisioned"
        }
        if [ ! -e pull-secret ] && [ -e "$( xdg-user-dir DOWNLOAD )/pull-secret" ]
        then
            mv "$( xdg-user-dir DOWNLOAD )/pull-secret" "$PWD/pull-secret"
        fi
        # [ ! -e pull-secret ] && die "Can not find 'pull-secret' file"

        [ ! -e crc-linux-amd64.tar.xz ] && [ ! -e "$( xdg-user-dir DOWNLOAD )/crc-linux-amd64.tar.xz" ] && {
            echo "ERROR:File '${PWD}/crc-linux-amd64.tar.xz' does not exist"
            echo "Opening browser for getting it at: https://cloud.redhat.com/openshift/install/crc/installer-provisioned"
            echo "Please download the linux version and copy it to '${PWD}' directory."
            open-url "https://cloud.redhat.com/openshift/install/crc/installer-provisioned"
            echo "INFO: Rerun 'install-local-tools' when the file has been downloaded"
            exit 1
        }
        [ -e crc-linux-amd64.tar.xz ] \
        || mv "$( xdg-user-dir DOWNLOAD )/crc-linux-amd64.tar.xz" "crc-linux-amd64.tar.xz"

        # Extracting crc binary
        [ ! -e crc-linux ] || rm -rf crc-linux
        mkdir crc-linux
        try tar xf crc-linux-amd64.tar.xz -C ./crc-linux --strip-components 1 \
        && ( [ -e "${HOME}/.local/bin" ] || mkdir -p "${HOME}/.local/bin" ) \
        && cp -f crc-linux/crc "${HOME}/.local/bin/crc" \
        && rm -rf crc-linux

        echo ">> Deploying VM with OpenShift"
        verbose crc setup || die "Setting up CodeReady Containers"

        cat <<EOF
>> Local CodeReadyContainer cluster
To create the local cluster run 'crc start --pull-secret-file $PWD/pull-secret'
Run 'crc console --credentials' if you forget them.
EOF
    fi
}


command -v "yml2json" &>/dev/null || {
	pip3 install yml2json
}

command -v "json2yaml" &>/dev/null || {
	pip3 install json2yaml
}

echo ">> Installing git-hooks"
if [ "${GITHUB_WORKSPACE}" == "" ]; then
    echo ">> Installing git-hooks"
    cp -f ./devel/pre-commit.sh .git/hooks/pre-commit
    unset GOPATH

    # shellcheck disable=SC1090
    source "$HOME/.bashrc"
    [ "$GOPATH" == "" ] && {
        cat <<EOF
Append the below at the end of your .bashrc file:
# Golang support
export GOPATH="$HOME/go"
export PATH="$GOPATH/bin:$PATH"
EOF
    }
fi

exit 0
