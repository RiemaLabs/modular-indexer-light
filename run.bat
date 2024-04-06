@echo off
setlocal enabledelayedexpansion

set "execFile=modular-indexer-light.exe"
set "command=.\modular-indexer-light.exe"
set "configExampleFile=config.example.json"
set "configFile=config.json"
set "bitcoinRPC=https://bitcoin-mainnet-archive.allthatnode.com"

if not exist "%configExampleFile%" (
    echo %configExampleFile% not found
    pause
    exit /b 1
)
copy /y %configExampleFile% %configFile%

if not exist "%execFile%" (
    echo %execFile% not found
    pause
    exit /b 1
)

set /p endpoint=Please enter a bitcoin rpc: 
echo %endpoint% | findstr /r /c:"^http[s]*:\/\/.*" > nul
if %errorlevel% neq 0 (
    echo Invalid bitcoin rpc, default bitcoin rpc %bitcoinRPC% will be used.
    set "endpoint=%bitcoinRPC%"
)
powershell -Command "(gc %configFile%) -replace '%bitcoinRPC%', '%endpoint%' | Out-File -encoding ASCII %configFile%"

set /p gasCoupon=Please enter a Gas Coupon: 
set gasCoupon=!gasCoupon!
if not "!gasCoupon!"=="" (
    if "!gasCoupon:~30,1!"=="" if not "!gasCoupon:~29,1!"=="" (
        powershell -Command "(gc %configFile%) -replace 'YourGasCoupon', '!gasCoupon!' | Out-File -encoding ASCII %configFile%"
    ) else (
        echo Invalid Gas Coupon
        pause
        exit /b 1
    )
) else (
    echo Gas Coupon required!
    pause
    exit /b 1
)

set /p name=Please enter indexer name: 
if "!name!"=="" (
    echo Use randomly generated name
    set /a name=%RANDOM%%RANDOM%
)
powershell -Command "(gc %configFile%) -replace 'YourOwnLightIndexerName', '!name!' | Out-File -encoding ASCII %configFile%"

echo start modular-indexer-light....
%command%