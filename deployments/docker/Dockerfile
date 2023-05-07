FROM golang:1.19 as build

RUN mkdir /app
WORKDIR /app

ADD ../../ ./
RUN go build -o web_backend ./cmd/web/main.go

FROM busybox

RUN mkdir -p /go/bin
WORKDIR /go/bin

COPY --from=build /app/web_backend .

EXPOSE 8080

CMD ["./web_backend"]