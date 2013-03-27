/*
    Utilities for managing directories/files
*/

package dirutils

import (
    "os"
    "github.com/skirkpatrick/svc/meta"
)

const (
    objectDir = ".svc"
    permissions = 0777
    metafileName = "metadata"
)

// Check parent directories until root for existing repo.
// Probably a more elegant way to do this.
func RecursivelyCheckForRepo(file *os.File) (found bool, dir string) {
    found = false
    origDir := file.Name()
    dir = file.Name()

    // Check for existing repo up to root
    for dir != "/" {
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
    }
    return false
}


// Create new object directory and initialize metadata file
func InitializeRepo() {
    // Create object directory
    err := os.Mkdir(objectDir, permissions)
    if err != nil { panic(err) }

    // Create metadata file
    file, err := os.Create(objectDir + "/" +  metafileName)
    if err != nil { panic(err) }
    defer file.Close()

    // Initialize metadata file
    meta.InitializeMetafile(file)
}
