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
)

const (
    helpMessage = "Usage: %s [command]\n\nValid commands are:\n" +
                  "commit\t\tCommit current changes\n" +
                  "help\t\tDisplay this message\n"
)

func main() {
    if len(os.Args) < 2 {
        displayHelp()
        return
    }

    //Check command
    switch os.Args[1] {
        case "commit":
        case "help":
            displayHelp()
        default:
            displayHelp()
    }
}


func displayHelp() {
    fmt.Printf(helpMessage, os.Args[0])
}
