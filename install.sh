#!/usr/bin/env sh

echo "Welcomt to SVC version control system installer."
echo "This script is going to install and set up SVC version control system."

which go > /dev/null
if [ $? -ne 0 ]; then
  echo "GO is not installed on your machine. to install SVC version control " \
       "system GO must be installed: GO is an open source programming environment " \
       "that makes it easy to build simple, reliable, and efficient software."

  answer=""
  while [ "$answer" != "Y" -a "$answer" != "N" ]; do
    echo "Do you want to install GO? (Y/N)"
    read answer
    answer=`echo $answer | tr yn YN`
  done

  if [ "$answer" == "Y" ]; then
    echo "Installing GO"
    git clone https://github.com/skirkpatrick/GetGo.git
    GetGo/getgo.sh
    rm -rf GetGo
    echo "Done Installing GO"
  else
    echo "Error: Cannot install SVC without installing GO"
    exit 1
  fi
else
  echo "GO is installed"
fi

echo "Installing SVC"
go get github.com/skirkpatrick/svc
go install
mv svc /usr/local/bin/

echo "Done Installing SVC"
echo "For more information about SVC version control system type:"
echo "svc help"
