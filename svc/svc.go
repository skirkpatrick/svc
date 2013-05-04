/*
    Simple Version Control

    SVC is meant to be an introductory tool for users new to version control.

    The main goal of SVC is simplicity. To achieve this, the ideology of
    "convention over configuration" is used. SVC will make assumptions with
    its commands to minimize the use of argument flags.

    Though SVC will not be as powerful or as customizable as other source
    control tools such as Git or SVN, it should be very easy to learn and use.

    Since SVC is intended for newer programmers, collaboration tools are not
    a priority. Instead, SVC will focus on getting users used to the idea of
    commiting, branching, and using common source control tools for solo
    projects.
*/
package main

import (
    "fmt"
    "os"
    "strconv"
    "github.com/skirkpatrick/svc/initialize"
    "github.com/skirkpatrick/svc/status"
    "github.com/skirkpatrick/svc/commit"
    "github.com/skirkpatrick/svc/reset"
    "github.com/skirkpatrick/svc/revert"
    "github.com/skirkpatrick/svc/remove"
    "github.com/skirkpatrick/svc/log"
    "github.com/skirkpatrick/svc/branch"
)

const (
    //Add new commands to this list
    helpMessage = "Usage: %s <command> [<flags>]\n\nValid commands are:\n" +
                  "init\t\tInitialize new SVC repo in current directory\n" +
                  "commit\t\tCommit current changes\n" +
                  "status\t\tList branch status\n" +
                  "branch\t\tList branches\n" +
                  "branch <b>\tCreate new feature branch\n" +
                  "merge <b>\tMerge branch <b> into this branch\n" +
                  "log\t\tList commit history\n" +
                  "reset\t\tReset all current changes\n" +
                  "revert <v>\tRevert to <v> commits ago\n" +
                  "delete\t\tRemove SVC repo in current directory\n" +
                  "tutorial\tLearn to use SVC\n" +
                  "help\t\tDisplay this message\n"
)

func main() {
    if len(os.Args) < 2 {
        displayHelp()
        return
    }

    //Check command
    //Add new commands to this list
    switch os.Args[1] {
        case "init":
            initialize.Initialize()
        case "commit":
            commit.Commit()
        case "status":
            status.Status()
        case "branch":
            if len(os.Args) == 3 {
                branch.Branch(os.Args[2])
            } else {
                branch.Display()
            }
        case "merge":
        case "log":
            log.Display()
        case "reset":
            reset.Reset()
        case "revert":
            if len(os.Args) != 3 {
                displayHelp()
                return
            }
            n, err := strconv.Atoi(os.Args[2])
            if err != nil {
                fmt.Println(err)
                os.Exit(2)
            }
            revert.Revert(n)
        case "delete":
            remove.Remove()
        case "tutorial":
        case "help":
            displayHelp()
        default:
            displayHelp()
    }
}


func displayHelp() {
    fmt.Printf(helpMessage, os.Args[0])
}
