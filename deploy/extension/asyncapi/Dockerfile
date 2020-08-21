FROM golang:1.15-alpine as builder

ENV BASE_APP_DIR=/go/src/github.com/kyma-project/rafter \
    GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR ${BASE_APP_DIR}

# Copy the Go Modules manifests
COPY ./go.mod .
COPY ./go.sum .

# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

#
# Copy files
#
COPY . ${BASE_APP_DIR}/

#
# Build app
#
RUN go build -a -o main cmd/extension/asyncapi/main.go \
    && mkdir /app \
    && mv ./main /app/main \
    && if [ -f ${BASE_APP_DIR}/licenses ]; then mv ${BASE_APP_DIR}/licenses /app/licenses; fi

# Get latest CA certs
FROM alpine:latest as certs
RUN apk --update add ca-certificates

FROM scratch

LABEL source = git@github.com:kyma-project/rafter.git

COPY --from=builder /app /app
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

ENTRYPOINT ["/app/main"]