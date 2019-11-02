FROM golang:1.12-alpine as builder
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
ARG VERSION=0.0.1

RUN apk update && apk add git
RUN git config --global core.compression 9

# build
WORKDIR /Users/tomcio/inf/go/src/github.com/knopt/k8s-sched-extender

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN go install -ldflags "-s -w -X main.version=$VERSION"

# runtime image
FROM gcr.io/google_containers/ubuntu-slim:0.14
COPY --from=builder /Users/tomcio/inf/go/src/github.com/knopt/k8s-sched-extender /usr/bin/k8s-sched-extender
ENTRYPOINT ["/usr/bin/k8s-sched-extender"]
