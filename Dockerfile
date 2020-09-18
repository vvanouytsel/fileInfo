FROM centos:centos7 AS builder

RUN set -xe \
    && yum -y install git gcc glibc-static \
    && curl -sS -L -o go.tar.gz https://golang.org/dl/go1.14.8.linux-amd64.tar.gz \
    && tar -xf go.tar.gz -C /usr/local

ENV PATH="/usr/local/go/bin:${PATH}"

WORKDIR /src
COPY . /src/

RUN set -xe \
    && go build -o /bin/fileInfo


FROM centos:centos7

COPY . /src/
COPY --from=builder /bin/fileInfo /bin/

WORKDIR /src/
ENTRYPOINT [ "/bin/fileInfo" ]
