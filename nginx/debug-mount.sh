#!/bin/bash

echo "=== Debugging volume mount ==="
echo "Files in root:"
docker run --rm -v $(pwd)/ssl-setup.sh:/ssl-setup.sh certbot/certbot ls -la /

echo "=== Files in ssl-setup.sh location ==="
docker run --rm -v $(pwd)/ssl-setup.sh:/ssl-setup.sh certbot/certbot ls -la /ssl-setup.sh

echo "=== Trying to run script ==="
docker run --rm -v $(pwd)/ssl-setup.sh:/ssl-setup.sh certbot/certbot /bin/sh -c "ls -la /ssl-setup.sh && chmod +x /ssl-setup.sh && /ssl-setup.sh"
