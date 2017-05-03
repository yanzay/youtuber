#!/bin/sh
set -e

echo "Building for ARM"
GOOS=linux GOARCH=arm GOARM=5 go build -v

echo "Stopping service"
ssh osmc@osmc "sudo sv stop youtuber"

echo "Copying"
scp youtuber osmc@osmc:~

echo "Starting service"
ssh osmc@osmc "sudo sv start youtuber"
