# SPDX-License-Identifier: Apache-2.0

################################################################################
##    docker build --no-cache --target certs -t vela-manifest-tool:certs .    ##
################################################################################

FROM alpine:3.24.2@sha256:294b683cb724975bec92580e1e685676bd4b50bda910ddb8c51d4cabeaec77e6 as certs

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
