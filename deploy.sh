#!/bin/sh
set -e

echo "Building for ARM"
GOOS=linux GOARCH=arm GOARM=5 go build -v

echo "Stopping service"
ssh osmc@yanzay.com -p 2222 "sudo sv stop youtuber"

echo "Copying"
scp -P 2222 youtuber osmc@yanzay.com:~

echo "Starting service"
ssh osmc@yanzay.com -p 2222 "sudo sv start youtuber"
