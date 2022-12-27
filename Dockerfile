FROM golang:1.17-alpine3.14 as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ENV CGO_ENABLED 0
RUN go build ./cmd/give-me-bnb

FROM python:3.9-slim
RUN apt update && apt install -y tor netcat git ffmpeg libsm6 libxext6 wget && \
    apt clean
WORKDIR /app
COPY docker-entrypoint.sh .
RUN chmod +x docker-entrypoint.sh
COPY --from=builder /app/give-me-bnb .
ENV PATH="${PWD}:${PATH}"
ENTRYPOINT ["docker-entrypoint.sh"]
CMD ["give-me-bnb"]
