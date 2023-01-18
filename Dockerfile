FROM golang:latest
RUN mkdir /app
ADD . /app/
WORKDIR /app/application/
RUN go get
RUN go build -o main .
CMD ["/app/application/main"]