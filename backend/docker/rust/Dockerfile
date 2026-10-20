FROM rust:alpine

RUN adduser -D -u 10001 runner
WORKDIR /app
USER runner

CMD ["sh", "-c", "rustc -O -o solution solution.rs && ./solution"]
