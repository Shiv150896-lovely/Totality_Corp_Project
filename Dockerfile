FROM golang:1.16-alpine

WORKDIR /app

COPY . .

RUN go mod tidy
RUN go build -o main .

EXPOSE 50051

CMD ["./main"]

