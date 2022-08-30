FROM amd64/golang:1.17-alpine AS build-env

ENV GO111MODULE=auto
ENV GOOS=linux
ENV GOARCH=amd64

RUN apk update && apk add bash ca-certificates git gcc libc-dev

WORKDIR /src
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY internal internal
COPY dist dist
COPY main.go main.go

RUN go build -o launcher main.go

FROM amd64/ubuntu:bionic AS run-env

RUN apt update && apt install bash ca-certificates git gcc libc-dev -y && apt install musl -y

WORKDIR /srv

COPY --from=build-env /src/launcher .
COPY --from=build-env /src/dist ./dist

CMD ./launcher --port-grpc=25565 --port-rest=8080 --port-swagger_ui=80
