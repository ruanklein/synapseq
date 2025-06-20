#!/bin/bash

# Build directory
BUILD_DIR="$PWD/build"

# Installer directory
INSTALLER_DIR="$BUILD_DIR/macos-installer"

# Documentation directory
DOC_DIR="$PWD/docs"

# Source common library
. $BUILD_DIR/lib.sh

APP_NAME="SynapSeq"
SYNAPSEQ_BINARY="$BUILD_DIR/dist/synapseq-macos-universal"
PNG_SOURCE="$BUILD_DIR/assets/mac-icon.png"
ICON_NAME="app_icon"
DMG_NAME="SynapSeq-Installer.dmg"

section_header "Creating macOS Application Bundle..."

# Clean build directory
rm -rf $INSTALLER_DIR

# Check if the binary exists
if [ ! -f "$SYNAPSEQ_BINARY" ]; then
    error "$SYNAPSEQ_BINARY not found. Execute ./macos-build-synapseq.sh first."
    exit 1
fi

# Check if the PNG file exists
if [ ! -f "$PNG_SOURCE" ]; then
    error "$PNG_SOURCE not found!"
    exit 1
fi

# Create build directory if it doesn't exist
create_dir_if_not_exists "$INSTALLER_DIR"

# Create temporary AppleScript with the dialog
info "Creating AppleScript handler..."
cat > $INSTALLER_DIR/synapseq.applescript <<EOF
on showAppNotInstalledAlert()
    set appPath to POSIX path of (path to me)

    if appPath starts with "/Volumes/" then
        display dialog "Please move the app to the Applications folder before using it." buttons {"OK"} default button 1
        error number -128
    end if
end showAppNotInstalledAlert

on initialize()
	set documentsFolder to POSIX path of (path to documents folder)
	set targetFolder to documentsFolder & "SynapSeq"
	set appBase to "/Applications/SynapSeq.app/Contents/Resources"
    set desktopFolder to POSIX path of (path to desktop folder)
    set linkPath to desktopFolder & "SynapSeq Files"

    set folderExists to (do shell script "test -d " & quoted form of targetFolder & " && echo yes || echo no")
	if folderExists is "yes" then return

    set noticePath to ((POSIX path of (path to me)) & "Contents/Resources/NOTICE.txt")
    set noticeText to do shell script "cat " & quoted form of noticePath
    set noticeText to noticeText & return & "––" & return & return & "IMPORTANT: By clicking 'OK, I AGREE', you confirm that you have read and accepted the license terms."

    set userChoice to button returned of (display dialog noticeText buttons {"View License", "OK, I Agree", "Cancel"} default button "View License" with title "GPL License")

	if userChoice is "View License" then
		set licensePath to ((POSIX path of (path to me)) & "Contents/Resources/COPYING.txt")
		tell application "TextEdit"
			activate
			open POSIX file licensePath
		end tell
        error number -128
	end if

    if userChoice is "Cancel" then
        error number -128
    end if

	do shell script "mkdir -p " & quoted form of targetFolder

	set fileMap to {¬
	    {"Documentation", appBase & "/docs"}, ¬
	    {"Samples", appBase & "/samples"}, ¬
	    {"License.txt", appBase & "/COPYING.txt"}, ¬
        {"Notice.txt", appBase & "/NOTICE.txt"}, ¬
	    {"ChangeLog.txt", appBase & "/ChangeLog.txt"}}

	repeat with pair in fileMap
		set fileName to item 1 of pair
		set targetPath to item 2 of pair
		set filePath to targetFolder & "/" & fileName

		set fileExists to (do shell script "test -e " & quoted form of filePath & " && echo yes || echo no")

		if fileExists is "no" then
			do shell script "cp -R " & quoted form of targetPath & " " & quoted form of filePath
		end if
	end repeat

    set linkExists to (do shell script "test -e " & quoted form of linkPath & " && echo yes || echo no")
    if linkExists is "no" then
	    do shell script "ln -s " & quoted form of targetFolder & " " & quoted form of linkPath
    end if
end initialize

on run
    showAppNotInstalledAlert()
    initialize()

    set documentsFolder to POSIX path of (path to documents folder)
    set examplesPath to documentsFolder & "SynapSeq/Examples"
    set iconPath to ((POSIX path of (path to me)) & "Contents/Resources/app_icon.icns")
    set dialogText to "Please open a .sbg file using this application."

    set examplesExists to (do shell script "test -d " & quoted form of examplesPath & " && echo yes || echo no")
    
    if examplesExists is "yes" then
        set userChoice to button returned of (display dialog dialogText buttons {"OK", "View Examples"} default button "OK" with title "$APP_NAME" with icon POSIX file iconPath)
        if userChoice is "View Examples" then
            tell application "Finder"
                open (POSIX file examplesPath as alias)
            end tell
        end if
    else
        display dialog dialogText buttons {"OK"} default button "OK" with title "$APP_NAME" with icon POSIX file iconPath
    end if

end run

on open theFiles
    showAppNotInstalledAlert()
    initialize()

    set filePath to POSIX path of (item 1 of theFiles)
    set appPath to POSIX path of (path to me)
    set synapseqPath to quoted form of (appPath & "Contents/Resources/bin/synapseq")
    
    tell application "System Events"
        set fileName to name of (POSIX file filePath as alias)
    end tell
    
    set mainChoice to button returned of (display dialog "What would you like to do with '" & fileName & "'?" buttons {"Play", "Edit", "More..."} default button "Play" with title "$APP_NAME" with icon POSIX file ((POSIX path of (path to me)) & "Contents/Resources/app_icon.icns"))
    
    set oldDelimiters to AppleScript's text item delimiters
    set AppleScript's text item delimiters to "/"

    set pathComponents to text items of filePath
    set fileNameBase to item -1 of pathComponents
    set dirPath to ""

    repeat with i from 1 to ((count of pathComponents) - 1)
        set dirPath to dirPath & (item i of pathComponents) & "/"
    end repeat

    set AppleScript's text item delimiters to "."

    set fileNameParts to text items of fileNameBase
    set outputFileName to (item 1 of fileNameParts) & ".wav"
    set outputFilePath to quoted form of (dirPath & outputFileName)

    set AppleScript's text item delimiters to oldDelimiters
    
    if mainChoice is "Play" then
        tell application "Terminal"
            activate
            set terminalWindow to do script ("cd " & quoted form of dirPath & "; " & synapseqPath & " " & quoted form of filePath & "; exit")
        end tell
    else if mainChoice is "Edit" then
        tell application "TextEdit"
            activate
            open (POSIX file filePath as alias)
        end tell
    else if mainChoice is "More..." then
        set convertChoice to button returned of (display dialog "Conversion options:" buttons {"Convert to WAV", "Cancel"} default button "Convert to WAV" with title "$APP_NAME" with icon POSIX file ((POSIX path of (path to me)) & "Contents/Resources/app_icon.icns"))
        if convertChoice is "Convert to WAV" then
            tell application "Terminal"
                activate
                set terminalWindow to do script ("cd " & quoted form of dirPath & "; " & synapseqPath & " --output-wav-file " & outputFilePath & " " & quoted form of filePath & "; exit")
                repeat while busy of terminalWindow is true
                    delay 0.5
                end repeat
                tell terminalWindow to close
            end tell
        end if
    end if
end open
EOF

# Compile the AppleScript into an .app
info "Compiling AppleScript into application bundle..."
osacompile -o "$INSTALLER_DIR/$APP_NAME.app" $INSTALLER_DIR/synapseq.applescript > /dev/null

if [ $? -ne 0 ]; then
    error "Failed to compile AppleScript into application bundle!"
    exit 1
fi

# Copy documentation
info "Copying documentation to application bundle..."
create_dir_if_not_exists "$INSTALLER_DIR/$APP_NAME.app/Contents/Resources/docs"
cp "$PWD/ChangeLog.txt" "$INSTALLER_DIR/$APP_NAME.app/Contents/Resources/ChangeLog.txt"

# Copy COPYING.txt
info "Copying COPYING.txt to application bundle..."
cp "$PWD/COPYING.txt" "$INSTALLER_DIR/$APP_NAME.app/Contents/Resources/COPYING.txt"

# Copy NOTICE.txt
info "Copying NOTICE.txt to application bundle..."
cp "$PWD/NOTICE.txt" "$INSTALLER_DIR/$APP_NAME.app/Contents/Resources/NOTICE.txt"

# Convert *.md to *.txt
pandoc -f markdown -t plain "$DOC_DIR/USAGE.md" -o "$INSTALLER_DIR/$APP_NAME.app/Contents/Resources/docs/USAGE.txt"

# Copy examples
info "Copying examples to application bundle..."
create_dir_if_not_exists "$INSTALLER_DIR/$APP_NAME.app/Contents/Resources/samples"
cp -R "$PWD/samples"/* "$INSTALLER_DIR/$APP_NAME.app/Contents/Resources/samples"

# Copy the binary
info "Copying binary to application bundle..."
create_dir_if_not_exists "$INSTALLER_DIR/$APP_NAME.app/Contents/Resources/bin"
cp "$SYNAPSEQ_BINARY" "$INSTALLER_DIR/$APP_NAME.app/Contents/Resources/bin/synapseq"
chmod +x "$INSTALLER_DIR/$APP_NAME.app/Contents/Resources/bin/synapseq"

# Generate the icon from the PNG file
section_header "Generating application icon..."
create_dir_if_not_exists "$INSTALLER_DIR/iconset"
sips -z 16 16 "$PNG_SOURCE" --out "$INSTALLER_DIR/iconset/icon_16x16.png" > /dev/null
sips -z 32 32 "$PNG_SOURCE" --out "$INSTALLER_DIR/iconset/icon_16x16@2x.png" > /dev/null
sips -z 32 32 "$PNG_SOURCE" --out "$INSTALLER_DIR/iconset/icon_32x32.png" > /dev/null
sips -z 64 64 "$PNG_SOURCE" --out "$INSTALLER_DIR/iconset/icon_32x32@2x.png" > /dev/null
sips -z 128 128 "$PNG_SOURCE" --out "$INSTALLER_DIR/iconset/icon_128x128.png" > /dev/null
sips -z 256 256 "$PNG_SOURCE" --out "$INSTALLER_DIR/iconset/icon_128x128@2x.png" > /dev/null
sips -z 256 256 "$PNG_SOURCE" --out "$INSTALLER_DIR/iconset/icon_256x256.png" > /dev/null
sips -z 512 512 "$PNG_SOURCE" --out "$INSTALLER_DIR/iconset/icon_512x512.png" > /dev/null
sips -z 1024 1024 "$PNG_SOURCE" --out "$INSTALLER_DIR/iconset/icon_512x512@2x.png" > /dev/null

# Convert the .iconset to .icns
mv "$INSTALLER_DIR/iconset" "$INSTALLER_DIR/$ICON_NAME.iconset"
iconutil -c icns "$INSTALLER_DIR/$ICON_NAME.iconset"

# Move the .icns to Resources
mv "$INSTALLER_DIR/$ICON_NAME.icns" "$INSTALLER_DIR/$APP_NAME.app/Contents/Resources/"

# Add file association and icon reference to Info.plist
info "Configuring application bundle Info.plist..."
cat > "$INSTALLER_DIR/$APP_NAME.app/Contents/Info.plist" <<EOF
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>CFBundleExecutable</key>
    <string>droplet</string>
    <key>CFBundleIdentifier</key>
    <string>com.synapseq</string>
    <key>CFBundleName</key>
    <string>$APP_NAME</string>
    <key>CFBundlePackageType</key>
    <string>APPL</string>
    <key>CFBundleVersion</key>
    <string>1.0</string>
    <key>CFBundleIconFile</key>
    <string>app_icon</string>
    <key>CFBundleDocumentTypes</key>
    <array>
        <dict>
            <key>CFBundleTypeName</key>
            <string>SPSQ File</string>
            <key>CFBundleTypeExtensions</key>
            <array>
                <string>spsq</string>
            </array>
            <key>CFBundleTypeRole</key>
            <string>Editor</string>
            <key>LSItemContentTypes</key>
            <array>
                <string>com.synapseq.spsq</string>
            </array>
            <key>CFBundleTypeOSTypes</key>
            <array>
                <string>spsq</string>
            </array>
            <key>LSHandlerRank</key>
            <string>Owner</string>
            <key>LSIsAppleDefaultForType</key>
            <true/>
            <key>CFBundleTypeIconFile</key>
            <string>app_icon</string>
            <key>CFBundleTypeMIMETypes</key>
            <array>
                <string>application/x-synapseq</string>
            </array>
        </dict>
    </array>
    <key>UTImportedTypeDeclarations</key>
    <array>
        <dict>
            <key>UTTypeIdentifier</key>
            <string>com.synapseq.spsq</string>
            <key>UTTypeDescription</key>
            <string>SynapSeq Sequence File</string>
            <key>UTTypeConformsTo</key>
            <array>
                <string>public.data</string>
            </array>
            <key>UTTypeTagSpecification</key>
            <dict>
                <key>public.filename-extension</key>
                <array>
                    <string>spsq</string>
                </array>
                <key>public.mime-type</key>
                <array>
                    <string>application/x-synapseq</string>
                </array>
            </dict>
            <key>UTTypeIconFile</key>
            <string>app_icon</string>
        </dict>
    </array>
    <key>UTExportedTypeDeclarations</key>
    <array>
        <dict>
            <key>UTTypeIdentifier</key>
            <string>com.synapseq.spsq</string>
            <key>UTTypeDescription</key>
            <string>SynapSeq Sequence File</string>
            <key>UTTypeConformsTo</key>
            <array>
                <string>public.data</string>
            </array>
            <key>UTTypeTagSpecification</key>
            <dict>
                <key>public.filename-extension</key>
                <array>
                    <string>spsq</string>
                </array>
                <key>public.mime-type</key>
                <array>
                    <string>application/x-synapseq</string>
                </array>
            </dict>
            <key>UTTypeIconFile</key>
            <string>app_icon</string>
        </dict>
    </array>
    <key>NSAppleEventsUsageDescription</key>
    <string>This application needs to control the Terminal to process .spsq files.</string>
</dict>
</plist>
EOF

# Create an entitlements file to allow Apple Events
cat > "$INSTALLER_DIR/entitlements.plist" <<EOF
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>com.apple.security.automation.apple-events</key>
    <true/>
</dict>
</plist>
EOF

# Sign the application with ad-hoc signing
codesign --force --deep --sign - "$INSTALLER_DIR/$APP_NAME.app" --entitlements "$INSTALLER_DIR/entitlements.plist" > /dev/null 2>&1

# Create a temporary folder for the DMG content
section_header "Creating DMG installer..."
create_dir_if_not_exists "$INSTALLER_DIR/dmg"
mv "$INSTALLER_DIR/$APP_NAME.app" "$INSTALLER_DIR/dmg/"

# Clean all dmg files in dist
rm -f $BUILD_DIR/dist/*.dmg

# Create the DMG with background using create-dmg
create-dmg \
  --volname "$APP_NAME" \
  --window-pos 200 120 \
  --window-size 800 400 \
  --icon-size 100 \
  --icon "$APP_NAME.app" 200 190 \
  --hide-extension "$APP_NAME.app" \
  --app-drop-link 600 185 \
  --background "$BUILD_DIR/assets/dmg-background.png" \
  "$BUILD_DIR/dist/$DMG_NAME" \
  "$INSTALLER_DIR/dmg" > /dev/null

if [ $? -ne 0 ]; then
    error "Failed to create DMG!"
    exit 1
fi

# Configure the DMG icon
DeRez -only icns "$BUILD_DIR/assets/VolumeIcon.icns" > "$INSTALLER_DIR/icns.rsrc"
Rez -append "$INSTALLER_DIR/icns.rsrc" -o "$BUILD_DIR/dist/$DMG_NAME"

# Set the DMG icon
SetFile -a C "$BUILD_DIR/dist/$DMG_NAME"

# Remove temporary files
info "Cleaning up temporary files..."
rm -rf $INSTALLER_DIR

success "Application bundle created and packaged in $BUILD_DIR/dist/$DMG_NAME successfully!"