# SPDX-License-Identifier: Apache-2.0

################################################################################
##    docker build --no-cache --target certs -t vela-manifest-tool:certs .    ##
################################################################################

FROM alpine:3.20.1@sha256:b89d9c93e9ed3597455c90a0b88a8bbb5cb7188438f70953fede212a0c4394e0 as certs

RUN apk add --update --no-cache ca-certificates

#################################################################
##    docker build --no-cache -t vela-manifest-tool:local .    ##
#################################################################

FROM mplatform/manifest-tool:alpine-v2.1.6@sha256:96db9e944c50a5f7514394af4e44f764725645cfd2efef92d87095b0016a55ae

COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

WORKDIR /workspace

RUN mkdir /root/.docker

COPY release/vela-manifest-tool /bin/vela-manifest-tool

ENTRYPOINT [ "/bin/vela-manifest-tool" ]
