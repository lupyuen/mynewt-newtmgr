<!--
#
# Licensed to the Apache Software Foundation (ASF) under one
# or more contributor license agreements.  See the NOTICE file
# distributed with this work for additional information
# regarding copyright ownership.  The ASF licenses this file
# to you under the Apache License, Version 2.0 (the
# "License"); you may not use this file except in compliance
# with the License.  You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing,
# software distributed under the License is distributed on an
# "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
#  KIND, either express or implied.  See the License for the
# specific language governing permissions and limitations
# under the License.
#
-->

# Updates for PineTime Companion App

See [`newtmgr.dart`](newtmgr.dart) for newtmgr converted from Go to Dart.

# Newtmgr

Newt Manager (newtmgr) is the application tool that enables a user to communicate with and manage
remote devices running the Mynewt OS. It uses a connection profile to establish a connection with
a device and sends command requests to the device.
The newtmgr tool documentation can be found under [/docs](/docs) which are
published at http://mynewt.apache.org/latest/os/modules/devmgmt/newtmgr.html

### Building and Running

Build and run the newtmgr tool as follows:

```bash
sudo apt install graphviz

# Download Newt Manager on Raspberry Pi or Pinebook Pro
cd ~/go
mkdir -p src/mynewt.apache.org
cd src/mynewt.apache.org/
git clone https://github.com/lupyuen/mynewt-newtmgr
mv mynewt-newtmgr newtmgr

# Build Newt Manager on Raspberry Pi or Pinebook Pro
cd ~/go/src/mynewt.apache.org/newtmgr/newtmgr
export GO111MODULE=on
go build

# Run Newt Manager on Raspberry Pi or Pinebook Pro
cd ~/go/src/mynewt.apache.org/newtmgr/newtmgr

sudo ./newtmgr conn add pinetime type=ble connstring="peer_name=pinetime" 2> /dev/null

sudo ./newtmgr image list -c pinetime 2> trace.out

go tool trace trace.out
```