FROM golang:1.5.1-wheezy
RUN apt-get update
RUN apt-get install -y supervisor
ENV GOPATH=/go
ENV GOARCH=amd64
RUN go get github.com/fsouza/go-dockerclient
RUN go get github.com/nlopes/slack
