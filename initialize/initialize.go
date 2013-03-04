/*
    Functions to initialize empty SVC repository in current directory
*/
package initialize

import (
    "fmt"
    "os"
    "github.com/skirkpatrick/svc/dirutils"
)

// Initializes an empty repo in current directory.
func Initialize() {
    curDir, err := os.Getwd()
    if err != nil { panic(err) }
    file, err := os.Open(curDir)
    if err != nil { panic(err) }

    // Check if already in repo
    exists, dir := dirutils.RecursivelyCheckForRepo(file)
    if exists {
        fmt.Printf("Found existing SVC repo in %s\n", dir)
    }

    fmt.Printf("Initializing empty repo in %s\n", curDir)

    // Make .svc directory and create metadata file
    dirutils.InitializeRepo()
}
