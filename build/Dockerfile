# https://catalog.redhat.com/software/containers/ubi8/ubi-minimal/5c359a62bed8bd75a2c3fba8?tag=8.2
FROM registry.access.redhat.com/ubi8/ubi-minimal:8.2

ENV OPERATOR=/usr/local/bin/idmocp \
    USER_UID=1001 \
    USER_NAME=idmocp

# install operator binary
COPY build/_output/bin/freeipa-operator ${OPERATOR}

COPY build/bin /usr/local/bin
RUN  /usr/local/bin/user_setup

ENTRYPOINT ["/usr/local/bin/entrypoint"]

USER ${USER_UID}
