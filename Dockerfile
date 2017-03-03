FROM alpine
RUN apk add --no-cache ca-certificates
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
ADD brickbot brickbot
ADD settings.json settings.json
ADD templates templates
CMD ["/brickbot", "-c", "settings.json"]
