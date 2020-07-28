# https://gitlab.devsensation.es/snippets/10
##
# This include file provide minimal
##

INDLUDE_LIST=(${INCLUDE_LIST[@]} "verbose")
declare -g DEPS_VERBOSE=()

VERBOSE_ERROR=1
VERBOSE_WARNING=2
VERBOSE_INFO=3
VERBOSE_DEBUG=4
VERBOSE_DEFAULT=$VERBOSE_WARNING
VERBOSE=${VERBOSE:-$VERBOSE_DEFAULT}

##
# Print out the arguments to the standard error.
# $*
##
function yield
{
    echo "$*" >&2
} # yield


##
# Primitive to print out debug messages in stderr.
# $* The message to print out.
##
function debug-msg
{
    # shellcheck disable=SC2086
    [ $VERBOSE -ge $VERBOSE_DEBUG ] && yield "DEBUG:${FUNCNAME[1]}:$*"
} # info-msg


##
# Primitive to print out informative messages in stderr.
# $* The message to print out.
##
function info-msg
{
    # shellcheck disable=SC2086
    [ $VERBOSE -ge $VERBOSE_INFO ] && yield "INFO:${FUNCNAME[1]}:$*"
} # info-msg


##
# Primitive to print out warning messages in stderr.
# $* The message to print out.
##
function warning-msg
{
    # shellcheck disable=SC2086
    [ $VERBOSE -ge $VERBOSE_WARNING ] && yield "WARNING:${FUNCNAME[1]}:$*"
} # warning-msg


##
# Primitive to print out error messages in stderr.
# $* The message to print out.
##
function error-msg
{
    # shellcheck disable=SC2086
    [ $VERBOSE -ge $VERBOSE_ERROR ] && yield "ERROR:${FUNCNAME[1]}:$*"
} # error-msg


##
# Terminate the script execution with the last return
# code or 127 if no error code in the last operation.
# $@ Message to be printed as error into the standard
# error output.
##
function die
{
    local ERR=$?
    [ $ERR -eq 0 ] && ERR=127
    error-msg "$@"
    exit $ERR
} # die


##
# Try the operation and if it fails, finish the script
# execution.
##
function try
{
    "$@" || die "Trying '$*'"
} # try


##
# Assert some command or expression
##
function assert
{
    "$@" || die "Assertion failed at ${FUNCNAME[1]}: $*"
} # assert

