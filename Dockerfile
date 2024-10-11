FROM golang:1.23-bookworm AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags='-s' -o /bookshop ./cmd/api


FROM debian:bookworm-slim AS build-release-stage

WORKDIR /

COPY --from=build-stage /bookshop /bookshop
EXPOSE 4000

ENTRYPOINT ["/bookshop"]