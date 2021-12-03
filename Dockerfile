FROM golang:1.17-alpine

WORKDIR /go/src/app
COPY . .

# server port
EXPOSE 5000

RUN chmod +x command
RUN /bin/sh command --build

CMD ["./bin/chapi"]