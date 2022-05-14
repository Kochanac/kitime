FROM --platform=linux/amd64 golang:alpine as builder
RUN mkdir /build
ADD . /build/
WORKDIR /build
RUN go build -o head cmd/main.go

FROM --platform=linux/amd64 alpine 
COPY --from=builder /build/head /app/
WORKDIR /app
CMD ["./head"]
