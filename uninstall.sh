#!/usr/bin/env sh

echo "Welcome to SVC version control system uninstaller."
echo "WARNING: This script is going to uninstall SVC version control system."

GITHUB_DIR=$GOPATH/src/github.com
SKIRKPATRICK_DIR=$GITHUB_DIR/skirkpatrick
SVC_DIR=$SKIRKPATRICK_DIR/svc

answer=""
while [ "$answer" != "Y" -a "$answer" != "N" ]; do
  echo "Do you want to continue? (Y/N)"
  read answer
  answer=`echo $answer | tr yn YN`
done

if [ "$answer" == "Y" ]; then
  echo "Uninstalling SVC"
  echo "To reinstall SVC, run the installation script"
  if [ "`ls $GITHUB_DIR | wc -l`" == "       1" ]; then
    rm -rf $GITHUB_DIR
  elif [ "`ls $SKIRKPATRICK_DIR | wc -l`" == "       1" ]; then
    rm -rf $SKIRKPATRICK_DIR
  else
    rm -rf $SVC_DIR
  fi
  echo "SVC was succesfully uninstalled"
else
  echo "Aborting..."
  exit 1
fi