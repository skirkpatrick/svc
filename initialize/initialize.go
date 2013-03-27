/*
    Functions to initialize empty SVC repository in current directory
*/
package initialize

import (
    "fmt"
    "github.com/skirkpatrick/svc/dirutils"
    "github.com/skirkpatrick/svc/meta"
)

// Initializes an empty repo in current directory.
func Initialize() {
    curDir, file, err := dirutils.GetCurrentDirectory()
    if err != nil { panic(err) }

    // Check if already in repo
    exists, dir := dirutils.RecursivelyCheckForRepo(file)
    if exists {
        fmt.Printf("Found existing SVC repo in %s\n", dir)
        return
    }

    fmt.Printf("Initializing empty repo in %s\n", curDir)

    // Make .svc directory and create metadata file
    dirutils.InitializeRepo()
    meta.InitializeRepo()
}
