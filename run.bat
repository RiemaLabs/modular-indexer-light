@echo off  
setlocal enabledelayedexpansion  
  
set "execFile=light-indexer.exe"  
set "configExampleFile=config.example.json"  
set "configFile=config.json"  
set "bitcoinRPC=https://bitcoin-mainnet-archive.allthatnode.com"  
set /a "randName=!RANDOM!%RANDOM%!RANDOM!%RANDOM%!RANDOM!%RANDOM!"  
randName=!randName:~0,6!  
  
if not exist "%configExampleFile%" (  
    echo %configExampleFile% not found  
    exit /b 1  
)  
  
copy "%configExampleFile%" "%configFile%" >nul  
  
if not exist "%execFile%" (  
    echo %execFile% not found  
    exit /b 1  
)  
  
set /p "endpoint=Please enter a Bitcoin RPC: "  
if not "!endpoint!"=="" (  
    echo !endpoint! | findstr /r /c:"^https?://[^/\s]+/?[^/\s]*\.?[^/\s]*$" >nul  
    if !errorlevel! eq 0 (  
        powershell -Command "(Get-Content '%configFile%').replace('%bitcoinRPC%', '!endpoint!') | Set-Content '%configFile%'"  
    ) else (  
        echo Invalid Bitcoin RPC, will use default Bitcoin RPC %bitcoinRPC%  
    )  
)  
  
set /p "report=Would you like to upload verified checkpoint to DA? [yes/no] "  
if /i "!report!"=="yes" (  
    set /p "gasCoupon=Please enter a Gas Coupon: "  
    if "!gasCoupon!"=="" (  
        echo Gas coupon is needed  
        exit /b 1  
    )  
    powershell -Command "(Get-Content '%configFile%').replace('YourGasCoupon', '!gasCoupon!') | Set-Content '%configFile%'"  
  
    set /p "name=Please enter indexer name: "  
    if "!name!"=="" (  
        set "name=!randName!"  
    )  
    powershell -Command "(Get-Content '%configFile%').replace('YourOwnLightIndexerName', '!name!') | Set-Content '%configFile%'"  
) else (  
    set "command=!execFile! --report=false"  
)  
  
echo Starting indexer...  
!command!  
endlocal