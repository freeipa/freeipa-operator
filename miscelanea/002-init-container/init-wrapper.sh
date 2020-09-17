#!/bin/bash -xv

function run_init_wrapper
{
    unset INIT_WRAPPER
    # shellcheck disable=SC2086
    exec /usr/local/sbin/init "$@" ${init_extra_args}
}


function run_bash
{
    unset INIT_WRAPPER
    exec /bin/bash "$@"
}

init_extra_args=""
[ "${DEBUG_TRACE}" == "1" ] && init_extra_args="${init_extra_args} --verbose"
[ "${DEBUG_TRACE}" == "2" ] && {
    init_extra_args="${init_extra_args} --debug"
    export SYSTEMD_LOG_LEVEL="debug"
    export SYSTEMD_LOG_TARGET="console"
    export SYSTEMD_LOG_COLOR="no"
}

# shellcheck disable=SC2086
[ "${INIT_WRAPPER}" == "1" ] && run_init_wrapper "$@" ${init_extra_args}
[ "${INIT_WRAPPER}" == "2" ] && run_bash "$@"
unset INIT_WRAPPER

# shellcheck disable=SC2086
exec /usr/sbin/original/init "$@" ${init_extra_args}
