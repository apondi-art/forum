#!/bin/bash

# Stop forum-app container if running
docker stop forum-app-container >/dev/null 2>&1 || true

# Remove the container
docker rm forum-app-container >/dev/null 2>&1 || true

# Delete the forum-app image if exists
docker rmi forum-app >/dev/null 2>&1 || true

# Build docker image
echo "Building forum-app Docker image..."
docker build -t forum-app .

# Fix for the CMD error in the Dockerfile
if [ $? -ne 0 ]; then
  echo "Build failed. Attempting to fix Dockerfile..."
  # Try to fix the command line in the Dockerfile
  sed -i 's/CMD\["\.\/forum-app"\]e/CMD ["\.\/forum-app"]/' Dockerfile
  echo "Retrying build with fixed Dockerfile..."
  docker build -t forum-app .
fi

# Run the image on forum-app-container on port 8080
echo "Starting forum-app container..."
docker run -p 8080:8080 --name forum-app-container forum-app

# Note: You might want to add -d flag to run in detached mode
# docker run -d -p 8080:8080 --name forum-app-container forum-app