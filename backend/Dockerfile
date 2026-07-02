FROM golang:1.26-alpine AS build

WORKDIR /src
RUN apk add --no-cache git ca-certificates
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/airclean-server ./cmd/server

FROM alpine:3.22
WORKDIR /app
RUN apk add --no-cache ca-certificates tzdata
COPY --from=build /out/airclean-server /app/airclean-server
EXPOSE 8080
ENTRYPOINT ["/app/airclean-server"]

