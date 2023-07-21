#docker build -t go-hello-world-web .
#docker run -dit --name go-hello-world-web -p 3000:3000 go-hello-world-web

FROM golang:1.18 As builder

WORKDIR /app
COPY . .

WORKDIR /app/cmd
RUN go build -o go-hello-world-web -ldflags '-linkmode external -w -extldflags "-static"'

FROM alpine

WORKDIR /app
COPY --from=builder /app/cmd/go-hello-world-web .
COPY --from=builder /app/init.sh .

CMD ["/app/go-hello-world-web"]
#CMD ["/app/init.sh"]