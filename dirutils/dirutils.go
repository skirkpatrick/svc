/*
    Utilities for managing directories/files
*/

package dirutils

import (
    "os"
    "fmt"
)

const (
    ObjectDir = ".svc"
    Permissions = 0777
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
        if name == ObjectDir { return true }
    }
    return false
}


// OpenRepo returns a *File to the current repo base directory
func OpenRepo() (*os.File, error) {
    curDir, dir, err := GetCurrentDirectory()
    if err != nil { return nil, err }
    exists, repoDir := RecursivelyCheckForRepo(dir)
    if !exists { return nil, fmt.Errorf("No repo found in %s", curDir) }
    dir, err = os.Open(repoDir)
    return dir, err
}


// GetCurrentDirectory returns a *File to the current directory
func GetCurrentDirectory() (string, *os.File, error) {
    curDir, err := os.Getwd()
    if err != nil { return curDir, nil, err }
    file, err := os.Open(curDir)
    return curDir, file, err
}


// GetDirectoryContents returns a slice of files/folders in dir
func GetDirectoryContents(dir *os.File) ([]os.FileInfo, error) {
    return dir.Readdir(0)
}


// DeleteIfEmpty deletes a directory if the contents are empty
func DeleteIfEmpty(dir string) error {
    dirp, err := os.Open(dir)
    if err != nil { return err }
    contents, err := dirp.Readdirnames(1)
    if err != nil { return err }
    if len(contents) == 0 {
        err = os.Remove(dir)
    }
    return err
}


// InitializeRepo creates new object directory
func InitializeRepo() {
    err := os.Mkdir(ObjectDir, Permissions)
    if err != nil { panic(err) }
}


// RemoveRepo deletes the current repo's metadata, leaving working files.
func RemoveRepo() error {
    repo, err := OpenRepo()
    if err != nil {
        return err
    }
    defer repo.Close()
    err = os.RemoveAll(repo.Name() + "/" + ObjectDir)
    return err
}
