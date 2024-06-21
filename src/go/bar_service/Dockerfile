# Build the manager binary
FROM golang:1.21 AS builder
ARG TARGETOS
ARG TARGETARCH

WORKDIR /workspace
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download
COPY cmd/main.go cmd/main.go
COPY svcctl/ svcctl/
COPY barsvc/ barsvc/

RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} go build -a cmd/main.go

FROM alpine:3.19
WORKDIR /
COPY --from=builder /workspace/main barsvc
USER 65532:65532

ENTRYPOINT ["/barsvc"]
