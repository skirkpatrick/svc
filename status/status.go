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
)


// Status prints the current repo's status
func Status() {
    // Read metadata
    repo, err := meta.Open()
    if err != nil { panic(err) }

    // Open base directory
    repoDir, err := dirutils.OpenRepo()
    if err != nil { panic(err) }

    printStatus(repo, repoDir)
}


// GetIgnore reads the .svcignore file
func GetIgnore () {
    // TODO
}


// GetFileStatus maps file names to their status.
func GetFileStatus(repo *meta.Repo, dir *os.File, prefix string) map[string]int {
    // GetIgnore
    files, err := dirutils.GetDirectoryContents(dir)
    if err != nil { panic(err) }

    filemap := make(map[string]int, len(files))

    // Recursively check directories (can be made concurrent)
    for _, file := range files {
        if file.IsDir() {
            recurseDir, err := os.Open(prefix + file.Name())
            if err != nil { panic(err) }
            stats := GetFileStatus(repo, recurseDir, prefix + file.Name() + "/")
            for key, value := range stats {
                filemap[key] = value
            }
        } else {
            // Check if modified or ignored
            filemap[prefix + file.Name()] = MODIFIED
        }
    }
    return filemap
}


// printStatus prints the repo's current status
func printStatus(repo *meta.Repo, repoDir *os.File) {
    files := GetFileStatus(repo, repoDir, "")
    for key, value := range files {
        fmt.Printf("%s: %v\n", key, value)
    }
}
