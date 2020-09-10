# --- Build qdr-print
FROM centos:8 AS builder

RUN yum install epel-release -y && \
        yum update -y --setopt=tsflags=nodocs && \
        yum install git golang qpid-proton-c-devel --setopt=tsflags=nodocs -y && \
        yum clean all && go get qpid.apache.org/amqp qpid.apache.org/electron

ENV D=/home/qdr-print

WORKDIR $D
COPY . $D/

RUN     go build qdr-print.go && \
        mv qdr-print /tmp/

# --- end build, create qdr-print runtime layer ---
FROM centos:8

LABEL io.k8s.display-name="Simple QDR Printer" \
      io.k8s.description="Reads data from AMQP via proton and dumps to stdout."

RUN yum install epel-release -y && \
        yum update -y --setopt=tsflags=nodocs && \
        yum install qpid-proton-c jq --setopt=tsflasgs=nodocs -y && \
        yum clean all && \
        rm -rf /var/cache/yum

COPY --from=builder /tmp/qdr-print /

CMD ["/qdr-print"]
