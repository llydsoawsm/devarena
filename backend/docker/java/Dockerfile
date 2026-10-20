FROM eclipse-temurin:21-jre-alpine

RUN apk add --no-cache openjdk21 && \
    adduser -D -u 10001 runner

WORKDIR /app
USER runner

CMD ["sh", "-c", "javac Solution.java && java Solution"]
