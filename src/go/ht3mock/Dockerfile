FROM golang:1.22 AS builder

WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
WORKDIR /app/cmd
ENV GOCACHE=/root/.cache/go-build
RUN --mount=type=cache,target="/root/.cache/go-build" go build -o ht3mok

FROM alpine
RUN mkdir /app
WORKDIR /app
COPY --from=builder /app/cmd/ht3mok .
ENTRYPOINT ht3mok -svr