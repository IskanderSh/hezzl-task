FROM golang:alpine AS builder

WORKDIR /usr/local/src

ADD go.mod .

COPY . .

RUN go build -o ./bin/application cmd/main.go

FROM alpine AS runner

COPY --from=builder /usr/local/src/bin/application .
COPY config/prod.yaml .

ENV CONFIG_PATH=./prod.yaml

CMD ["/application"]