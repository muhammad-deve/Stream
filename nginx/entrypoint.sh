#!/bin/bash

DOMAIN=${DOMAIN:-freetvchannels.online}

echo "Starting nginx with automatic SSL configuration..."

# Check if SSL certificate exists
if [ -d "/etc/letsencrypt/live/$DOMAIN" ]; then
    echo "SSL certificate found! Starting with HTTPS..."
    exec nginx -g "daemon off;"
else
    echo "No SSL certificate found. This should not happen with the new setup..."
    echo "Starting with HTTP-only..."
    exec nginx -g "daemon off;"
fi
