FROM golang:1.21 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /simgameserver
RUN export SIM_GAME_DB_PW=valmet865
RUN mkdir logs/

FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /simgameserver /simgameserver

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/simgameserver"]
