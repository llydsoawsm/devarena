FROM node:22-alpine

WORKDIR /app
RUN adduser -D -u 10001 runner
USER runner

CMD ["node", "solution.js"]
