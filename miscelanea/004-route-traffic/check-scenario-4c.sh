#!/bin/bash


function verbose
{
  echo "$@" >&2
  "$@"
}

for port in 31080 31443 31389 31636 31088 31464
do
  STDOUT="$( echo -n "hello world" | nc -4 test.apps.permanent.idmocp.lab.eng.rdu2.redhat.com "${port}" )"
  echo -ne "Port ${port}/tcp: "
  if [ "${STDOUT}" == "hello world" ]; then echo "Success"
  else echo ">> Failed"
  fi
done

for port in 31088 31464
do
  STDOUT="$( ( echo -n "hello world"; sleep 1 ) | nc -4u -w1 test.apps.permanent.idmocp.lab.eng.rdu2.redhat.com "${port}" )"
  echo -ne "Port ${port}/udp: "
  if [ "${STDOUT}" == "hello world" ]; then echo "Success"
  else echo ">> Failed"
  fi
done

verbose oc logs pod/poc-004-4c -c http  --tail 7
verbose oc logs pod/poc-004-4c -c https --tail 7
verbose oc logs pod/poc-004-4c -c ldap  --tail 7
verbose oc logs pod/poc-004-4c -c ldaps --tail 7
verbose oc logs pod/poc-004-4c -c kerberos-tcp --tail 7
verbose oc logs pod/poc-004-4c -c kerberos-admin-tcp --tail 7

verbose oc logs pod/poc-004-4c -c kerberos-udp --tail 7
verbose oc logs pod/poc-004-4c -c kerberos-admin-udp --tail 7
