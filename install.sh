#!/bin/bash
set -e
set -x

# command check
commandExist() {
    if ! command -v $1 >/dev/null 2>&1; then
        echo "The $1 command is necessary. Please install the $1 software first."
        exit 1
    fi
}

# handle WSL system
if uname -r | grep -qEi "(Microsoft|WSL)"; then
    echo "WSL detected. This script is not fully compatible with WSL. Please download the Windows runner instead by clicking this link"
    exit 1
fi

commandExist "unzip"
commandExist "curl"

system=$(uname -s)
system=$(echo "$system" | tr '[:upper:]' '[:lower:]')

arch=$(uname -m)
if [ $arch = "x86_64" ]; then
    arch="amd64"
fi

# latest_release=$(curl -s "https://api.github.com/repos/RiemaLabs/modular-indexer-light/releases/latest")
version=$(curl -s "https://api.github.com/repos/RiemaLabs/modular-indexer-light/releases/latest" | grep "tag_name" | head -n 1 | awk -F ":" '{print $2}' | sed 's/\"//g;s/,//g;s/ //g')

zipfile="light-indexer-$system-$arch.zip"

download_url="https://github.com/RiemaLabs/modular-indexer-light/releases/download/$version/$zipfile"
curl -SfLO $download_url
if [ -f "$zipfile" ]; then
    unzip "$zipfile"
else
    echo "Failed to download the $zipfile."
    exit 1
fi

rm -f $zipfile
/bin/bash run.sh

# sh -c "$(curl -fsSL https://raw.githubusercontent.com/RiemaLabs/modular-indexer-light/main/install.sh)"
# sh -c "$(wget https://raw.githubusercontent.com/RiemaLabs/modular-indexer-light/main/install.sh -O -)"
