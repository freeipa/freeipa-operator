#!/bin/bash

echo -ne "\\nPort 80/tcp: "
# shellcheck disable=SC2016
oc exec pod/poc-004-2-client -it -c client -- /bin/bash -c 'RESPONSE="$( echo -n "hello" | nc -4 -w2  poc-004-2 80 )"; if [ "${RESPONSE:0:5}" == "hello" ]; then echo OK; else echo FAILURE; fi'
oc logs pod/poc-004-2-services -c http --tail 7

echo -ne "\\nPort 443/tcp: "
# shellcheck disable=SC2016
oc exec pod/poc-004-2-client -it -c client -- /bin/bash -c 'RESPONSE="$( echo -n "hello" | nc -4 -w2  poc-004-2 443 )"; if [ "${RESPONSE:0:5}" == "hello" ]; then echo OK; else echo FAILURE; fi'
oc logs pod/poc-004-2-services -c https --tail 7

echo -ne "\\nPort 389/tcp: "
# shellcheck disable=SC2016
oc exec pod/poc-004-2-client -it -c client -- /bin/bash -c 'RESPONSE="$( echo -n "hello" | nc -4 -w2  poc-004-2 389 )"; if [ "${RESPONSE:0:5}" == "hello" ]; then echo OK; else echo FAILURE; fi'
oc logs pod/poc-004-2-services -c ldap --tail 7

echo -ne "\\nPort 636/tcp: "
# shellcheck disable=SC2016
oc exec pod/poc-004-2-client -it -c client -- /bin/bash -c 'RESPONSE="$( echo -n "hello" | nc -4 -w2  poc-004-2 636 )"; if [ "${RESPONSE:0:5}" == "hello" ]; then echo OK; else echo FAILURE; fi'
oc logs pod/poc-004-2-services -c ldaps --tail 7

echo -ne "\\nPort 88/tcp: "
# shellcheck disable=SC2016
oc exec pod/poc-004-2-client -it -c client -- /bin/bash -c 'RESPONSE="$( echo -n "hello" | nc -4 -w2  poc-004-2 88 )"; if [ "${RESPONSE:0:5}" == "hello" ]; then echo OK; else echo FAILURE; fi'
oc logs pod/poc-004-2-services -c kerberos-tcp --tail 7

echo -ne "\\nPort 464/tcp: "
# shellcheck disable=SC2016
oc exec pod/poc-004-2-client -it -c client -- /bin/bash -c 'RESPONSE="$( echo -n "hello" | nc -4 -w2  poc-004-2 464 )"; if [ "${RESPONSE:0:5}" == "hello" ]; then echo OK; else echo FAILURE; fi'
oc logs pod/poc-004-2-services -c kerberos-admin-tcp --tail 7



echo -ne "\\nPort 88/udp: "
# shellcheck disable=SC2016
oc exec pod/poc-004-2-client -it -c client -- /bin/bash -c 'RESPONSE="$( echo -n "hello" | nc -4u -w2  poc-004-2 88 )"; if [ "${RESPONSE:0:5}" == "hello" ]; then echo OK; else echo FAILURE; fi'
oc logs pod/poc-004-2-services -c kerberos-udp --tail 7

echo -ne "\\nPort 464/udp: "
# shellcheck disable=SC2016
oc exec pod/poc-004-2-client -it -c client -- /bin/bash -c 'RESPONSE="$( echo -n "hello" | nc -4u -w2  poc-004-2 464 )"; if [ "${RESPONSE:0:5}" == "hello" ]; then echo OK; else echo FAILURE; fi'
oc logs pod/poc-004-2-services -c kerberos-admin-udp --tail 7
