FROM node:22-alpine

WORKDIR /app
RUN npm install -g tsx && adduser -D -u 10001 runner
USER runner

CMD ["tsx", "solution.ts"]
