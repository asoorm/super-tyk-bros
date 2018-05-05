#!/bin/bash

# https://github.com/veandco/go-sdl2-packaging-example

# build for mac os
go build

sudo rm -rf ./super-tyk-bros.app
mkdir -p super-tyk-bros.app/Contents/{MacOS,Frameworks}

# check for linked sdl libs
# change links to sdl libs to local frameworks directory
# copy libs to frameworks dir
# copy binary to MacOS dir
otool -L ./super-tyk-bros | \
    awk '{print $1;}' | \
    grep sdl2 | \
    while read lib ; \
    do install_name_tool -change $lib @executable_path/../Frameworks/${lib##*/} super-tyk-bros ;
    sudo cp $lib ./super-tyk-bros.app/Contents/Frameworks/ ;
    done
mv super-tyk-bros ./super-tyk-bros.app/Contents/MacOS/
cp -r assets ./super-tyk-bros.app/Contents/MacOS/

cat > ./super-tyk-bros.app/Contents/Info.plist <<- EOM
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>CFBundleExecutable</key>
	<string>super-tyk-bros</string>
	<key>CFBundleIdentifier</key>
	<string>co.asoorm.super-tyk-bros</string>
	<key>CFBundleVersion</key>
	<string>1.0</string>
	<key>CFBundleDisplayName</key>
	<string>Super Tyk Bros.</string>
	<key>LSRequiresIPhoneOS</key>
	<string>false</string>
</dict>
</plist>
EOM
