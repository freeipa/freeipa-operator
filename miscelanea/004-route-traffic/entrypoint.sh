#!/bin/sh

ECHO_PROTOCOL="${ECHO_PROTOCOL:-udp}"
ECHO_ADDRESS="${ECHO_ADDRESS:-0.0.0.0}"
ECHO_PORT="${ECHO_PORT:-8007}"
ECHO_MAX_CLIENTS="${ECHO_MAX_CLIENTS:-1}"

case "${ECHO_PROTOCOL}" in
    tcp )
        exec ./tcp-server.py
        ;;
    udp )
        exec ./udp-server.py
        ;;
    * )
        echo "ECHO_PROTOCOL '${ECHO_PROTOCOL}' unknown or unsupported"
        exit 1
        ;;
esac
