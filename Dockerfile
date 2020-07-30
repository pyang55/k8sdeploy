FROM golang:1.13.4-alpine
RUN apk add git
ENV GOOS linux
ENV GOARCH amd64
ENV CGO_ENABLED=0
WORKDIR /go/src/k8sdeploy
COPY . /go/src/k8sdeploy
RUN go get
RUN go build -o k8sdeploy
