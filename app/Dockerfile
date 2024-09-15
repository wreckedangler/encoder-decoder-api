FROM golang:1.22-alpine

WORKDIR /app

COPY go.mod go.sum ./
COPY app/ ./app/
COPY main.go ./
COPY .env ./

# Abhängigkeiten herunterladen
RUN go mod download

# Exponiere den Port 8080
EXPOSE 8080

# Baue die API
RUN go build -o main .

# Starte die API
CMD ["./main"]
