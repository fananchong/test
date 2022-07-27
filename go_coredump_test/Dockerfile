FROM alpine:latest
ARG TMPDIR=/tmp/ WORKDIR=/myworkdir/
ENV GOTRACEBACK=crash TMPDIR=${TMPDIR} WORKDIR=${WORKDIR}
WORKDIR ${WORKDIR}
COPY docker-entrypoint.sh /
COPY go_coredump_test ${TMPDIR}
ENTRYPOINT ["/docker-entrypoint.sh"]
