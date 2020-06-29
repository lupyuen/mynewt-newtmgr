# List Firmware
set -e -x

cd newtmgr
sudo ./newtmgr conn add pinetime type=ble connstring="peer_name=pinetime"

sudo ./newtmgr image list -c pinetime --loglevel debug
