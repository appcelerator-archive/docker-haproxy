FROM haproxy:1.6-alpine

ENV GOPATH /go
ENV PATH $PATH:/go/bin

COPY ./ /go/src/github.com/appcelerator/docker-haproxy
RUN echo "@edge http://nl.alpinelinux.org/alpine/edge/main" >> /etc/apk/repositories && \
    echo "@testing http://nl.alpinelinux.org/alpine/edge/testing" >> /etc/apk/repositories && \
    echo "@community http://nl.alpinelinux.org/alpine/edge/community" >> /etc/apk/repositories
RUN apk update && \
    apk -v --virtual build-deps add git make bash go@community musl-dev && \
    apk -v add curl go@community && \
    go version && \
    cd $GOPATH/src/github.com/appcelerator/docker-haproxy && \
    go get -u github.com/Masterminds/glide/... && \
    glide install && \
    rm -rf vendor/github.com/appcelerator/amp/vendor && \
    make install && \
    echo amp-haproxy built && \
    rm /go/bin/glide && \
    apk del binutils-libs binutils gmp isl libgomp libatomic libgcc pkgconf pkgconfig mpfr3 mpc1 libstdc++ gcc go && \
    cd / && rm -rf /go/src /go/pkg /var/cache/apk/* /root/.cache /root/.glide && \
    chmod +x $GOPATH/bin/*

COPY haproxy-main.cfg.tpt /usr/local/etc/haproxy/haproxy-main.cfg.tpt
COPY haproxy-stack.cfg.tpt /usr/local/etc/haproxy/haproxy-stack.cfg.tpt
COPY haproxy.cfg /usr/local/etc/haproxy/haproxy.cfg

HEALTHCHECK --interval=5s --timeout=10s --retries=12 CMD curl http://localhost/healthcheck

CMD ["/go/bin/docker-haproxy"]
