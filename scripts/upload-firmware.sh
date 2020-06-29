# Upload Firmware
set -e -x

cd newtmgr
sudo ./newtmgr image upload -c pinetime \
   ~/my_sensor_app_1.1.img
