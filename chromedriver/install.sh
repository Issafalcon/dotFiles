SCRIPT_DIR=$( cd ${0%/*} && pwd -P )

echo $SCRIPT_DIR

curl "https://storage.googleapis.com/chrome-for-testing-public/130.0.6723.69/linux64/chromedriver-linux64.zip" -o "${SCRIPT_DIR}/chromedriver-linux64.zip"
unzip  "${SCRIPT_DIR}/chromedriver-linux64.zip" -d "${SCRIPT_DIR}/"
sudo mv "${SCRIPT_DIR}/chromedriver-linux64/chromedriver" /usr/bin/chromedriver

rm -fR "${SCRIPT_DIR}/chromedriver-linux64"
rm -f "${SCRIPT_DIR}/chromedriver-linux64.zip"
