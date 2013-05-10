#Simple Version Control

[![Build Status](https://travis-ci.org/skirkpatrick/svc.png)](https://travis-ci.org/skirkpatrick/svc)

##What is Simple Version Control?

Simple Version Control (SVC for short) is a source control system aimed at those new to source control. SVC is intended to be simple and intuitive so the user can quickly grasp the basic concepts of developing with version control, without the hassel of learning complicated commands or the inner workings of the system.

##What the heck *is* version control?

Since you're on Github, chances are you already know the answer to this question. Version control, also known as source control or revision control, is software that keeps track of changes to files. It is typically used for tracking changes to source code and most source control systems provide much more functionality than simply keeping track of changes, such as allowing the user to rollback changes or work on multiple features at the same time.

##Why not just use Git?

Git is an extremely powerful and customizable tool that I feel every programmer should know how to use. Unfortunately, Git also has a very steep learning curve, especially for someone who has never used source control before. That's where SVC comes in. SVC is meant to give users an introduction to version control and allow them to become comfortable integrating version control into their workflow before moving on to a more powerful version control system, such as Git.

##Contributing

Find bugs, take a look at the issues, create new issues, submit pull requests. I'm open to any and all suggestions!

##Getting started

###Installation

Currently, only Linux is supported, though I'm pretty sure it will work fine on Mac OS X, too.

If you don't have Go installed, you need it. I recommend using my [install script](https://github.com/skirkpatrick/GetGo), or you can do it manually.

Assuming you have [Go](http://golang.org/doc/install) installed and have a workspace setup through [GOPATH](http://golang.org/doc/code.html):

`$ go get github.com/skirkpatrick/svc`


###Basic usage

Run `$ svc help` to see the full list of commands

All commands begin with `svc` followed by the command name and any arguments, if applicable.

####Create a repo

`$ svc init`

####Commit changes

`$ svc commit`

####Check status since last commit

`$ svc status`

####Reset to state of last commit

`$ svc reset`

####Revert 5 commits

`$ svc revert 5`

####View change log

`$ svc log`

####View list of all branches

`$ svc branch`

####Create a new branch or checkout existing branch

`$ svc branch branch_name`

Branch names can contain any characters that folders/files can, including white space

####Delete repo

`$ svc delete`

Deleting a repo will leave the working directory's contents in their current state
