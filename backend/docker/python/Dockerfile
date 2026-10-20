FROM python:3.12-alpine

WORKDIR /app
RUN adduser -D -u 10001 runner
USER runner

CMD ["python", "solution.py"]
