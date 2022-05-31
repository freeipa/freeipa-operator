# https://github.com/operator-framework/operator-sdk/blob/09c3aa14625965af9f22f513cd5c891471dbded2/images/custom-scorecard-tests/Dockerfile
FROM golang:1.17 AS builder

COPY . /src
WORKDIR /src
RUN go build -mod vendor -o bin/custom-scorecard-tests cmd/scorecard-tests/main.go

FROM registry.access.redhat.com/ubi8/ubi-minimal:8.6-751

ENV HOME=/opt/custom-scorecard-tests \
    USER_NAME=custom-scorecard-tests \
    USER_UID=1001

RUN echo "${USER_NAME}:x:${USER_UID}:0:${USER_NAME} user:${HOME}:/sbin/nologin" >> /etc/passwd

WORKDIR ${HOME}

ARG BIN=/src/bin/custom-scorecard-tests
COPY --from=builder $BIN /usr/local/bin/custom-scorecard-tests

ENTRYPOINT ["/usr/local/bin/custom-scorecard-tests"]

USER ${USER_UID}
