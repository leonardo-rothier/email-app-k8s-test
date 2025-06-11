FROM golang:1.24.3

WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code.
COPY *.go ./

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /email-site

EXPOSE 80

# Run
CMD ["/email-site"]
