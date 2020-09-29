FROM docker.io/freeipa/freeipa-server:centos-8

# Workaround for running systemd mode from podman
RUN mkdir /usr/sbin/original \
    && ln -svf /usr/lib/systemd/systemd /usr/sbin/original/init \
    && rm /usr/sbin/init
COPY init-wrapper.sh /usr/sbin/init
COPY init-data /usr/loca/sbin/init

# Workaround to don't loose the machine-id between different executions
# Pass on $container_uuid to set up the machine-id when the file content is
# empty (first initialization of the volume)
RUN rm -f /etc/machine-id \
    && ( [ -e /data-template/etc ] || mkdir -p /data-template/etc ) \
    && ( [ -e /data/etc ] || mkdir -p /data/etc ) \
    && touch /data/etc/machine-id \
    && chown 0444 /data-template/etc/machine-id \
    && ln -svf /data/etc/machine-id /etc/machine-id \
    && touch /data-template/etc/machine-id \
    && rm -vf /data/etc/machine-id

# Workaround for setting PrivateTmp=* to PrivateTmp=no
# https://blog.iwakd.de/lxc-cap_sys_admin-jessie
#
# Extend to all the paths below:
#   /etc/systemd/system.control
#   /run/systemd/system.control
#   /run/systemd/transient
#   /etc/systemd/system
#   /run/systemd/system
#   /run/systemd/generator
#   /usr/lib/systemd/system
# RUN DIRS="/etc/systemd/system.control /run/systemd/system.control /run/systemd/transient /etc/systemd/system /run/systemd/system /run/systemd/generator /usr/lib/systemd/system /usr/lib/systemd/portable"; \
#     for item in $( grep -HRl ^PrivateTmp= $DIRS ); do \
#       sed -i -e s/^PrivateTmp=yes/PrivateTmp=no/g -e s/^PrivateTmp=true/PrivateTmp=no/g ${item}; \
#     done; \
#     for item in $( grep -HRl ^PrivateDevices= $DIRS ); do \
#       sed -i -e s/^PrivateDevices=yes/PrivateDevices=no/g -e s/^PrivateDevices=true/PrivateDevices=no/g $item; \
#     done; \
#     for item in $( grep -HRl ^PrivateNetwork= $DIRS ); do \
#       sed -i -e s/^PrivateNetwork=yes/PrivateNetwork=no/g -e s/^PrivateNetwork=true/PrivateNetwork=no/g $item; \
#     done; \
#     for item in $( grep -HRl ^ProtectSystem= $DIRS ); do \
#       sed -i -e s/^ProtectSystem=strict/ProtectSystem=/g -e s/^ProtectSystem=full/ProtectSystem=/g $item; \
#     done; \
#     for item in $( grep -HRl ^ProtectHome= $DIRS ); do \
#       sed -i -e s/^ProtectHome=yes/ProtectHome=/g $item; \
#     done; \
#     for item in $( grep -HRl ^ReadOnlyDirectories= $DIRS ); do \
#       sed -i -e s/^ReadOnlyDirectories=strict/ReadOnlyDirectories=/g -e s/^ReadOnlyDirectories=full/ReadOnlyDirectories=/g $item; \
#     done; \
#     for item in $( grep -HRl -e ^SystemCallFilter= -e ^RestrictAddressFamilies= $DIRS ); do \
#       sed -i -e s/^.Service.\$/[Service]\\nNoNewPrivileges=yes/g $item; \
#     done

VOLUME /data

ENTRYPOINT ["/usr/sbin/init"]
