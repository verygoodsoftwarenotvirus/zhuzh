# build stage
FROM golang:1.26-trixie AS build-stage

WORKDIR /go/src/github.com/verygoodsoftwarenotvirus/zhuzh/backend

COPY . .

RUN ./scripts/build.sh -o /server github.com/verygoodsoftwarenotvirus/zhuzh/backend/cmd/services/api

# final stage
FROM debian:bullseye

RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates
COPY --from=build-stage /server /server

ENTRYPOINT ["/server", "serve"]
