#!/bin/bash -e

echo "Update /etc/hosts"
sudo bash -c 'echo "127.1.2.3 spinnaker" >> /etc/hosts'

