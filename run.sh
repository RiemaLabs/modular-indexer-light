#!/bin/bash
set -e
# set -x

execFile="light-indexer"
command="./$execFile"
configExampleFile="config.example.json"
configFile="config.json"
bitcoinRPC="https://bitcoin-mainnet-archive.allthatnode.com"
randName=$(cat /dev/urandom | tr -dc 'a-zA-Z0-9' | fold -w 6 | head -n 1)

if [[ ! -f "$configExampleFile" ]]; then 
    echo "$configExampleFile not found"
    exit 1
fi
cp $configExampleFile $configFile

if [[ ! -f "$execFile" ]]; then
    echo "$execfile not found"
    exit 1
fi

read -p "please enter a bitcoin rpc: " endpoint
if [[ ! "$endpoint" =~ ^(http|https)://([^/\s]+)/?([^/\s]*)\.?([^/\s]*)?$ ]]; then  
    echo "invalid bitcoin rpc, will use default bitcoin rpc $bitcoinRPC"
else  
    sed -i "s|$bitcoinRPC|$endpoint|g" "$configFile"
fi
 

read -p "would you like upload verified checkpoint to DA ? [yes/no] " report
report=$(echo "$report" | xargs) 

if [[ "$report" == "yes" ]]; then  
    read -p "please enter a Gas Coupon: " gasCoupon
    gasCoupon=$(echo "$gasCoupon" | xargs)
    if [[ $gasCoupon == "" ]]; then
        echo "gas coupon is needed"
        exit 1
    fi
    sed -i "s/YourGasCoupon/$gasCoupon/g" "$configFile"

    read -p "please enter indexer name: " name
    if [[ $name == "" ]]; then
        echo "use randomly generated name"
        $name=$randName
    fi
    sed -i "s/YourOwnLightIndexerName/$name/g" "$configFile"

else 
    command="$command --report=false"
fi

echo "start indexer...."
$command
