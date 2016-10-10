FROM haproxy:1.6-alpine

RUN apk add --no-cache ca-certificates

ENV GOLANG_VERSION 1.7.1
ENV GOLANG_SRC_URL https://golang.org/dl/go$GOLANG_VERSION.src.tar.gz
ENV GOLANG_SRC_SHA256 2b843f133b81b7995f26d0cb64bbdbb9d0704b90c44df45f844d28881ad442d3

# https://golang.org/issue/14851
COPY no-pic.patch /

RUN set -ex \
        && apk add --no-cache --virtual .build-deps \
                bash \
                gcc \
                musl-dev \
                openssl \
                go \
        \
        && export GOROOT_BOOTSTRAP="$(go env GOROOT)" \
        \
        && wget -q "$GOLANG_SRC_URL" -O golang.tar.gz \
        && echo "$GOLANG_SRC_SHA256  golang.tar.gz" | sha256sum -c - \
        && tar -C /usr/local -xzf golang.tar.gz \
        && rm golang.tar.gz \
        && cd /usr/local/go/src \
        && patch -p2 -i /no-pic.patch \
        && ./make.bash \
        \
        && rm -rf /*.patch \
        && apk del .build-deps

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH" 

ENV GOPATH /go
ENV PATH $PATH:/go/bin

COPY ./ /go/src/github.com/appcelerator/docker-haproxy
RUN apk update && \
    apk --virtual build-deps add git make bash && \
    apk add curl && \
    cd $GOPATH/src/github.com/appcelerator/docker-haproxy && \
    go get -u github.com/Masterminds/glide/... && \
    glide install && \
    rm -rf vendor/github.com/appcelerator/amp/vendor && \
    make install && \
    cd / && apk del build-deps && rm -rf $GOPATH/src /var/cache/apk/* $GOPATH/pkg /root/.cache /root/.glide && \
    chmod +x $GOPATH/bin/*

COPY haproxy-main.cfg.tpt /usr/local/etc/haproxy/haproxy-main.cfg.tpt
COPY haproxy-stack.cfg.tpt /usr/local/etc/haproxy/haproxy-stack.cfg.tpt

HEALTHCHECK --interval=5s --timeout=10s --retries=12 CMD curl http://localhost/healthcheck

CMD ["/go/bin/docker-haproxy"]
