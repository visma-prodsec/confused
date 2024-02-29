FROM golang:latest

RUN go install github.com/knavesec/confused@latest
ENTRYPOINT ["confused"]