# SCC and capabilities for initContainer init-volume

With this proof of concept is investigated the SCC profile needed to allow
the init-volume initContainer to be executed completely, with no issues
about permissions that could avoid the initialisation of the `/data` volume.

Along this trip different error messages were found and fixed.

Thanks to Matt Dorn and Alexandre Menezes for helping out to get this working
and letting me get a better understunding about Security Container Constraints.
Higly recommended to read [this article](https://www.openshift.com/blog/managing-sccs-in-openshift) written by Alexandre.

## Quick start

1. Create a new namespace (just once):

   ```shell
   oc new-project my-cool-namespace
   ```

1. Modify internal namespace ocurrencies by (just once):

   ```shell
   sed -i s/avisiedo-init-container/my-cool-namespace/g poc-003-admin.yaml
   ```

1. Delete, deploy and get-info about the proof of concept by:

   ```shell
   make app-delete app-deploy get-info
   ```

## Errors found

Errors found by trial and error while investigating the permissions for running
the initContainer that initialise the /data volume in OpenShift, and following
the systemd container interface document.

### tar: var/lib/chrony: Cannot change ownership to uid 285, gid 283: Operation not permitted

This error is happening in the context of init-data script:

```raw
+ /usr/local/bin/populate-volume-from-template /tmp
tar: var/lib/chrony: Cannot change ownership to uid 285, gid 283: Operation not permitted
tar: Exiting with failure status due to previous errors
```

Fixed by adding the `CHOWN` capability.

```raw
Make arbitrary changes to file UIDs and GIDs (see chown(2)).
```

Quick clarification here. Openshift clusters use chronyd at each node to
synchronize clocks between the different nodes in the cluster. The below
can be seen in one of the nodes:

```raw
sh-4.4# systemctl status chronyd
● chronyd.service - NTP client/server
   Loaded: loaded (/usr/lib/systemd/system/chronyd.service; enabled; vendor preset: enabled)
   Active: active (running) since Thu 2020-10-08 09:00:53 UTC; 1 weeks 5 days ago
     Docs: man:chronyd(8)
           man:chrony.conf(5)
  Process: 1215 ExecStartPost=/usr/libexec/chrony-helper update-daemon (code=exited, status=0/SUCCESS)
  Process: 1206 ExecStart=/usr/sbin/chronyd $OPTIONS (code=exited, status=0/SUCCESS)
 Main PID: 1209 (chronyd)
    Tasks: 1
   Memory: 2.8M
      CPU: 3.253s
   CGroup: /system.slice/chronyd.service
           └─1209 /usr/sbin/chronyd

Oct 20 16:06:09 permanent-bdd7p-worker-9r4b6 chronyd[1209]: Source 71.114.67.173 replaced with 2606:4700:f1::123
Oct 20 16:40:55 permanent-bdd7p-worker-9r4b6 chronyd[1209]: Source 2001:550:3402:3:8733:2637:dfce:2d18 replaced with 2604:e8c0:3::5
Oct 20 17:14:39 permanent-bdd7p-worker-9r4b6 chronyd[1209]: Source 74.6.168.72 replaced with 64.142.54.12
Oct 20 18:00:56 permanent-bdd7p-worker-9r4b6 chronyd[1209]: Source 2001:470:1f06:50::2 replaced with 2001:470:88f6::
Oct 20 18:41:03 permanent-bdd7p-worker-9r4b6 chronyd[1209]: Source 2606:4700:f1::123 replaced with 2001:19f0:8001:afd:5400:1ff:fe9d:cba
Oct 20 19:16:14 permanent-bdd7p-worker-9r4b6 chronyd[1209]: Source 2604:e8c0:3::5 replaced with 216.230.228.242
Oct 20 19:49:38 permanent-bdd7p-worker-9r4b6 chronyd[1209]: Source 64.142.54.12 replaced with 2001:470:0:50::2
Oct 20 20:36:01 permanent-bdd7p-worker-9r4b6 chronyd[1209]: Source 2001:470:88f6:: replaced with 108.59.2.24
Oct 20 21:15:48 permanent-bdd7p-worker-9r4b6 chronyd[1209]: Source 2001:19f0:8001:afd:5400:1ff:fe9d:cba replaced with 216.218.254.202
Oct 20 21:51:36 permanent-bdd7p-worker-9r4b6 chronyd[1209]: Source 216.230.228.242 replaced with 74.6.168.72
```

### tar: var/lib/chrony: Cannot change mode to rwxr-xr-x: Operation not permitted

This error is happening in the context of init-data script:

```raw
+ /usr/local/bin/populate-volume-from-template /tmp
tar: var/lib/chrony: Cannot change mode to rwxr-xr-x: Operation not permitted
tar: Exiting with failure status due to previous errors
```

Fixed by adding the `FOWNER` capability.

```raw
* Bypass permission checks on operations that normally require the filesystem UID of the process to match the UID of the file (e.g., chmod(2), utime(2)), excluding  those  operations  covered  by  CAP_DAC_OVERRIDE  and  CAP_DAC_READ_SEARCH;
* set inode flags (see ioctl_iflags(2)) on arbitrary files;
* set Access Control Lists (ACLs) on arbitrary files;
* ignore directory sticky bit on file deletion;
* modify user extended attributes on sticky directory owned by any user;
* specify O_NOATIME for arbitrary files in open(2) and fcntl(2).
```

### rm: cannot remove '/run/rpcbind/*': Permission denied

Fixed by adding the `DAC_OVERRIDE` capability.

```raw
Bypass file read, write, and execute permission checks.  (DAC is an abbreviation of "discretionary access control".)
```

### PR_SET_MM_ARG_START failed, proceeding without: Operation not permitted

Error on initing systemd, at the very beginning.

Fixed by adding the `SYS_RESOURCE` capability.

```raw
* Use reserved space on ext2 filesystems;
* make ioctl(2) calls controlling ext3 journaling;
* override disk quota limits;
* increase resource limits (see setrlimit(2));
* override RLIMIT_NPROC resource limit;
* override maximum number of consoles on console allocation;
* override maximum number of keymaps;
* allow more than 64hz interrupts from the real-time clock;
* raise msg_qbytes limit for a System V message queue above the limit in /proc/sys/kernel/msgmnb (see msgop(2) and msgctl(2));
* allow the RLIMIT_NOFILE resource limit on the number of "in-flight" file descriptors to be bypassed when passing file descriptors to another process via a UNIX domain socket (see unix(7));
* override the /proc/sys/fs/pipe-size-max limit when setting the capacity of a pipe using the F_SETPIPE_SZ fcntl(2) command.
* use F_SETPIPE_SZ to increase the capacity of a pipe above the limit specified by /proc/sys/fs/pipe-max-size;
* override /proc/sys/fs/mqueue/queues_max limit when creating POSIX message queues (see mq_overview(7));
* employ the prctl(2) PR_SET_MM operation;
```

### systemd-tmpfiles-setup.service: Failed to set invocation ID on control group /kubepods.slice/kubepods-besteffort.slice/kubepods-besteffort-podb5886319_457f_4071_b775_834cd8b166f0.slice/crio-11b9d211f76f0b2d0b4de4effd7126fcfac4c428359d2148671a783d11f3c803.scope/system.slice/systemd-tmpfiles-setup.service, ignoring: Operation not permitted

Following the breadcrumbs I can see that the error message happen when
using setxattr on the path indicated:

- https://github.com/systemd/systemd/blob/69c0807432fa4fbfbf507a53872664cd26715559/src/basic/cgroup-util.c#L583

  ```c
  int cg_set_xattr(const char *controller, const char *path, const char *name, const void *value, size_t size, int flags) {
          _cleanup_free_ char *fs = NULL;
          int r;
  
          assert(path);
          assert(name);
          assert(value || size <= 0);
  
          r = cg_get_path(controller, path, NULL, &fs);
          if (r < 0)
                  return r;
  
          if (setxattr(fs, name, value, size, flags) < 0)
                  return -errno;
  
          return 0;
  }
  ```

Searching at `man capabilities` can be found the below:

```raw
File capabilities
       Since kernel 2.6.24, the kernel supports associating capability sets with an executable file using setcap(8).  The file capability sets are stored in an extended attribute (see setxattr(2) and xattr(7))  named  secu‐
       rity.capability.  Writing to this extended attribute requires the CAP_SETFCAP capability.  The file capability sets, in conjunction with the capability sets of the thread, determine the capabilities of a thread after
       an execve(2).

       The three file capability sets are:

       Permitted (formerly known as forced):
              These capabilities are automatically permitted to the thread, regardless of the thread's inheritable capabilities.

       Inheritable (formerly known as allowed):
              This set is ANDed with the thread's inheritable set to determine which inheritable capabilities are enabled in the permitted set of the thread after the execve(2).

       Effective:
              This is not a set, but rather just a single bit.  If this bit is set, then during an execve(2) all of the new permitted capabilities for the thread are also raised in the effective set.  If  this  bit  is  not
              set, then after an execve(2), none of the new permitted capabilities is in the new effective set.

              Enabling  the  file effective capability bit implies that any file permitted or inheritable capability that causes a thread to acquire the corresponding permitted capability during an execve(2) (see the trans‐
              formation rules described below) will also acquire that capability in its effective set.  Therefore, when assigning capabilities to a file (setcap(8), cap_set_file(3), cap_set_fd(3)), if we specify the  effec‐
              tive flag as being enabled for any capability, then the effective flag must also be specified as enabled for all other capabilities for which the corresponding permitted or inheritable flags is enabled.
```

It should be fixed by `SETFCAP`, but still getting the error message, and
it is only removed when the `SYS_ADMIN` capability is added. This is the big
one to try to remove as this provide too many privileges. This will be
investigated into a separate investigation.

```raw
* Perform a range of system administration operations including: quotactl(2), mount(2), umount(2), pivot_root(2), swapon(2), swapoff(2), sethostname(2), and setdomainname(2);
* perform privileged syslog(2) operations (since Linux 2.6.37, CAP_SYSLOG should be used to permit such operations);
* perform VM86_REQUEST_IRQ vm86(2) command;
* perform IPC_SET and IPC_RMID operations on arbitrary System V IPC objects;
* override RLIMIT_NPROC resource limit;
* perform operations on trusted and security Extended Attributes (see xattr(7));
* use lookup_dcookie(2);
* use ioprio_set(2) to assign IOPRIO_CLASS_RT and (before Linux 2.6.25) IOPRIO_CLASS_IDLE I/O scheduling classes;
* forge PID when passing socket credentials via UNIX domain sockets;
* exceed /proc/sys/fs/file-max, the system-wide limit on the number of open files, in system calls that open files (e.g., accept(2), execve(2), open(2), pipe(2));
* employ CLONE_* flags that create new namespaces with clone(2) and unshare(2) (but, since Linux 3.8, creating user namespaces does not require any capability);
* call perf_event_open(2);
* access privileged perf event information;
* call setns(2) (requires CAP_SYS_ADMIN in the target namespace);
* call fanotify_init(2);
* call bpf(2);
* perform privileged KEYCTL_CHOWN and KEYCTL_SETPERM keyctl(2) operations;
* perform madvise(2) MADV_HWPOISON operation;
* employ the TIOCSTI ioctl(2) to insert characters into the input queue of a terminal other than the caller's controlling terminal;
* employ the obsolete nfsservctl(2) system call;
* employ the obsolete bdflush(2) system call;
* perform various privileged block-device ioctl(2) operations;
* perform various privileged filesystem ioctl(2) operations;
* perform privileged ioctl(2) operations on the /dev/random device (see random(4));
* install a seccomp(2) filter without first having to set the no_new_privs thread attribute;
* modify allow/deny rules for device control groups;
* employ the ptrace(2) PTRACE_SECCOMP_GET_FILTER operation to dump tracee's seccomp filters;
* employ the ptrace(2) PTRACE_SETOPTIONS operation to suspend the tracee's seccomp protections (i.e., the PTRACE_O_SUSPEND_SECCOMP flag);
* perform administrative operations on many device drivers.
* Modify autogroup nice values by writing to /proc/[pid]/autogroup (see sched(7)).
```

### Can't acquire effective CAP_SETPCAP bit, ignoring: Operation not permitted

This happens on this context, when launching `systemd-journald`.

```raw
Can't acquire effective CAP_SETPCAP bit, ignoring: Operation not permitted
systemd-journald.service: Failed to drop capabilities: Operation not permitted
systemd-journald.service: Failed at step CAPABILITIES spawning /usr/lib/systemd/systemd-journald: Operation not permitted
```

Fixed by adding `SETPCAP` capability. From `man capabilities` can be read:

```raw
       CAP_SETPCAP
              If  file  capabilities are supported (i.e., since Linux 2.6.24): add any capability from the calling thread's bounding set to its inheritable set; drop capabilities from the bounding set (via prctl(2) PR_CAPB‐
              SET_DROP); make changes to the securebits flags.

              If file capabilities are not supported (i.e., kernels before Linux 2.6.24): grant or remove any capability in the caller's permitted capability set to or from any other process.  (This property of  CAP_SETPCAP
              is not available when the kernel is configured to support file capabilities, since CAP_SETPCAP has entirely different semantics for such kernels.)
```

### Failed to create control group scope: Permission denied

```raw
+ exec /usr/sbin/original/init --show-status=false
systemd 239 running in system mode. (+PAM +AUDIT +SELINUX +IMA -APPARMOR +SMACK +SYSVINIT +UTMP +LIBCRYPTSETUP +GCRYPT +GNUTLS +ACL +XZ +LZ4 +SECCOMP +BLKID +ELFUTILS +KMOD +IDN2 -IDN +PCRE2 default-hierarchy=legacy)
Detected virtualization container-other.
Detected architecture x86-64.
Set hostname to <poc-003>.
Initializing machine ID from container UUID.
Failed to mount /etc/machine-id: Operation not permitted
Failed to add address 127.0.0.1 to loopback interface: Operation not permitted
Failed to add address ::1 to loopback interface: Operation not permitted
Failed to bring loopback interface up: Operation not permitted
Failed to bump AF_UNIX datagram queue length, ignoring: Read-only file system
Found cgroup on /sys/fs/cgroup/systemd, legacy hierarchy
Using cgroup controller name=systemd. File system hierarchy is at /sys/fs/cgroup/systemd/kubepods.slice/kubepods-besteffort.slice/kubepods-besteffort-poda50f5782_9931_42f6_8e88_9f334e3a2838.slice/crio-d169ae56934ebe6d5659d080aa0bc3e88ea0dd8cbaef9e1aea0296554a546d57.scope.
Release agent already installed.
Failed to create /kubepods.slice/kubepods-besteffort.slice/kubepods-besteffort-poda50f5782_9931_42f6_8e88_9f334e3a2838.slice/crio-d169ae56934ebe6d5659d080aa0bc3e88ea0dd8cbaef9e1aea0296554a546d57.scope/init.scope control group: Permission denied
Failed to allocate manager object: Permission denied
Failed to read pids.max attribute of cgroup root, ignoring: No such file or directory
[!!!!!!] Failed to allocate manager object, freezing.
```

More specific:

```raw
Failed to create /kubepods.slice/kubepods-besteffort.slice/kubepods-besteffort-poda50f5782_9931_42f6_8e88_9f334e3a2838.slice/crio-d169ae56934ebe6d5659d080aa0bc3e88ea0dd8cbaef9e1aea0296554a546d57.scope/init.scope control group: Permission denied
```

Fixed by:

- Setting `allowHostDirVolumePlugin: true`.
- Enabling sebool `container_manage_cgroup`.

> It has been seen the worker nodes had `container_manage_cgroup` disabled.

### krb5kdc.service: Can't open PID file /var/run/krb5kdc.pid (yet?) after start: No such file or directory

Fixed by mounting /var/run as an in-memory emptyDir with rw permissions.

### dirsrv@APPS-PERMANENT-IDMOCP-LAB-ENG-RDU2-REDHAT-COM.service: Can't open PID file /run/dirsrv/slapd-APPS-PERMANENT-IDMOCP-LAB-ENG-RDU2-REDHAT-COM.pid (yet?) after start: No such file or directory

Fixed by mounting /var/run/dirsrv as an in-memory emptyDir with rw permissions.

### The ipa-server-install command failed, exception: CalledProcessError: Command '['systemctl', 'start', 'dirsrv@APPS-PERMANENT-IDMOCP-LAB-ENG-RDU2-REDHAT-COM']' returned non-zero exit status 1

Fixed by adding the capabilities: `SETUID`, `SETGID`, `KILL` and `NET_BIND_SERVICE`.

- SETUID

  ```raw
  * Make arbitrary manipulations of process UIDs (setuid(2), setreuid(2), setresuid(2), setfsuid(2));
  * forge UID when passing socket credentials via UNIX domain sockets;
  * write a user ID mapping in a user namespace (see user_namespaces(7)).
  ```

- SETGID

  ```raw
  - Make arbitrary manipulations of process GIDs and supplementary GID list;
  - forge GID when passing socket credentials via UNIX domain sockets;
  - write a group ID mapping in a user namespace (see user_namespaces(7)).
  ```

- KILL

  ```raw
  Bypass permission checks for sending signals (see kill(2)).  This includes use of the ioctl(2) KDSIGACCEPT operation.
  ```

- NET_BIND_SERVICE

  ```raw
  Bind a socket to Internet domain privileged ports (port numbers less than 1024).
  ```

### PR_FILE_NOT_FOUND_ERROR: File not found

```raw
Starting external process
args=['/usr/bin/certutil', '-d', 'sql:/etc/dirsrv/slapd-APPS-PERMANENT-IDMOCP-LAB-ENG-RDU2-REDHAT-COM/', '-L', '-n', 'APPS.PERMANENT.IDMOCP.LAB.ENG.RDU2.REDHAT.COM IPA CA', '-a', '-f', '/etc/dirsrv/slapd-APPS-PERMANENT-IDMOCP-LAB-ENG-RDU2
-REDHAT-COM/pwdfile.txt']
Process finished, return code=255
stdout=
stderr=certutil: Could not find cert: APPS.PERMANENT.IDMOCP.LAB.ENG.RDU2.REDHAT.COM IPA CA
: PR_FILE_NOT_FOUND_ERROR: File not found

Starting external process
args=['/usr/bin/certutil', '-d', 'sql:/etc/dirsrv/slapd-APPS-PERMANENT-IDMOCP-LAB-ENG-RDU2-REDHAT-COM/', '-N', '-f', '/etc/dirsrv/slapd-APPS-PERMANENT-IDMOCP-LAB-ENG-RDU2-REDHAT-COM/pwdfile.txt', '-@', '/etc/dirsrv/slapd-APPS-PERMANENT-ID
MOCP-LAB-ENG-RDU2-REDHAT-COM/pwdfile.txt']
Process finished, return code=0
stdout=
stderr=
```

**How is this impacting the install process?**

Thanks @frasertweedale and @tiran.

> **PR_FILE_NOT_FOUND_ERROR** error is expected and harmless. It is a normal
> part of 389DS NSSDB initialisation. The error symbol is comming from
> [Netscape Portable Runtime](https://developer.mozilla.org/es/docs/NSPR)
> library.

### No valid Negotiate header in server response

> I have not detected what fixed this, but the current state install
> ipa-client with no issues. I leave this just in case it is replayed.

```raw
Starting external process
args=['/usr/sbin/ipa-client-install', '--on-master', '--unattended', '--domain', 'apps.permanent.idmocp.lab.eng.rdu2.redhat.com', '--server', 'poc-003.apps.permanent.idmocp.lab.eng.rdu2.redhat.com', '--realm', 'APPS.PERMANENT.IDMOCP.LAB.ENG.RDU2.REDHAT.COM', '--hostname', 'poc-003.apps.permanent.idmocp.lab.eng.rdu2.redhat.com', '--no-ntp', '--no-ssh', '--no-sshd']
This program will set up IPA client.
Version 4.8.4

Using existing certificate '/etc/ipa/ca.crt'.
Client hostname: poc-003.apps.permanent.idmocp.lab.eng.rdu2.redhat.com
Realm: APPS.PERMANENT.IDMOCP.LAB.ENG.RDU2.REDHAT.COM
DNS Domain: apps.permanent.idmocp.lab.eng.rdu2.redhat.com
IPA Server: poc-003.apps.permanent.idmocp.lab.eng.rdu2.redhat.com
BaseDN: dc=apps,dc=permanent,dc=idmocp,dc=lab,dc=eng,dc=rdu2,dc=redhat,dc=com

Configured sudoers in /data/etc/authselect/user-nsswitch.conf
Configured /etc/sssd/sssd.conf
No valid Negotiate header in server response
The ipa-client-install command failed. See /var/log/ipaclient-install.log for more information
httpd.service: Got notification message from PID 3065 (READY=1, STATUS=Total requests: 1; Idle/Busy workers 100/0;Requests/sec: 0.0208; Bytes served/sec: 203 B/sec)
Process finished, return code=1
```

Executed command:

```shell
/usr/sbin/ipa-client-install --on-master \
                             --unattended \
                             --domain apps.permanent.idmocp.lab.eng.rdu2.redhat.com \
                             --server poc-003.apps.permanent.idmocp.lab.eng.rdu2.redhat.com \
                             --realm APPS.PERMANENT.IDMOCP.LAB.ENG.RDU2.REDHAT.COM \
                             --hostname poc-003.apps.permanent.idmocp.lab.eng.rdu2.redhat.com \
                             --no-ntp \
                             --no-ssh \
                             --no-sshd
```

More information from ipaclient-install.log file:

```raw
2020-10-14T20:56:37Z DEBUG failed to find session_cookie in persistent storage for principal 'host/poc-003.apps.permanent.idmocp.lab.eng.rdu2.redhat.com@APPS.PERMANENT.IDMOCP.LAB.ENG.RDU2.REDHAT.COM'
2020-10-14T20:56:37Z DEBUG trying https://poc-003.apps.permanent.idmocp.lab.eng.rdu2.redhat.com/ipa/json
2020-10-14T20:56:37Z DEBUG Created connection context.rpcclient_140622652047312
2020-10-14T20:56:37Z DEBUG [try 1]: Forwarding 'schema' to json server 'https://poc-003.apps.permanent.idmocp.lab.eng.rdu2.redhat.com/ipa/json'
2020-10-14T20:56:37Z DEBUG New HTTP connection (poc-003.apps.permanent.idmocp.lab.eng.rdu2.redhat.com)
2020-10-14T20:56:37Z DEBUG HTTP connection destroyed (poc-003.apps.permanent.idmocp.lab.eng.rdu2.redhat.com)
Traceback (most recent call last):
  File "/usr/lib/python3.6/site-packages/ipaclient/remote_plugins/__init__.py", line 126, in get_package
    plugins = api._remote_plugins
AttributeError: 'API' object has no attribute '_remote_plugins'

During handling of the above exception, another exception occurred:

Traceback (most recent call last):
  File "/usr/lib/python3.6/site-packages/ipalib/rpc.py", line 724, in single_request
    if not self._auth_complete(response):
  File "/usr/lib/python3.6/site-packages/ipalib/rpc.py", line 677, in _auth_complete
    message=u"No valid Negotiate header in server response")
ipalib.errors.KerberosError: No valid Negotiate header in server response
2020-10-14T20:56:37Z DEBUG Destroyed connection context.rpcclient_140622652047312
2020-10-14T20:56:37Z DEBUG   File "/usr/lib/python3.6/site-packages/ipapython/admintool.py", line 179, in execute
    return_value = self.run()
  File "/usr/lib/python3.6/site-packages/ipapython/install/cli.py", line 340, in run
    return cfgr.run()
  File "/usr/lib/python3.6/site-packages/ipapython/install/core.py", line 360, in run
    return self.execute()
  File "/usr/lib/python3.6/site-packages/ipapython/install/core.py", line 386, in execute
    for rval in self._executor():
  File "/usr/lib/python3.6/site-packages/ipapython/install/core.py", line 431, in __runner
    exc_handler(exc_info)
  File "/usr/lib/python3.6/site-packages/ipapython/install/core.py", line 460, in _handle_execute_exception
    self._handle_exception(exc_info)
  File "/usr/lib/python3.6/site-packages/ipapython/install/core.py", line 450, in _handle_exception
    six.reraise(*exc_info)
  File "/usr/lib/python3.6/site-packages/six.py", line 693, in reraise
    raise value
  File "/usr/lib/python3.6/site-packages/ipapython/install/core.py", line 421, in __runner
    step()
  File "/usr/lib/python3.6/site-packages/ipapython/install/core.py", line 418, in <lambda>
    step = lambda: next(self.__gen)
  File "/usr/lib/python3.6/site-packages/ipapython/install/util.py", line 81, in run_generator_with_yield_from
    six.reraise(*exc_info)
  File "/usr/lib/python3.6/site-packages/six.py", line 693, in reraise
    raise value
  File "/usr/lib/python3.6/site-packages/ipapython/install/util.py", line 59, in run_generator_with_yield_from
    value = gen.send(prev_value)
  File "/usr/lib/python3.6/site-packages/ipapython/install/core.py", line 655, in _configure
    next(executor)
  File "/usr/lib/python3.6/site-packages/ipapython/install/core.py", line 431, in __runner
    exc_handler(exc_info)
  File "/usr/lib/python3.6/site-packages/ipapython/install/core.py", line 460, in _handle_execute_exception
    self._handle_exception(exc_info)
  File "/usr/lib/python3.6/site-packages/ipapython/install/core.py", line 518, in _handle_exception
    self.__parent._handle_exception(exc_info)
  File "/usr/lib/python3.6/site-packages/ipapython/install/core.py", line 450, in _handle_exception
    six.reraise(*exc_info)
  File "/usr/lib/python3.6/site-packages/six.py", line 693, in reraise
    raise value
  File "/usr/lib/python3.6/site-packages/ipapython/install/core.py", line 515, in _handle_exception
    super(ComponentBase, self)._handle_exception(exc_info)
  File "/usr/lib/python3.6/site-packages/ipapython/install/core.py", line 450, in _handle_exception
    six.reraise(*exc_info)
  File "/usr/lib/python3.6/site-packages/six.py", line 693, in reraise
    raise value
  File "/usr/lib/python3.6/site-packages/ipapython/install/core.py", line 421, in __runner
    step()
  File "/usr/lib/python3.6/site-packages/ipapython/install/core.py", line 418, in <lambda>
    step = lambda: next(self.__gen)
  File "/usr/lib/python3.6/site-packages/ipapython/install/util.py", line 81, in run_generator_with_yield_from
    six.reraise(*exc_info)
  File "/usr/lib/python3.6/site-packages/six.py", line 693, in reraise
    raise value
  File "/usr/lib/python3.6/site-packages/ipapython/install/util.py", line 59, in run_generator_with_yield_from
    value = gen.send(prev_value)
  File "/usr/lib/python3.6/site-packages/ipapython/install/common.py", line 65, in _install
    for unused in self._installer(self.parent):
  File "/usr/lib/python3.6/site-packages/ipaclient/install/client.py", line 3818, in main
    install(self)
  File "/usr/lib/python3.6/site-packages/ipaclient/install/client.py", line 2531, in install
    _install(options)
  File "/usr/lib/python3.6/site-packages/ipaclient/install/client.py", line 2843, in _install
    api.finalize()
  File "/usr/lib/python3.6/site-packages/ipalib/plugable.py", line 743, in finalize
    self.__do_if_not_done('load_plugins')
  File "/usr/lib/python3.6/site-packages/ipalib/plugable.py", line 430, in __do_if_not_done
    getattr(self, name)()
  File "/usr/lib/python3.6/site-packages/ipalib/plugable.py", line 622, in load_plugins
    for package in self.packages:
  File "/usr/lib/python3.6/site-packages/ipalib/__init__.py", line 954, in packages
    ipaclient.remote_plugins.get_package(self),
  File "/usr/lib/python3.6/site-packages/ipaclient/remote_plugins/__init__.py", line 134, in get_package
    plugins = schema.get_package(server_info, client)
  File "/usr/lib/python3.6/site-packages/ipaclient/remote_plugins/schema.py", line 553, in get_package
    schema = Schema(client)
  File "/usr/lib/python3.6/site-packages/ipaclient/remote_plugins/schema.py", line 402, in __init__
    fingerprint, ttl = self._fetch(client, ignore_cache=read_failed)
  File "/usr/lib/python3.6/site-packages/ipaclient/remote_plugins/schema.py", line 427, in _fetch
    schema = client.forward(u'schema', **kwargs)['result']
  File "/usr/lib/python3.6/site-packages/ipalib/rpc.py", line 1149, in forward
    return self._call_command(command, params)
  File "/usr/lib/python3.6/site-packages/ipalib/rpc.py", line 1125, in _call_command
    return command(*params)
  File "/usr/lib/python3.6/site-packages/ipalib/rpc.py", line 1279, in _call
    return self.__request(name, args)
  File "/usr/lib/python3.6/site-packages/ipalib/rpc.py", line 1246, in __request
    verbose=self.__verbose >= 3,
  File "/usr/lib64/python3.6/xmlrpc/client.py", line 1154, in request
    return self.single_request(host, handler, request_body, verbose)
  File "/usr/lib/python3.6/site-packages/ipalib/rpc.py", line 724, in single_request
    if not self._auth_complete(response):
  File "/usr/lib/python3.6/site-packages/ipalib/rpc.py", line 677, in _auth_complete
    message=u"No valid Negotiate header in server response")

2020-10-14T20:56:37Z DEBUG The ipa-client-install command failed, exception: KerberosError: No valid Negotiate header in server
2020-10-14T20:56:37Z ERROR No valid Negotiate header in server response
2020-10-14T20:56:37Z ERROR The ipa-client-install command failed. See /var/log/ipaclient-install.log for more information
```

## BPF Firewalling not supported on this manager, proceeding without

```raw
BPF firewalling not supported on this manager, proceeding without.
BPF firewalling not supported on this manager, proceeding without.
BPF firewalling not supported on this manager, proceeding without.
BPF firewalling not supported on this manager, proceeding without.
BPF firewalling not supported on this manager, proceeding without.
BPF firewalling not supported on this manager, proceeding without.
BPF firewalling not supported on this manager, proceeding without.
BPF firewalling not supported on this manager, proceeding without.
BPF firewalling not supported on this manager, proceeding without.
BPF firewalling not supported on this manager, proceeding without.
BPF firewalling not supported on this manager, proceeding without.
BPF firewalling not supported on this manager, proceeding without.
BPF firewalling not supported on this manager, proceeding without.
BPF firewalling not supported on this manager, proceeding without.
BPF firewalling not supported on this manager, proceeding without.
BPF firewalling not supported on this manager, proceeding without.
BPF firewalling not supported on this manager, proceeding without.
BPF firewalling not supported on this manager, proceeding without.
BPF firewalling not supported on this manager, proceeding without.
BPF firewalling not supported on this manager, proceeding without.
BPF firewalling not supported on this manager, proceeding without.
BPF firewalling not supported on this manager, proceeding without.
BPF firewalling not supported on this manager, proceeding without.
```

No blocking. Skipping this error message.

## Current process state when ipa-server-install finish

This is when running with `no-exit` argument; then I get the below from the
init-volume container.

```raw
[root@poc-003 /]# ps fax
    PID TTY      STAT   TIME COMMAND
   3879 pts/1    Ss     0:00 /bin/bash
   3893 pts/1    R+     0:00  \_ ps fax
      1 ?        Ss     0:02 /usr/sbin/original/init --show-status=false
    245 ?        S      0:00 /usr/bin/coreutils --coreutils-prog-shebang=tail /usr/bin/tail --silent -n 0 -f --retry /var/log/ipa-server-configure-first.log /var/log/ipa-server-run.log
    269 ?        Ss     0:00 /usr/bin/dbus-daemon --system --address=systemd: --nofork --nopidfile --systemd-activation --syslog-only
    639 ?        Ss     0:00 /usr/sbin/kadmind -P /var/run/kadmind.pid
    663 ?        Ss     0:00 /usr/libexec/platform-python -I /usr/libexec/ipa/ipa-custodia /etc/ipa/custodia/custodia.conf
   1699 ?        Ss     0:00 /usr/sbin/certmonger -S -p /var/run/certmonger.pid -n -d2
   2456 ?        Z      0:00  \_ [certmonger] <defunct>
   2798 ?        Ssl    0:18 /usr/lib/jvm/jre-1.8.0-openjdk/bin/java -classpath /usr/share/tomcat/bin/bootstrap.jar:/usr/share/tomcat/bin/tomcat-juli.jar:/usr/share/java/ant.jar:/usr/share/java/ant-launcher.jar:/usr/lib/jvm/java/lib/tools
   2984 ?        Ssl    0:00 /usr/sbin/gssproxy -D
   3064 ?        Ss     0:00 /usr/sbin/httpd -DFOREGROUND
   3070 ?        S      0:00  \_ /usr/sbin/httpd -DFOREGROUND
   3072 ?        Sl     0:00  \_ (wsgi:kdcproxy) -DFOREGROUND
   3073 ?        Sl     0:00  \_ (wsgi:kdcproxy) -DFOREGROUND
   3074 ?        Sl     0:02  \_ (wsgi:ipa)      -DFOREGROUND
   3075 ?        Sl     0:02  \_ (wsgi:ipa)      -DFOREGROUND
   3076 ?        Sl     0:02  \_ (wsgi:ipa)      -DFOREGROUND
   3077 ?        Sl     0:02  \_ (wsgi:ipa)      -DFOREGROUND
   3078 ?        Sl     0:00  \_ /usr/sbin/httpd -DFOREGROUND
   3079 ?        Sl     0:00  \_ /usr/sbin/httpd -DFOREGROUND
   3080 ?        Sl     0:00  \_ /usr/sbin/httpd -DFOREGROUND
   3629 ?        Sl     0:00  \_ /usr/sbin/httpd -DFOREGROUND
   3353 ?        Ss     0:00 /usr/sbin/oddjobd -n -p /var/run/oddjobd.pid -t 300
   3568 ?        Ssl    0:04 /usr/sbin/ns-slapd -D /etc/dirsrv/slapd-APPS-PERMANENT-IDMOCP-LAB-ENG-RDU2-REDHAT-COM -i /run/dirsrv/slapd-APPS-PERMANENT-IDMOCP-LAB-ENG-RDU2-REDHAT-COM.pid
   3607 ?        Ss     0:00 /usr/sbin/krb5kdc -P /var/run/krb5kdc.pid -w 4
   3608 ?        S      0:00  \_ /usr/sbin/krb5kdc -P /var/run/krb5kdc.pid -w 4
   3609 ?        S      0:00  \_ /usr/sbin/krb5kdc -P /var/run/krb5kdc.pid -w 4
   3610 ?        S      0:00  \_ /usr/sbin/krb5kdc -P /var/run/krb5kdc.pid -w 4
   3611 ?        S      0:00  \_ /usr/sbin/krb5kdc -P /var/run/krb5kdc.pid -w 4
   3715 ?        Ss     0:00 /usr/sbin/sssd -i --logger=files
   3716 ?        S      0:00  \_ /usr/libexec/sssd/sssd_be --domain implicit_files --uid 0 --gid 0 --logger=files
   3717 ?        S      0:00  \_ /usr/libexec/sssd/sssd_be --domain apps.permanent.idmocp.lab.eng.rdu2.redhat.com --uid 0 --gid 0 --logger=files
   3719 ?        S      0:00  \_ /usr/libexec/sssd/sssd_nss --uid 0 --gid 0 --logger=files
   3720 ?        S      0:00  \_ /usr/libexec/sssd/sssd_pam --uid 0 --gid 0 --logger=files
   3721 ?        S      0:00  \_ /usr/libexec/sssd/sssd_ifp --uid 0 --gid 0 --logger=files
   3722 ?        S      0:00  \_ /usr/libexec/sssd/sssd_sudo --uid 0 --gid 0 --logger=files
   3723 ?        S      0:00  \_ /usr/libexec/sssd/sssd_pac --uid 0 --gid 0 --logger=files
```

This is what I get in the main container, when I run with `exit-on-finished`
argument:

```raw
[root@poc-003 /]# ps fax
    PID TTY      STAT   TIME COMMAND
    254 pts/0    Ss     0:00 /bin/bash
    269 pts/0    R+     0:00  \_ ps fax
      1 ?        Ss     0:00 /usr/sbin/original/init
     24 ?        Ss     0:00 /usr/lib/systemd/systemd-journald
     32 ?        Ss     0:00 /usr/bin/dbus-daemon --system --address=systemd: --nofork --nopidfile --systemd-activation --syslog-only
     33 ?        Ss     0:00 /usr/sbin/oddjobd -n -p /var/run/oddjobd.pid -t 300
     34 ?        Ss     0:00 /usr/sbin/certmonger -S -p /var/run/certmonger.pid -n -d2
     36 ?        Ss     0:00 /usr/sbin/sssd -i --logger=files
     54 ?        S      0:00  \_ /usr/libexec/sssd/sssd_be --domain implicit_files --uid 0 --gid 0 --logger=files
     55 ?        S      0:00  \_ /usr/libexec/sssd/sssd_be --domain apps.permanent.idmocp.lab.eng.rdu2.redhat.com --uid 0 --gid 0 --logger=files
     56 ?        S      0:00  \_ /usr/libexec/sssd/sssd_nss --uid 0 --gid 0 --logger=files
     57 ?        S      0:00  \_ /usr/libexec/sssd/sssd_pam --uid 0 --gid 0 --logger=files
     58 ?        S      0:00  \_ /usr/libexec/sssd/sssd_ifp --uid 0 --gid 0 --logger=files
     59 ?        S      0:00  \_ /usr/libexec/sssd/sssd_sudo --uid 0 --gid 0 --logger=files
     60 ?        S      0:00  \_ /usr/libexec/sssd/sssd_pac --uid 0 --gid 0 --logger=files
     40 ?        Ssl    0:00 /usr/sbin/gssproxy -D
```

> The main container has not started FreeIPA as desired, but the scope of this
> investigation is the initContainers rather than the main container, so far.

## References

- [capabilities.7](https://man7.org/linux/man-pages/man7/capabilities.7.html).
- [Managing SCCs in OpenShift](https://www.openshift.com/blog/managing-sccs-in-openshift)
