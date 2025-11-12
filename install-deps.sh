#!/bin/bash

# Install dependencies for CI/CD pipeline
echo "Installing frontend dependencies..."
cd front
npm install

echo "Installing backend dependencies..."
cd ../back/app
go mod download

echo "Dependencies installed successfully!"
echo "Next steps:"
echo "1. Configure GitHub secrets in your repository"
echo "2. Set up your production server"
echo "3. Push to main branch to trigger deployment"
