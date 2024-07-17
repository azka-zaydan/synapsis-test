ARG GO_VERSION=1.20
# Builder
FROM golang:${GO_VERSION}-alpine as builder

RUN apk update && \
    apk --update add git make build-base

WORKDIR /app

COPY . .

RUN go generate ./...
RUN go build -o goBinary .

# Distribution
FROM alpine:latest

RUN apk update && apk --no-cache add ca-certificates && \
    apk --update --no-cache add tzdata

ENV TZ=Asia/Jakarta

WORKDIR /app 

EXPOSE 3000

COPY --from=builder /app/goBinary /app
COPY --from=builder /app/.env /app
COPY --from=builder /app/docs /app/docs

CMD /app/goBinary