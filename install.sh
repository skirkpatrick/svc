#!/usr/bin/env sh

echo "Welcome to SVC version control system installer."
echo "This script is going to install and set up SVC version control system."

which go > /dev/null
if [ $? -ne 0 ]; then
  echo "Go is not installed on your machine. To install SVC version control " \
       "system Go must be installed: Go is an open source programming environment " \
       "that makes it easy to build simple, reliable, and efficient software."

  answer=""
  while [ "$answer" != "Y" -a "$answer" != "N" ]; do
    echo "Do you want to install Go? (Y/N)"
    read answer
    answer=`echo $answer | tr yn YN`
  done

  if [ "$answer" == "Y" ]; then
    echo "Installing Go"
    wget https://raw.github.com/skirkpatrick/GetGo/master/getgo.sh
    source getgo.sh
    rm getgo.sh
    echo "Done Installing Go"
  else
    echo "Error: Cannot install SVC without installing Go"
    exit 1
  fi
else
  echo "Go is installed"
fi

echo "Installing SVC"
go get github.com/skirkpatrick/svc

echo "Done Installing SVC"
echo "For more information about SVC version control system type:"
echo "svc help"
