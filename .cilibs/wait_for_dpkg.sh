#!/bin/bash -e

sleep 10
while systemctl status apt-daily >/dev/null || systemctl status apt-daily-upgrade >/dev/null || sudo fuser /var/{lib/{dpkg,apt/lists},cache/apt/archives}/lock; do
echo "waiting 30s for dpkg locks..."
sleep 30
done
