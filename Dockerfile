FROM golang:1.5.1-wheezy
RUN apt-get update
RUN apt-get install -y supervisor
ENV GOPATH=/go
ENV GOARCH=amd64
RUN go get github.com/fsouza/go-dockerclient
RUN go get github.com/nlopes/slack
ADD . /server
RUN cd /server && go build -o cortana
CMD ["/usr/bin/supervisord", "-n", "-c", "/server/supervisord.conf"]
