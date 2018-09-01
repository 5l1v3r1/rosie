#!/bin/bash

# Creates the static go asset archives
# You'll need wget, tar, and unzip commands

GO_VER="1.11"
BLOAT_FILES="AUTHORS CONTRIBUTORS PATENTS VERSION favicon.ico robots.txt CONTRIBUTING.md LICENSE README.md ./doc ./test"


REPO_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
WORK_DIR=`mktemp -d`

echo "-----------------------------------------------------------------"
echo $WORK_DIR
echo "-----------------------------------------------------------------"
cd $WORK_DIR

# --- Darwin --- 
wget https://dl.google.com/go/go$GO_VER.darwin-amd64.tar.gz
tar xvf go$GO_VER.darwin-amd64.tar.gz

cd go
rm -rf $BLOAT_FILES
zip -r ../src.zip ./src  # Zip up /src we only need to do this once
rm -rf ./src
rm -f ./pkg/tool/darwin_amd64/doc
rm -f ./pkg/tool/darwin_amd64/tour
rm -f ./pkg/tool/darwin_amd64/test2json
cd ..
cp -vv src.zip $REPO_DIR/server/assets/src.zip
rm -f src.zip

zip -r darwin-go.zip ./go
cp -vv darwin-go.zip $REPO_DIR/server/assets/darwin/go.zip

rm -rf ./go
rm -f darwin-go.zip go$GO_VER.darwin-amd64.tar.gz


# --- Linux --- 
wget https://dl.google.com/go/go$GO_VER.linux-amd64.tar.gz
tar xvf go$GO_VER.linux-amd64.tar.gz
cd go
rm -rf $BLOAT_FILES
rm -rf ./src
rm -f ./pkg/tool/linux_amd64/doc
rm -f ./pkg/tool/linux_amd64/tour
rm -f ./pkg/tool/linux_amd64/test2json
cd ..
zip -r linux-go.zip ./go
cp -vv linux-go.zip $REPO_DIR/server/assets/linux/go.zip
rm -rf ./go
rm -f linux-go.zip go$GO_VER.linux-amd64.tar.gz

# --- Windows --- 
wget https://dl.google.com/go/go$GO_VER.windows-amd64.zip
unzip go$GO_VER.windows-amd64.zip
cd go
rm -rf $BLOAT_FILES
rm -rf ./src
rm -f ./pkg/tool/windows_amd64/doc.exe
rm -f ./pkg/tool/windows_amd64/tour.exe
rm -f ./pkg/tool/windows_amd64/test2json.exe
cd ..
zip -r windows-go.zip ./go
cp -vv windows-go.zip $REPO_DIR/server/assets/windows/go.zip
rm -rf ./go
rm -f windows-go.zip go$GO_VER.windows-amd64.zip

# end
echo -e "clean up: $WORK_DIR"
rm -rf $WORK_DIR
echo -e "\n[*] All done\n"
