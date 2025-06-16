# syntax=docker/dockerfile:1

##
## STEP 1 - BUILD
##
FROM golang:1.19-alpine AS build
RUN mkdir /app
COPY . /app/
WORKDIR /app/application/
RUN ls -la
RUN go mod download
RUN go build -o /godocker


##
## STEP 2 - DEPLOY
##
FROM alpine:3.18.2

WORKDIR /

COPY --from=build /godocker /godocker
COPY --from=build /app/application/.env env/.env
EXPOSE 8080

ENTRYPOINT ["/godocker"]


#FROM golang:latest
#RUN mkdir /app
#ADD . /app/
#WORKDIR /app/application/
#RUN go get
#RUN go build -o main .
#CMD ["/app/application/main"]