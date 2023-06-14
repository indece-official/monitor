# indece Monitor
# Copyright (C) 2023 indece UG (haftungsbeschränkt)
#
# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation, either version 3 of the License or any
# later version.
#
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
# GNU General Public License for more details.
#
# You should have received a copy of the GNU General Public License
# along with this program. If not, see <https://www.gnu.org/licenses/>.

FROM node:18.16.0-alpine3.17 as buildfrontend

ARG BUILD_VERSION=""

USER root

RUN mkdir -p /app

WORKDIR /app

COPY . .

WORKDIR /app/frontend

RUN npm ci
RUN npm run lint
#RUN npm run test
RUN npm run build

FROM golang:1.20.4-alpine3.17 as buildbackend

USER root

WORKDIR /go/src/app

COPY . .

COPY --from=buildfrontend /app/frontend/build ./frontend/build

ARG BUILD_VERSION=""
ENV PROJECT_NAME "indece-monitor"
ENV GO111MODULE "on"

WORKDIR /go/src/app/backend

RUN apk add --no-cache make git build-base pkgconfig protobuf

RUN go version && \
    go env && \
    go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v1.12.4 && \
    go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.30.0 && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0 && \
    go install github.com/securego/gosec/v2/cmd/gosec@v2.16.0

WORKDIR /go/src/app

RUN make --always-make copy_frontend && \
    make --always-make build_backend

USER nobody

FROM alpine:3.17.0

USER nobody

WORKDIR /app

COPY --from=buildbackend /etc/passwd /etc/passwd
COPY --from=buildbackend /go/src/app/backend/dist/bin .

ENTRYPOINT ["/app/indece-monitor"]
