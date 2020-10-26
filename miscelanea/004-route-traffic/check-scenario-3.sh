#!/bin/bash

echo -ne "\\nPort 80/tcp: "
# shellcheck disable=SC2016
oc exec pod/poc-004-3-client -it -c client -- /bin/bash -c 'if [ "$( echo "hello" | nc -4 -w2  poc-004-3 80 )" == "hello" ]; then echo OK; else echo FAILURE; fi'
oc logs pod/poc-004-3-services -c http --tail 7

echo -ne "\\nPort 443/tcp: "
# shellcheck disable=SC2016
oc exec pod/poc-004-3-client -it -c client -- /bin/bash -c 'if [ "$( echo "hello" | nc -4 -w2  poc-004-3 443 )" == "hello" ]; then echo OK; else echo FAILURE; fi'
oc logs pod/poc-004-3-services -c https --tail 7

echo -ne "\\nPort 389/tcp: "
# shellcheck disable=SC2016
oc exec pod/poc-004-3-client -it -c client -- /bin/bash -c 'if [ "$( echo "hello" | nc -4 -w2  poc-004-3 389 )" == "hello" ]; then echo OK; else echo FAILURE; fi'
oc logs pod/poc-004-3-services -c ldap --tail 7

echo -ne "\\nPort 636/tcp: "
# shellcheck disable=SC2016
oc exec pod/poc-004-3-client -it -c client -- /bin/bash -c 'if [ "$( echo "hello" | nc -4 -w2  poc-004-3 636 )" == "hello" ]; then echo OK; else echo FAILURE; fi'
oc logs pod/poc-004-3-services -c ldaps --tail 7

echo -ne "\\nPort 88/tcp: "
# shellcheck disable=SC2016
oc exec pod/poc-004-3-client -it -c client -- /bin/bash -c 'if [ "$( echo "hello" | nc -4 -w2  poc-004-3 88 )" == "hello" ]; then echo OK; else echo FAILURE; fi'
oc logs pod/poc-004-3-services -c kerberos-tcp --tail 7

echo -ne "\\nPort 464/tcp: "
# shellcheck disable=SC2016
oc exec pod/poc-004-3-client -it -c client -- /bin/bash -c 'if [ "$( echo "hello" | nc -4 -w2  poc-004-3 464 )" == "hello" ]; then echo OK; else FAILURE; fi'
oc logs pod/poc-004-3-services -c kerberos-admin-tcp --tail 7



echo -ne "\\nPort 88/udp: "
# shellcheck disable=SC2016
oc exec pod/poc-004-3-client -it -c client -- /bin/bash -c 'if [ "$( echo "hello" | nc -4u -w1  poc-004-3 88 )" == "hello" ]; then echo OK; else FAILURE; fi'
oc logs pod/poc-004-3-services -c kerberos-udp --tail 7

echo -ne "\\nPort 464/udp: "
# shellcheck disable=SC2016
oc exec pod/poc-004-3-client -it -c client -- /bin/bash -c 'if [ "$( echo "hello" | nc -4u -w1  poc-004-3 464 )" == "hello" ]; then echo OK; else FAILURE; fi'
oc logs pod/poc-004-3-services -c kerberos-admin-udp --tail 7
