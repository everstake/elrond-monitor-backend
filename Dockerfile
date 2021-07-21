FROM golang:1.15
WORKDIR /app
COPY . /app
RUN go build -o app .
CMD ["./app"]
