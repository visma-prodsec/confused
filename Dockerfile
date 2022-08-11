FROM golang:latest

RUN go get -u github.com/visma-prodsec/confused
ENTRYPOINT ["confused"]
