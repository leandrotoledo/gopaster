#!/bin/sh

# Check if server.crt and server.key exist, otherwise generate self-signed certificates
if [ ! -f /app/certs/server.crt ] || [ ! -f /app/certs/server.key ]; then
    echo "Generating self-signed certificates..."
    mkdir -p /app/certs
    openssl req -x509 -newkey rsa:4096 -keyout /app/certs/server.key -out /app/certs/server.crt -days 365 -nodes -subj "/CN=localhost"
fi

# Start the Go application
./gopaster
