FROM golang:1.24-alpine

WORKDIR /app
RUN adduser -D -u 10001 runner
USER runner

CMD ["go", "run", "solution.go"]
