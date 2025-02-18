# Step 1: Use the official Go image to build the app
FROM golang:1.23-alpine as builder

# Step 2: Install SQLite3 and other dependencies
RUN apk update && apk add --no-cache \
    sqlite \
    sqlite-dev \
    bash \
    gcc \
    musl-dev \
    make

# Step 3: Set up a working directory for the app
WORKDIR /app

# Step 4: Copy go.mod and go.sum to ensure dependencies are downloaded first
COPY go.mod go.sum ./
RUN go mod tidy

# Step 5: Copy the rest of the application files to the container
COPY . .

# Step 6: Build the Go application
RUN go build -o forum-app ./cmd/main.go

# Step 7: Use a minimal Alpine image to run the app
FROM alpine:latest

# Step 8: Install SQLite runtime on the final image
RUN apk add --no-cache sqlite

# Step 9: Set the working directory
WORKDIR /app

# Step 10: Copy the compiled binary and static files from the builder image
COPY --from=builder /app/forum-app /app/forum-app
COPY static/ /app/static/
COPY templates/ /app/templates/

# Step 11: Expose the port the app runs on
EXPOSE 8080

# Step 12: Run the app when the container starts
CMD ["./forum-app"]
