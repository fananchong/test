FROM alpine:latest
RUN apk --no-cache add bash
COPY ./race ./docker-entrypoint.sh /
ENTRYPOINT ["/docker-entrypoint.sh"]
