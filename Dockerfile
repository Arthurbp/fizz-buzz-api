FROM golang:1.14

WORKDIR /go/src/fizz-buzz-api
COPY . .

RUN go get ./
RUN go build -o fizz-buzz-api

EXPOSE 8080

CMD [ "fizz-buzz-api" ]