FROM alpine:latest
RUN apk --no-cache add bash
COPY ./race ./docker_entrypoint.sh /
ENTRYPOINT ["/docker_entrypoint.sh"]
