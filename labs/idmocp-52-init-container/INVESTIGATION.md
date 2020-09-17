# Init Container Investigation

Running from cluster node:

```shell
oc get nodes
oc debug node/permanent-4f5zq-worker-hpks5
chroot /host

podman pull quay.io/avisied0/freeipa-openshift-container:latest

mkdir -p /tmp/freeipa

cat > /tmp/freeipa/init <<EOF
#!/bin/bash
export IPA_SERVER_HOSTNAME="\$( hostname )"
export IPA_SERVER_IP="\$( ip addr show dev eth0 | grep inet\ | awk '{ print \$2}' | awk -F '/' '{ print \$1 }' )"
export PASSWORD="administrator"
export IPA_SERVER_INSTALL_OPTS=" --ds-password=\${PASSWORD} --admin-password=\${PASSWORD} --ip-address=\${IPA_SERVER_IP} --domain=www.freeipa.test --realm=WWW.FREEIPA.TEST --no-host-dns --setup-kra --setup-adtrust --mkhomedir --no-ntp --no-sshd --no-ssh"


# find /data-template \
# | while read -r source
# do
#   target="/data\${source##/data-template}"
#   [ -d "\${source}" ] && mkdir -p "\${target}"
#   [ ! -d "\${source}" ] && cp -vf "\${source}" "\${target}"
# done

# touch /data/etc/krb5.keytab
# touch /data/etc/named.keytab
# touch /data/etc/ntp.conf
# touch /data/etc/rndc.key
# touch /data/etc/yp.conf
# mkdir -p /tmp/var/tmp


# mkdir -p /data/etc/certmonger
# mkdir -p /data/etc/dirsrv
# mkdir -p /data/etc/gssproxy
# touch /data/etc/krb5.conf
# mkdir -p /data/etc/krb5.conf.d
# touch /data/etc/krb5.keytab
# touch /data/etc/named.conf
# touch /data/etc/named.keytab
# touch /data/etc/named.keytab
# mkdir -p /data/etc/openldap
# mkdir -p /data/etc/pam.d
# touch /data/etc/ntp.conf
# touch /data/etc/rndc.key
# touch /data/etc/yp.conf
# mkdir -p /data/etc/ipa
# mkdir -p /data/etc/samba
# mkdir -p /data/etc/sssd
# mkdir -p /data/etc/systemd/system
# mkdir -p /data/etc/sysconfig
# mkdir -p /data/etc/tmpfiles.d
# mkdir -p /data/etc/sysconfig

# mkdir -p /data/var/log
# mkdir -p /data/var/lib/authconfig
# mkdir -p /tmp/var/lib/certmonger
# mkdir -p /tmp/var/lib/chrony
# mkdir -p /data/var/lib/dirsrv
# mkdir -p /tmp/var/lib/gssproxy
# mkdir -p /data/var/lib/ipa
# mkdir -p /data/var/lib/ipa-client
# mkdir -p /data/var/lib/pki
# mkdir -p /data/var/lib/samba
# mkdir -p /data/var/lib/systemd
# mkdir -p /data/var/lib/ipa/sysrestore

# Generate /etc/machine-id
# systemd-machine-id-setup --commit
systemd-id128 new > /etc/machine-id

cat > /etc/nsswitch.conf <<< "hosts: files dns\n"
# Launch installer.sh in background
nohup /installer.sh &
# Pass the control to systemd (this process has PID 1)
# exec /lib/systemd/systemd    # Wrong init file
exec /usr/local/sbin/init      # Thanks #chaimes
# exec bash

EOF

cat > /tmp/freeipa/installer.sh <<EOF
#!/bin/bash
export IPA_SERVER_HOSTNAME="\$( hostname )"
export IPA_SERVER_IP="\$( ip addr show dev eth0 | grep inet\  | awk '{ print \$2}' | awk -F '/' '{ print \$1 }' )"
export PASSWORD="administrator"
export IPA_SERVER_INSTALL_OPTS=" --ds-password=\${PASSWORD} --admin-password=\${PASSWORD} --ip-address=\${IPA_SERVER_IP} --domain=www.freeipa.test --realm=WWW.FREEIPA.TEST --no-host-dns --setup-kra --setup-adtrust --mkhomedir --no-ntp --no-sshd --no-ssh"

function yield
{
  echo "\$*" >&2
}

function verbose
{
  yield "\$*"
  "\$@"
}

while ! pgrep systemd &>/dev/null
do
  sleep 5
done
# Change this for some indicator checking
sleep 5
verbose ipa-server-install \${IPA_SERVER_INSTALL_OPTS} --unattended --hostname="www.freeipa-master.test"

EOF

# chmod a+x /tmp/freeipa/installer.sh
# chmod a+x /tmp/freeipa/init

# container_id="$( podman run --systemd=true \
#                             -P -d \
#                             -h www.freeipa-master.test \
#                             -v "/sys/fs/cgroup:/sys/fs/cgroup:ro" \
#                             -v /data \
#                             -v "/tmp/freeipa/init:/init:z" \
#                             -v "/tmp/freeipa/installer.sh:/installer.sh:z" \
#                             --entrypoint "" \
#                             quay.io/avisied0/freeipa-openshift-container:dev-test / \
#                )"
sudo semanage sebool container_manage_cgroup true
container_id="$( podman run --systemd=true \
                            -P -d \
                            -h www.freeipa-master.test \
                            -v "/sys/fs/cgroup:/sys/fs/cgroup:rw" \
                            -v /data \
                            --entrypoint "" \
                            quay.io/avisied0/freeipa-openshift-container:dev-test /usr/sbin/init \
               )"
podman logs -f "${container_id}"

podman run --systemd=true \
                            -P -d \
                            -h www.freeipa-master.test \
                            -v "/sys/fs/cgroup:/sys/fs/cgroup:rw" \
                            -v /data \
                            -e DATA=/data \
                            -e HOST=www.freeipa-master.test \
                            --entrypoint "" \
                            docker.io/freeipa/freeipa-server:fedora-32 \
                            /usr/local/sbin/init /bin/install.sh
```


podman run -P -it \
           --systemd=true \
           -h www.freeipa-master.test \
           -v "/sys/fs/cgroup:/sys/fs/cgroup:ro" \
           -v /data \
           -v "/tmp/freeipa/init:/init:z" \
           -v "/tmp/freeipa/installer.sh:/installer.sh:z" \
           -e DATA=/data \
           -e HOST=www.freeipa-master.test \
           --entrypoint "" \
           docker.io/freeipa/freeipa-server:fedora-32 \
           /usr/local/sbin/init /bin/install.sh --hostname="www.freeipa-master.test"

## References

- [How to run systemd in a container](https://developers.redhat.com/blog/2019/04/24/how-to-run-systemd-in-a-container/).
- https://github.com/freeipa/freeipa-container/issues/301

## Tests


### In the VM

```shell
# In a VM
sudo ipa-server-install --unattended -p administrator -a administrator --realm=example.com --hostname=www.example.com --domain=example.com --no-host-dns --setup-adtrust --setup-kra --mkhomedir --no-ntp --no-ssh --no-sshd --no-dns-sshfp
```

### Interesting finds in travis-ci pipeline for freeipa-openshift-container

Commands used as starting point from the tests that Jan have in the pipeline.

```shell
# travis-ci.org for fedora-32
docker run --read-only --dns=127.0.0.1 -d --name freeipa-master -h ipa.example.test --sysctl net.ipv6.conf.all.disable_ipv6=0 --tmpfs /run --tmpfs /tmp -v /dev/urandom:/dev/random:ro -v /sys/fs/cgroup:/sys/fs/cgroup:ro -v /tmp/freeipa-test-6447/data:/data:Z -v /tmp/freeipa-test-6447/data/etc/machine-id:/etc/machine-id:ro,Z -e INIT_WRAPPER=1 -e PASSWORD=Secret123 local/freeipa-server:fedora-32 exit-on-finished -U -r EXAMPLE.TEST --setup-dns --no-forwarders --auto-reverse --allow-zone-overlap --no-ntp
```

### Tests launched directly in the node

Prepare the node with:

```shell
oc debug node/avisiedo7-x5mpl-worker-p55cj
chroot /host
setsebool container_manage_cgroup 1
getsebool container_manage_cgroup
```

```shell
# The one I am using in 'oc debug node/...'
podman run --read-only \
           --dns=127.0.0.1 \
           --interactive \
           --name freeipa-master \
           --hostname ipa.example.test \
           --sysctl net.ipv6.conf.all.disable_ipv6=0 \
           --tmpfs /run \
           --tmpfs /tmp \
           --volume /dev/urandom:/dev/random:ro \
           --volume /sys/fs/cgroup:/sys/fs/cgroup:ro \
           --volume /tmp/freeipa-test/data:/data:Z \
           --volume /etc/machine-id:/etc/machine-id:ro,Z \
           --env INIT_WRAPPER=2 \
           --env PASSWORD=Secret123 \
           quay.io/avisied0/freeipa-openshift-container:dev-test -U \
           -r EXAMPLE.TEST \
           --setup-dns \
           --no-forwarders \
           --auto-reverse \
           --allow-zone-overlap \
           --no-ntp

```

---

### It can not be launched directly from the node

```text
sh-4.4# podman --log-level debug run --rm --read-only            --dns=127.0.0.1            --interactive            --name freeipa-master            --hostname ipa.example.test            --sysctl net.ipv6.conf.all.disable_ipv6=0           --tmpfs /run            --tmpfs /tmp            --volume /dev/urandom:/dev/random:ro            --volume /sys/fs/cgroup:/sys/fs/cgroup:ro            --volume /tmp/freeipa-test/data:/data:Z            --volume /etc/machine-id:/etc/machie-id:ro,Z            --env INIT_WRAPPER=2            --env PASSWORD=Secret123      --entrypoint  ""  --workdir /  quay.io/avisied0/freeipa-openshift-container:dev-test ls -l
DEBU[0000] Reading configuration file "/usr/share/containers/libpod.conf" 
DEBU[0000] Merged system config "/usr/share/containers/libpod.conf": &{{false false false false false true} 0 {   [] [] []}  docker://  runc map[crun:[/usr/bin/crun /usr/local/bin/crun] runc:[/usr/bin/runc /usr/sbin/runc /usr/local/bin/runc /usr/local/sbin/runc /sbin/runc /bin/runc /usr/lib/cri-o-runc/sbin/runc /run/current-system/sw/bin/runc]] [crun runc] [crun] [] [/usr/libexec/podman/conmon /usr/local/libexec/podman/conmon /usr/local/lib/podman/conmon /usr/bin/conmon /usr/sbin/conmon /usr/local/bin/conmon /usr/local/sbin/conmon /run/current-system/sw/bin/conmon] [PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin] systemd   /var/run/libpod -1 false /etc/cni/net.d/ [/usr/libexec/cni /usr/lib/cni /usr/local/lib/cni /opt/cni/bin] podman []   k8s.gcr.io/pause:3.1 /pause false false  2048 shm    false} 
DEBU[0000] Using conmon: "/usr/bin/conmon"              
DEBU[0000] Initializing boltdb state at /var/lib/containers/storage/libpod/bolt_state.db 
DEBU[0000] Using graph driver overlay                   
DEBU[0000] Using graph root /var/lib/containers/storage 
DEBU[0000] Using run root /var/run/containers/storage   
DEBU[0000] Using static dir /var/lib/containers/storage/libpod 
DEBU[0000] Using tmp dir /var/run/libpod                
DEBU[0000] Using volume path /var/lib/containers/storage/volumes 
DEBU[0000] Set libpod namespace to ""                   
DEBU[0000] [graphdriver] trying provided driver "overlay" 
DEBU[0000] cached value indicated that overlay is supported 
DEBU[0000] cached value indicated that metacopy is not being used 
DEBU[0000] NewControl(/var/lib/containers/storage/overlay): nextProjectID = 2 
DEBU[0000] cached value indicated that native-diff is usable 
DEBU[0000] backingFs=xfs, projectQuotaSupported=true, useNativeDiff=true, usingMetacopy=false 
DEBU[0000] Initializing event backend journald          
DEBU[0000] using runtime "/usr/bin/runc"                
WARN[0000] Error initializing configured OCI runtime crun: no valid executable found for OCI runtime crun: invalid argument 
INFO[0000] Found CNI network crio-bridge (type=bridge) at /etc/cni/net.d/100-crio-bridge.conf 
INFO[0000] Found CNI network 200-loopback.conf (type=loopback) at /etc/cni/net.d/200-loopback.conf 
INFO[0000] Found CNI network podman (type=bridge) at /etc/cni/net.d/87-podman-bridge.conflist 
DEBU[0000] parsed reference into "[overlay@/var/lib/containers/storage+/var/run/containers/storage]quay.io/avisied0/freeipa-openshift-container:dev-test" 
DEBU[0000] parsed reference into "[overlay@/var/lib/containers/storage+/var/run/containers/storage]@470c3b534682ca55a9fc41a9d8fc2be9c27f2a15d24cb2a0db40a6a59423960f" 
DEBU[0000] User mount /dev/urandom:/dev/random options [ro] 
DEBU[0000] User mount /sys/fs/cgroup:/sys/fs/cgroup options [ro] 
DEBU[0000] User mount /tmp/freeipa-test/data:/data options [Z] 
DEBU[0000] User mount /etc/machine-id:/etc/machine-id options [ro Z] 
DEBU[0000] Using bridge netmode                         
DEBU[0000] Adding mount /proc                           
DEBU[0000] Adding mount /dev                            
DEBU[0000] Adding mount /dev/pts                        
DEBU[0000] Adding mount /dev/mqueue                     
DEBU[0000] Adding mount /sys                            
DEBU[0000] setting container name freeipa-master        
DEBU[0000] created OCI spec and options for new container 
DEBU[0000] Allocated lock 0 for container 7477a79b5cf4aa1760f4138c20742f218d991de0e49602a480eac2f64d8d968b 
DEBU[0000] parsed reference into "[overlay@/var/lib/containers/storage+/var/run/containers/storage]@470c3b534682ca55a9fc41a9d8fc2be9c27f2a15d24cb2a0db40a6a59423960f" 
DEBU[0000] created container "7477a79b5cf4aa1760f4138c20742f218d991de0e49602a480eac2f64d8d968b" 
DEBU[0000] container "7477a79b5cf4aa1760f4138c20742f218d991de0e49602a480eac2f64d8d968b" has work directory "/var/lib/containers/storage/overlay-containers/7477a79b5cf4aa1760f4138c20742f218d991de0e49602a480eac2f64d8d968b/userdata" 
DEBU[0000] container "7477a79b5cf4aa1760f4138c20742f218d991de0e49602a480eac2f64d8d968b" has run directory "/var/run/containers/storage/overlay-containers/7477a79b5cf4aa1760f4138c20742f218d991de0e49602a480eac2f64d8d968b/userdata" 
DEBU[0000] Creating new volume 28f3c352f7aeef62b441fd31f3a38931edfb8f2dd8e4111324eb149512e250ab for container 
DEBU[0000] Validating options for local driver          
DEBU[0000] New container created "7477a79b5cf4aa1760f4138c20742f218d991de0e49602a480eac2f64d8d968b" 
DEBU[0000] container "7477a79b5cf4aa1760f4138c20742f218d991de0e49602a480eac2f64d8d968b" has CgroupParent "machine.slice/libpod-7477a79b5cf4aa1760f4138c20742f218d991de0e49602a480eac2f64d8d968b.scope" 
DEBU[0000] Made network namespace at /var/run/netns/cni-d26f2fff-2b88-e687-b24f-c40f73b772bd for container 7477a79b5cf4aa1760f4138c20742f218d991de0e49602a480eac2f64d8d968b 
INFO[0000] Got pod network &{Name:freeipa-master Namespace:freeipa-master ID:7477a79b5cf4aa1760f4138c20742f218d991de0e49602a480eac2f64d8d968b NetNS:/var/run/netns/cni-d26f2fff-2b88-e687-b24f-c40f73b772bd Networks:[] RuntimeConfig:map[podman:{IP: PortMappings:[] Bandwidth:<nil> IpRanges:[]}]} 
INFO[0000] About to add CNI network cni-loopback (type=loopback) 
DEBU[0000] overlay: mount_data=lowerdir=/var/lib/containers/storage/overlay/l/S2K53IU53UA6QYUZS7FNZNC4FP:/var/lib/containers/storage/overlay/l/R7PYBTLPBJOTCOVDBXLTDIG5XZ,upperdir=/var/lib/containers/storage/overlay/529b01f7c138e0f4d767cbc55941f5036faed85280627e33dc416bb0c091d9bc/diff,workdir=/var/lib/containers/storage/overlay/529b01f7c138e0f4d767cbc55941f5036faed85280627e33dc416bb0c091d9bc/work,context="system_u:object_r:container_file_t:s0:c259,c815" 
DEBU[0000] mounted container "7477a79b5cf4aa1760f4138c20742f218d991de0e49602a480eac2f64d8d968b" at "/var/lib/containers/storage/overlay/529b01f7c138e0f4d767cbc55941f5036faed85280627e33dc416bb0c091d9bc/merged" 
DEBU[0000] Copying up contents from container 7477a79b5cf4aa1760f4138c20742f218d991de0e49602a480eac2f64d8d968b to volume 28f3c352f7aeef62b441fd31f3a38931edfb8f2dd8e4111324eb149512e250ab 
DEBU[0000] Creating dest directory: /var/lib/containers/storage/volumes/28f3c352f7aeef62b441fd31f3a38931edfb8f2dd8e4111324eb149512e250ab/_data 
DEBU[0000] Calling TarUntar(/var/lib/containers/storage/overlay/529b01f7c138e0f4d767cbc55941f5036faed85280627e33dc416bb0c091d9bc/merged/var/log/journal, /var/lib/containers/storage/volumes/28f3c352f7aeef62b441fd31f3a38931edfb8f2dd8e4111324eb149512e250ab/_data) 
DEBU[0000] TarUntar(/var/lib/containers/storage/overlay/529b01f7c138e0f4d767cbc55941f5036faed85280627e33dc416bb0c091d9bc/merged/var/log/journal /var/lib/containers/storage/volumes/28f3c352f7aeef62b441fd31f3a38931edfb8f2dd8e4111324eb149512e250ab/_data) 
DEBU[0000] Created root filesystem for container 7477a79b5cf4aa1760f4138c20742f218d991de0e49602a480eac2f64d8d968b at /var/lib/containers/storage/overlay/529b01f7c138e0f4d767cbc55941f5036faed85280627e33dc416bb0c091d9bc/merged 
INFO[0000] Got pod network &{Name:freeipa-master Namespace:freeipa-master ID:7477a79b5cf4aa1760f4138c20742f218d991de0e49602a480eac2f64d8d968b NetNS:/var/run/netns/cni-d26f2fff-2b88-e687-b24f-c40f73b772bd Networks:[] RuntimeConfig:map[podman:{IP: PortMappings:[] Bandwidth:<nil> IpRanges:[]}]} 
INFO[0000] About to add CNI network podman (type=bridge) 
DEBU[0000] [0] CNI result: Interfaces:[{Name:cni-podman0 Mac:4a:cc:92:c1:04:a0 Sandbox:} {Name:vethe5c4c9f3 Mac:da:a7:ab:74:eb:00 Sandbox:} {Name:eth0 Mac:5a:b5:04:0f:74:3c Sandbox:/var/run/netns/cni-d26f2fff-2b88-e687-b24f-c40f73b772bd}], IP:[{Version:4 Interface:0xc00039cb48 Address:{IP:10.88.0.13 Mask:ffff0000} Gateway:10.88.0.1}], Routes:[{Dst:{IP:0.0.0.0 Mask:00000000} GW:<nil>}], DNS:{Nameservers:[] Domain: Search:[] Options:[]} 
DEBU[0000] /etc/system-fips does not exist on host, not mounting FIPS mode secret 
DEBU[0000] Setting CGroups for container 7477a79b5cf4aa1760f4138c20742f218d991de0e49602a480eac2f64d8d968b to machine.slice:libpod:7477a79b5cf4aa1760f4138c20742f218d991de0e49602a480eac2f64d8d968b 
DEBU[0000] reading hooks from /usr/share/containers/oci/hooks.d 
DEBU[0000] reading hooks from /etc/containers/oci/hooks.d 
DEBU[0000] Created OCI spec for container 7477a79b5cf4aa1760f4138c20742f218d991de0e49602a480eac2f64d8d968b at /var/lib/containers/storage/overlay-containers/7477a79b5cf4aa1760f4138c20742f218d991de0e49602a480eac2f64d8d968b/userdata/config.json 
DEBU[0000] /usr/bin/conmon messages will be logged to syslog 
DEBU[0000] running conmon: /usr/bin/conmon               args="[--api-version 1 -s -c 7477a79b5cf4aa1760f4138c20742f218d991de0e49602a480eac2f64d8d968b -u 7477a79b5cf4aa1760f4138c20742f218d991de0e49602a480eac2f64d8d968b -r /usr/bin/runc -b /var/lib/containers/storage/overlay-containers/7477a79b5cf4aa1760f4138c20742f218d991de0e49602a480eac2f64d8d968b/userdata -p /var/run/containers/storage/overlay-containers/7477a79b5cf4aa1760f4138c20742f218d991de0e49602a480eac2f64d8d968b/userdata/pidfile -l k8s-file:/var/lib/containers/storage/overlay-containers/7477a79b5cf4aa1760f4138c20742f218d991de0e49602a480eac2f64d8d968b/userdata/ctr.log --exit-dir /var/run/libpod/exits --socket-dir-path /var/run/libpod/socket --log-level debug --syslog -i --conmon-pidfile /var/run/containers/storage/overlay-containers/7477a79b5cf4aa1760f4138c20742f218d991de0e49602a480eac2f64d8d968b/userdata/conmon.pid --exit-command /usr/bin/podman --exit-command-arg --root --exit-command-arg /var/lib/containers/storage --exit-command-arg --runroot --exit-command-arg /var/run/containers/storage --exit-command-arg --log-level --exit-command-arg debug --exit-command-arg --cgroup-manager --exit-command-arg systemd --exit-command-arg --tmpdir --exit-command-arg /var/run/libpod --exit-command-arg --runtime --exit-command-arg runc --exit-command-arg --storage-driver --exit-command-arg overlay --exit-command-arg --events-backend --exit-command-arg journald --exit-command-arg container --exit-command-arg cleanup --exit-command-arg --rm --exit-command-arg 7477a79b5cf4aa1760f4138c20742f218d991de0e49602a480eac2f64d8d968b]"
INFO[0000] Running conmon under slice machine.slice and unitName libpod-conmon-7477a79b5cf4aa1760f4138c20742f218d991de0e49602a480eac2f64d8d968b.scope 
DEBU[0000] Received: 2999340                            
INFO[0000] Got Conmon PID as 2999328                    
DEBU[0000] Created container 7477a79b5cf4aa1760f4138c20742f218d991de0e49602a480eac2f64d8d968b in OCI runtime 
DEBU[0000] Attaching to container 7477a79b5cf4aa1760f4138c20742f218d991de0e49602a480eac2f64d8d968b 
DEBU[0000] connecting to socket /var/run/libpod/socket/7477a79b5cf4aa1760f4138c20742f218d991de0e49602a480eac2f64d8d968b/attach 
DEBU[0000] Starting container 7477a79b5cf4aa1760f4138c20742f218d991de0e49602a480eac2f64d8d968b with command [ls -l] 
DEBU[0000] Started container 7477a79b5cf4aa1760f4138c20742f218d991de0e49602a480eac2f64d8d968b 
DEBU[0000] Enabling signal proxying                     
DEBU[0000] Cleaning up container 7477a79b5cf4aa1760f4138c20742f218d991de0e49602a480eac2f64d8d968b 
DEBU[0000] Tearing down network namespace at /var/run/netns/cni-d26f2fff-2b88-e687-b24f-c40f73b772bd for container 7477a79b5cf4aa1760f4138c20742f218d991de0e49602a480eac2f64d8d968b 
INFO[0000] Got pod network &{Name:freeipa-master Namespace:freeipa-master ID:7477a79b5cf4aa1760f4138c20742f218d991de0e49602a480eac2f64d8d968b NetNS:/var/run/netns/cni-d26f2fff-2b88-e687-b24f-c40f73b772bd Networks:[] RuntimeConfig:map[podman:{IP: PortMappings:[] Bandwidth:<nil> IpRanges:[]}]} 
INFO[0000] About to del CNI network podman (type=bridge) 
DEBU[0000] unmounted container "7477a79b5cf4aa1760f4138c20742f218d991de0e49602a480eac2f64d8d968b" 
DEBU[0000] Container 7477a79b5cf4aa1760f4138c20742f218d991de0e49602a480eac2f64d8d968b storage is already unmounted, skipping... 
DEBU[0001] Removed volume 28f3c352f7aeef62b441fd31f3a38931edfb8f2dd8e4111324eb149512e250ab 

```

