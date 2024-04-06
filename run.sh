#!/bin/bash
set -e
# set -x

execFile="modular-indexer-light"
command="./$execFile"
configExampleFile="config.example.json"
configFile="config.json"
bitcoinRPC="https://bitcoin-mainnet-archive.allthatnode.com"
randName=$(date +%s%N | base64 | head -c 6; echo)

if [[ ! -f "$configExampleFile" ]]; then
    echo "$configExampleFile not found"
    exit 1
fi
cp $configExampleFile $configFile

if [[ ! -f "$execFile" ]]; then
    echo "$execfile not found"
    exit 1
fi

read -p "Please enter a bitcoin rpc: " endpoint
if [[ ! "$endpoint" =~ ^(http|https)://([^/\s]+)/?([^/\s]*)\.?([^/\s]*)?$ ]]; then
    echo "Invalid bitcoin rpc, default bitcoin rpc $bitcoinRPC will be used."
else
    sed -i'' -e  "s|$bitcoinRPC|$endpoint|g" "$configFile"
fi

read -p "Please enter a Gas Coupon: " gasCoupon
gasCoupon=$(echo "$gasCoupon" | xargs)
if [[ 30 != ${#gasCoupon} ]]; then
    echo "Invalid Gas Coupon"
    exit 1
fi

if [[ $gasCoupon == "" ]]; then
    echo "Gas Coupon required!"
    exit 1
fi
sed -i'' -e "s/YourGasCoupon/$gasCoupon/g" "$configFile"

read -p "Please enter indexer name: " name
if [[ $name == "" ]]; then
    echo "Use randomly generated name"
    name=$randName
fi
sed -i'' -e "s/YourOwnLightIndexerName/$name/g" "$configFile"

echo "start modular-indexer-light...."
$command