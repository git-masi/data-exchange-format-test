#!/bin/bash

# Update package list and install git and wget
sudo apt update
sudo apt install -y git wget

# Download and install Go 1.20.5
wget https://golang.org/dl/go1.20.5.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.20.5.linux-amd64.tar.gz
rm go1.20.5.linux-amd64.tar.gz

# Add Go binary to PATH
echo "export PATH=$PATH:/usr/local/go/bin" >>~/.bashrc
source ~/.bashrc

# Clone your GitHub repository (replace with your actual GitHub repo URL)
git clone https://github.com/git-masi/data-exchange-format-test.git

# Navigate to the server directory and build the server
cd data-exchange-format-test/cmd/server
go build

# Run the server
./server
