FROM mcr.microsoft.com/dotnet/sdk:8.0-alpine

RUN adduser -D -u 10001 runner
WORKDIR /app
USER runner

CMD ["sh", "-c", "dotnet run --project ."]
