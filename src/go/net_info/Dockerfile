# Build the manager binary
FROM golang:1.22 AS builder
ARG TARGETOS
ARG TARGETARCH

WORKDIR /workspace
COPY go.mod go.mod
RUN go mod download
COPY cmd/main.go cmd/main.go

RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} go build -a cmd/main.go

FROM alpine:3.19
WORKDIR /
COPY --from=builder /workspace/main net_info
USER 65532:65532

ENTRYPOINT ["/net_info"]
