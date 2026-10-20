FROM alpine:3.20

RUN apk add --no-cache g++ build-base bash && \
    adduser -D -u 10001 runner

WORKDIR /app
USER runner

CMD ["sh", "-c", "g++ -O3 -o solution solution.cpp && ./solution"]
