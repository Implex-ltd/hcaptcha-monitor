@echo off
cd ../
go get
cd cmd/hcmonitor
go build -o ../../scripts/ .
cd ../../scripts/
cls
hcmonitor.exe
pause