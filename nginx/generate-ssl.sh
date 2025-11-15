#!/bin/bash

DOMAIN=${DOMAIN:-freetvchannels.online}
EMAIL=${EMAIL:-muhammadjonxudaynazarov226@gmail.com}

echo "Generating SSL certificate for $DOMAIN..."

# Start nginx first (for Let's Encrypt challenge)
docker compose up -d nginx

# Wait for nginx to start
sleep 5

# Generate SSL certificate
docker compose run --rm certbot certonly \
    --webconf \
    --webroot-path /var/www/certbot \
    -d $DOMAIN \
    --email $EMAIL \
    --agree-tos \
    --no-eff-email \
    --non-interactive

# Check if certificate was generated
if [ -d "certbot/conf/live/$DOMAIN" ]; then
    echo "✅ SSL certificate generated successfully!"
    echo "Now you can enable HTTPS by:"
    echo "1. Copy HTTPS config: cp conf.d/default-https.conf conf.d/default.conf"
    echo "2. Restart nginx: docker compose restart nginx"
else
    echo "❌ Failed to generate SSL certificate"
    echo "Check certbot logs: docker compose run --rm certbot certificates"
fi
