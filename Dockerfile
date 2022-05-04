FROM golang:1.17-alpine

WORKDIR /app

COPY entrypoint.sh /entrypoint.sh
COPY . .

RUN go mod download && \
    go build -o /circleci-trigger-action && \
    chmod +x /entrypoint.sh

CMD [ "/entrypoint.sh" ]