FROM alpine:latest
RUN apk --no-cache add bash
COPY ./race ./json.sh ./docker-entrypoint.sh /
ENTRYPOINT ["/docker-entrypoint.sh"]
