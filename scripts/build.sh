# Build Newt Manager
cd ~/go
mkdir -p src/mynewt.apache.org
cd src/mynewt.apache.org/
git clone https://github.com/apache/mynewt-newtmgr
mv mynewt-newtmgr newtmgr
cd newtmgr/newtmgr
export GO111MODULE=on
go build
