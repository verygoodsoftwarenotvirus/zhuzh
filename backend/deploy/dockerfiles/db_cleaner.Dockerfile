# build stage
FROM golang:1.26-trixie AS build-stage

WORKDIR /go/src/github.com/verygoodsoftwarenotvirus/zhuzh/backend

COPY . .

RUN go build -trimpath -o /action github.com/verygoodsoftwarenotvirus/zhuzh/backend/cmd/workers/db_cleaner

# final stage
FROM debian:bullseye

RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates
COPY --from=build-stage /action /action

ENTRYPOINT ["/action"]
