FROM golang as builder

WORKDIR /go/src/app

COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /go/bin/app

FROM scratch

COPY --from=builder /go/bin/app /app

CMD [ "/app" ]
