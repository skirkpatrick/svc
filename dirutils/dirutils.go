/*
    Utilities for managing directories/files
*/

package dirutils

import (
    "fmt"
    "os"
)

const (
    objectDir = ".svc"
)

// Check parent directories until root for existing repo.
// Probably a more elegant way to do this.
func RecursivelyCheckForRepo(file *os.File) (found bool, dir string) {
    found = false
    origDir := file.Name()
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

    // Change back to original directory
    err := os.Chdir(origDir)
    if err != nil { panic(err) }

    return
}


// Check current directory for existing repo.
func checkForRepo(file *os.File) bool {
    names, err := file.Readdirnames(0)
    if err != nil { panic(err) }
    for _, name := range names {
        if name == objectDir { return true }
        fmt.Printf("Reading %s\n", name)
    }
    return false
}
