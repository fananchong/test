FROM alpine:latest
RUN apk --no-cache add bash
COPY ./race ./date.sh ./docker-entrypoint.sh /
ENTRYPOINT ["/docker-entrypoint.sh"]
