FROM golang:1.22 AS builder

WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
WORKDIR /app/cmd
ENV GOCACHE=/root/.cache/go-build
RUN --mount=type=cache,target="/root/.cache/go-build" CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} go build -o ht3mock

FROM alpine
RUN mkdir /app
WORKDIR /app
COPY --from=builder /app/cmd/ht3mock ht3mock
CMD /app/ht3mock -svr