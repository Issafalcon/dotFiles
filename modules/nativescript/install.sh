#!/bin/bash
# NativeScript CLI setup for Linux
# https://docs.nativescript.org/setup/linux

set -euo pipefail

npm install -g nativescript

echo ""
echo "✓ NativeScript CLI installed: $(ns --version)"
echo ""
echo "WARNING: Android Studio is not installed automatically."
echo "Install it manually from https://developer.android.com/studio and, during"
echo "setup, select the Android SDK, Android SDK Platform, and Android Virtual"
echo "Device components. ANDROID_HOME is already set to \$HOME/Android/Sdk by"
echo "this module, so installing to the default location needs no extra config."
echo "Once installed, run 'ns doctor android' to verify the setup."
