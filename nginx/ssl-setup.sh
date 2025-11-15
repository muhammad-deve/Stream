#!/bin/bash

echo "Starting SSL generation for $DOMAIN..."

# Generate SSL certificate
if [ ! -d "/etc/letsencrypt/live/$DOMAIN" ]; then
    echo "Generating SSL certificate..."
    certbot certonly --webroot --webroot-path /var/www/certbot -d $DOMAIN --email $EMAIL --agree-tos --no-eff-email --non-interactive
    
    if [ $? -eq 0 ]; then
        echo "SSL certificate generated successfully!"
        # Switch to HTTPS config
        cp /nginx-configs/default-https.conf /nginx-configs/default.conf
        echo "Switched to HTTPS configuration"
    else
        echo "Failed to generate SSL certificate, keeping HTTP configuration"
        # Keep HTTP config - don't switch to HTTPS
        echo "Continuing with HTTP-only setup"
    fi
else
    echo "SSL certificate already exists, using HTTPS config"
    cp /nginx-configs/default-https.conf /nginx-configs/default.conf
fi
