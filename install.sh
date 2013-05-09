#!/bin/bash

echo  "Welcomt to svc version control system installer."
echo  "This script is going to install and set up svc version control system."
echo  "Select an operating system:"

select answer in "macosx" "linix"; do
  case $answer in
    macosx) echo  "Note: to install svc version control system GO must be installed:" \
                  "GO is an open source programming environment that makes it easy to build simple, reliable, and efficient software."

            echo  "Is GO installed on your machine?"
            select answer in "Yes" "No"; do
              case $answer in
                Yes)  break;;
                No)   echo  "Do you want to install GO?"
                      select answer in "Yes" "No"; do
                        case $answer in
                          Yes)  echo  "Is Homebrew installed on your machine?"
                                echo  "Note: to install Go, Homebrew must be installed." \
                                echo  "Homebrew isa package management system that simplifies the installation of software on the Mac OS X operating system." \
                                      " It is a free/open source software project to simplify installation of other free/open source software.";
                                select answer in "Yes" "No"; do
                                  case $answer in
                                    Yes)  break;;
                                    No)   echo  "Do you want to install Homebrew?"
                                          select answer in "Yes" "No"; do
                                            case $answer in
                                              Yes)  echo  "installing Homebrew"
                                                    ruby -e "$(curl -fsSL https://raw.github.com/mxcl/homebrew/go)"
                                                    break;;
                                              No)   exit;;
                                            esac
                                          done
                                          break;;
                                  esac
                                done
                                echo "installing go"
                                brew install go
                                echo "installing svc"
                                cd svc
                                go get github.com/skirkpatrick/svc/initialize
                                go build
                                cp svc /usr/local/bin
                                rm -rf svc
                                break;;
                          No)   exit;;
                        esac
                      done
                      break;;
              esac
            done
            break;;
    linix)  echo "Not yet implemented"
            exit 
  esac
done

echo "Done installing"
echo "For more information about svc version control system type:"
echo "svc help"
