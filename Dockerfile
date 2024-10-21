FROM golang:1.23.2-bookworm
WORKDIR /app
COPY ./go.mod .
COPY ./weather.go .
COPY ./go.sum .
RUN go build -o /app ./weather.go
CMD ["/app/weather"]