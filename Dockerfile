FROM golang:1.14-alpine

WORKDIR /go/src/ganjifa
COPY . .

RUN apk add --update nodejs npm
RUN cd ./webapp && npm install && npm run build
RUN cd ..

RUN go get -d -v ./...
RUN go install -v ./...

EXPOSE 80

CMD ["ganjifa"]