FROM public.ecr.aws/amazonlinux/amazonlinux:latest

RUN yum install -y go

ENV CGO_ENABLED=1

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY *.go .
RUN go build -ldflags="-s -w"