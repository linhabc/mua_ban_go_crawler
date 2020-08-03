FROM golang:latest

LABEL maintainer="linhnln"

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go build crawler_phone_using_worker.go data_type.go util.go db.go

CMD [ "./crawler_phone_using_worker" ]