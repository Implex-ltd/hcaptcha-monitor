clear
cd ..
go get -u
cd cmd/hcmonitor
go build -o ../../scripts/ .
cd ../../scripts/
echo "================================================="
./hcmonitor
