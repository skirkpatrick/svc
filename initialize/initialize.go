/*
    Functions to initialize empty SVC repository in current directory
*/
package initialize

import (
    "fmt"
    "os"
)

// Initializes an empty repo in current directory.
func Initialize() {
    curDir, err := os.Getwd()
    if err != nil { panic(err) }
    file, err := os.Open(curDir)
    if err != nil { panic(err) }

    // Check if already in repo
    exists, dir := recursivelyCheckForRepo(file)
    if exists {
        fmt.Printf("Found existing SVC repo in %s\n", dir)
    }

    fmt.Printf("Initializing empty repo in %s\n", curDir)
}


// Check parent directories until root for existing repo.
// Probably a more elegant way to do this.
func recursivelyCheckForRepo(file *os.File) (found bool, dir string) {
    found = false
    dir = file.Name()

    // Check for existing repo up to root
    for dir != "/" {
        fmt.Printf("In dir %s\n", dir)
        if checkForRepo(file) { return true, dir }

        // cd to parent
        err := os.Chdir("..")
        if err != nil { panic(err) }
        dir, err = os.Getwd()
        if err != nil { panic(err) }
        file, err = os.Open(dir)
        if err != nil { panic(err) }
    }

    // Check root
    if checkForRepo(file) { return true, dir }

    return
}


// Check current directory for existing repo.
func checkForRepo(file *os.File) bool {
    names, err := file.Readdirnames(0)
    if err != nil { panic(err) }
    for _, name := range names {
        if name == ".svc" { return true }
        fmt.Printf("Reading %s\n", name)
    }
    return false
}
