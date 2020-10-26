FROM docker.io/python:3.9.0-alpine


ENV ECHO_PROTOCOL=udp
ENV ECHO_ADDRESS=0.0.0.0
ENV ECHO_PORT=8007
ENV ECHO_MAX_CLIENTS=1


WORKDIR /usr/src/app

COPY . ./

ENTRYPOINT ["/usr/src/app/entrypoint.sh"]
HEALTHCHECK --interval=30s --timeout=30s --start-period=5s --retries=3 CMD [ "timeout", "3", "bash", "-c", "</dev/tcp/localhost/${ECHO_PORT}" ]
