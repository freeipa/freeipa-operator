#!/bin/bash

echo -ne "\\nPort 80/tcp: "
# shellcheck disable=SC2016
oc exec pod/poc-004-1 -it -c client -- /bin/bash -c 'if [ "$( echo "hello" | nc -4 -w1  localhost 80 )" == "hello" ]; then echo OK; else echo FAILURE; fi'
oc logs pod/poc-004-1 -c http --tail 7

echo -ne "\\nPort 443/tcp: "
# shellcheck disable=SC2016
oc exec pod/poc-004-1 -it -c client -- /bin/bash -c 'if [ "$( echo "hello" | nc -4 -w1  localhost 443 )" == "hello" ]; then echo OK; else echo FAILURE; fi'
oc logs pod/poc-004-1 -c https --tail 7

echo -ne "\\nPort 389/tcp: "
# shellcheck disable=SC2016
oc exec pod/poc-004-1 -it -c client -- /bin/bash -c 'if [ "$( echo "hello" | nc -4 -w1  localhost 389 )" == "hello" ]; then echo OK; else echo FAILURE; fi'
oc logs pod/poc-004-1 -c ldap --tail 7

echo -ne "\\nPort 636/tcp: "
# shellcheck disable=SC2016
oc exec pod/poc-004-1 -it -c client -- /bin/bash -c 'if [ "$( echo "hello" | nc -4 -w1  localhost 636 )" == "hello" ]; then echo OK; else echo FAILURE; fi'
oc logs pod/poc-004-1 -c ldaps --tail 7

echo -ne "\\nPort 88/tcp: "
# shellcheck disable=SC2016
oc exec pod/poc-004-1 -it -c client -- /bin/bash -c 'if [ "$( echo "hello" | nc -4 -w1  localhost 88 )" == "hello" ]; then echo OK; else echo FAILURE; fi'
oc logs pod/poc-004-1 -c kerberos-tcp --tail 7

echo -ne "\\nPort 464/tcp: "
# shellcheck disable=SC2016
oc exec pod/poc-004-1 -it -c client -- /bin/bash -c 'if [ "$( echo "hello" | nc -4 -w1  localhost 464 )" == "hello" ]; then echo OK; else FAILURE; fi'
oc logs pod/poc-004-1 -c kerberos-admin-tcp --tail 7



echo -ne "\\nPort 88/udp: "
# shellcheck disable=SC2016
oc exec pod/poc-004-1 -it -c client -- /bin/bash -c 'if [ "$( echo "hello" | nc -4u -w1  localhost 88 )" == "hello" ]; then echo OK; else FAILURE; fi'
oc logs pod/poc-004-1 -c kerberos-udp --tail 7

echo -ne "\\nPort 464/udp: "
# shellcheck disable=SC2016
oc exec pod/poc-004-1 -it -c client -- /bin/bash -c 'if [ "$( echo "hello" | nc -4u -w1  localhost 464 )" == "hello" ]; then echo OK; else FAILURE; fi'
oc logs pod/poc-004-1 -c kerberos-admin-udp --tail 7
