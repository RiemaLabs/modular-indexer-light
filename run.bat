@echo off
setlocal enabledelayedexpansion

set execFile=light-indexer.exe
set command=light-indexer.exe
set configFile=config.json
set configExampleFile=config.example.json
set bitcoinRPC=https://bitcoin-mainnet-archive.allthatnode.com

if not exist "%configExampleFile%" (
    echo %configExampleFile% not found
    exit /b 1
)
copy /y %configExampleFile% %configFile%

if not exist "%execFile%" (
    echo %execFile% not found
    exit /b 1
)

set /p endpoint=please enter a bitcoin rpc: 
echo %endpoint% | findstr /r /c:"^http[s]*:\/\/.*" > nul
if %errorlevel% neq 0 (
    echo invalid bitcoin rpc, will use default bitcoin rpc %bitcoinRPC%
    set endpoint=%bitcoinRPC%
)
powershell -Command "(gc %configFile%) -replace '%bitcoinRPC%', '%endpoint%' | Out-File -encoding ASCII %configFile%"

set /p report=would you like upload verified checkpoint to DA ? [yes/no] 
if /i "%report%"=="yes" (
    set /p gasCoupon=please enter a Gas Coupon: 
    if "!gasCoupon!"=="" (
        echo gas coupon is needed
        exit /b 1
    )
    powershell -Command "(gc %configFile%) -replace 'YourGasCoupon', '!gasCoupon!' | Out-File -encoding ASCII %configFile%"

    set /p name=please enter indexer name: 
    if "!name!"=="" (
        echo use randomly generated name
        set /a name=%RANDOM%
    )
    powershell -Command "(gc %configFile%) -replace 'YourOwnLightIndexerName', '!name!' | Out-File -encoding ASCII %configFile%"
) else (
    set command=%execFile% --report=false
)

echo start indexer....
%command%