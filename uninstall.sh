#!/usr/bin/env sh

echo "Welcome to SVC version control system uninstaller."
echo "WARNING: This script is going to uninstall SVC version control system."

DIR=$GOPATH/src/github.com/skirkpatrick

answer=""
while [ "$answer" != "Y" -a "$answer" != "N" ]; do
  echo "Do you want to continue? (Y/N)"
  read answer
  answer=`echo $answer | tr yn YN`
done

if [ "$answer" == "Y" ]; then
  echo "Uninstalling SVC"
  echo "To reinstall SVC, run the installation script"
  rm -rf $DIR/svc
  echo "SVC was succesfully uninstalled"
else
  echo "Aborting..."
  exit 1
fi