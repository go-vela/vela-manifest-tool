# SPDX-License-Identifier: Apache-2.0

################################################################################
##    docker build --no-cache --target certs -t vela-manifest-tool:certs .    ##
################################################################################

FROM alpine:3.21.0@sha256:21dc6063fd678b478f57c0e13f47560d0ea4eeba26dfc947b2a4f81f686b9f45 as certs

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
