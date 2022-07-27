FROM alpine:latest
ENV GOTRACEBACK=crash TMPDIR=/tmp/ WORKDIR=/myworkdir/
WORKDIR ${WORKDIR}
COPY docker-entrypoint.sh /
COPY go_coredump_test ${TMPDIR}
ENTRYPOINT ["/docker-entrypoint.sh"]
