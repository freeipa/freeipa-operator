FROM docker.io/freeipa/freeipa-server:centos-8

ENV KRB5_TRACE=/dev/console

# Workaround for running systemd mode from podman
RUN mkdir /usr/sbin/original \
    && ln -svf /usr/lib/systemd/systemd /usr/sbin/original/init \
    && rm /usr/sbin/init
COPY init-wrapper.sh /usr/sbin/init
COPY init-data /usr/loca/sbin/init

# Enable debug
COPY gssproxy.conf /data-template/etc/gssproxy/gssproxy.conf

VOLUME /data

ENTRYPOINT ["/usr/sbin/init"]
