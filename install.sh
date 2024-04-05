#!/bin/bash  
set -e

system=$(uname -s)
system=$(echo "$system" | tr '[:upper:]' '[:lower:]')
# latest_release=$(curl -s "https://api.github.com/repos/RiemaLabs/modular-indexer-light/releases/latest")  
version=$(wget -qO- -t1 -T2 "https://api.github.com/repos/RiemaLabs/modular-indexer-light/releases/latest" | grep "tag_name" | head -n 1 | awk -F ":" '{print $2}' | sed 's/\"//g;s/,//g;s/ //g')

zipfile="light-indexer-$system-amd64.zip"

download_url="https://github.com/RiemaLabs/modular-indexer-light/releases/download/$version/$zipfile"
wget -t2 -T2 -c $download_url
unzip $zipfile

rm -f $zipfile
/bin/bash run.sh
echo "最新 release 版本是: $download_url"