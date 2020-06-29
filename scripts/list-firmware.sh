# List Firmware on PineTime
set -e -x

cd newtmgr
sudo ./newtmgr conn add pinetime type=ble connstring="peer_name=pinetime"

sudo ./newtmgr image list -c pinetime --loglevel debug \
    >../logs/list-firmware.log 2>&1
cat ../logs/list-firmware.log
