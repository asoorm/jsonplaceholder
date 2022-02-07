FROM golang:1.17-alpine as build

ENV GO111MODULE=on
WORKDIR $GOPATH/src
#ARG VERSION
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

COPY . .

RUN go install -a -ldflags="-s -w" .

FROM alpine:3

ENV HOST=0.0.0.0
ENV PORT=8000

RUN apk --no-cache add ca-certificates
RUN adduser -D -g jsonplaceholder jsonplaceholder
USER jsonplaceholder

WORKDIR /opt/jsonplaceholder
COPY --from=build /go/bin/jsonplaceholder /opt/jsonplaceholder
COPY ./db ./db
USER jsonplaceholder

#EXPOSE $PORT

CMD ["./jsonplaceholder"]
