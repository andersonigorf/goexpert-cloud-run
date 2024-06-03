FROM golang:1.22.3
WORKDIR /app
COPY . .
CMD ["go", "run", "main.go"]
