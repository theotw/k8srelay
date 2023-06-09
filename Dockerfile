# Base image
FROM theotw/devstack:go1.19.7  as base
LABEL stage=build
WORKDIR /build
RUN mkdir -p mods
COPY ./ ./
ENV PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/usr/local/go/bin/:/build/mods/bin
ENV GOSUMDB=off
ARG IMAGE_TAG=latest
ARG CA_KEY
ARG CA_CERT
ARG CA_SUB
ARG GOPATH=/build/mods

RUN make buildlinux
RUN make buildtest

# Alpine base
FROM alpine:3.14 as alpine-base
# Go doesn't use /etc/hosts so we need to create this /etc/nsswitch.conf file to redirect to /etc/hosts in case we ever
# use kubernetes' hostAliases (which uses /etc/hosts)
RUN echo 'hosts: files dns' > /etc/nsswitch.conf


# Test image
FROM alpine-base as k8srelay-test
WORKDIR /build
ARG IMAGE_TAG=latest
ENV GOSUMDB=off
COPY --from=base /build/out/*.test ./

# relay server
FROM alpine-base as k8srelayserver
RUN mkdir out
COPY --from=base /build/out/k8srelayserver_amd64_linux ./k8srelayserver
COPY --from=base /build/out/k8srelay.crt ./out/k8srelay.crt
COPY --from=base /build/out/k8srelay.key ./out/k8srelay.key
COPY --from=base /build/myCA.pem ./out/myCA.pem

ENV GIN_MODE=release

ENTRYPOINT ["./k8srelayserver"]

# relaylet
FROM alpine-base as k8srelaylet
ARG IMAGE_TAG=latest
ENV GOSUMDB=off
COPY --from=base /build/out/k8srelaylet_amd64_linux ./k8srelaylet


ENV GIN_MODE=release
ENTRYPOINT ["./k8srelaylet"]