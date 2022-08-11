FROM golang:latest

RUN go install github.com/visma-prodsec/confused@latest
ENTRYPOINT ["confused"]
