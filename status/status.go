package status

import (
    "os"
    "fmt"
    "github.com/skirkpatrick/svc/dirutils"
    "github.com/skirkpatrick/svc/meta"
)


const (
    UNMODIFIED int = 0
    MODIFIED int = 1
    UNTRACKED int = 2
    NEW int = 3
)


// Status prints the current repo's status
func Status() {
    printStatus()
}


// GetIgnore reads the .svcignore file
func GetIgnore () {
    // TODO
}


// GetFileStatus maps file names to their status.
func GetFileStatus() map[string]int {
    // Read metadata
    repo, err := meta.Open()
    if err != nil { panic(err) }

    // Open base directory
    repoDir, err := dirutils.OpenRepo()
    if err != nil { panic(err) }

    files := getFileStatus(repo, repoDir, "")

    // Check if modified, unmodified, or new

    return files
}


// getFileStatus is a recursive function to map file names to status
func getFileStatus(repo *meta.Repo, dir *os.File, prefix string) map[string]int {
    files, err := dirutils.GetDirectoryContents(dir)
    if err != nil { panic(err) }

    filemap := make(map[string]int, len(files))

    // Recursively check directories (can be made concurrent)
    for _, file := range files {
        if file.IsDir() {
            recurseDir, err := os.Open(prefix + file.Name())
            if err != nil { panic(err) }
            stats := getFileStatus(repo, recurseDir, prefix + file.Name() + "/")
            for key, value := range stats {
                filemap[key] = value
            }
        } else {
            // Check if ignored
            filemap[prefix + file.Name()] = NEW
        }
    }
    return filemap
}


// printStatus prints the repo's current status
func printStatus() {
    files := GetFileStatus()
    for key, value := range files {
        fmt.Printf("%s: %v\n", key, value)
    }
}
