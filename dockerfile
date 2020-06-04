FROM golang:1.12.7
ENV GO111MODULE on
ENV GOPROXY https://goproxy.io
ADD . /code
WORKDIR /code
CMD ["go","run","main.go"]
