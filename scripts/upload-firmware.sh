# Upload Firmware to PineTime
set -e -x

cd newtmgr
sudo ./newtmgr image upload -c pinetime --loglevel debug \
   ~/my_sensor_app_1.1.img \
    >../logs/upload-firmware.log 2>&1
cat ../logs/upload-firmware.log
