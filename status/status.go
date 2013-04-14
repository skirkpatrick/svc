package status

import (
    "os"
    "fmt"
    "github.com/skirkpatrick/svc/dirutils"
    "github.com/skirkpatrick/svc/meta"
    "github.com/skirkpatrick/svc/crypto"
)


const (
    UNMODIFIED int = 0
    MODIFIED int = 1
    UNTRACKED int = 2
    NEW int = 3
    REMOVED int = 4
)


// Status prints the current repo's status
func Status() {
    printStatus()
}


// GetIgnore reads the .svcignore file
func GetIgnore (dir *os.File) map[string]bool {
    // TODO
    return make(map[string]bool)
}


// GetFileStatus maps file names to their status.
func GetFileStatus() (map[string]int, error) {
    // Read metadata
    repo, err := meta.Open()
    if err != nil {
        return nil, err
    }

    // Open base directory
    repoDir, err := dirutils.OpenRepo()
    if err != nil { panic(err) }

    files := readFiles(repo, repoDir, "")

    // Check if modified, unmodified, or new
    markStatus(repo, files)

    return files, nil
}


// readFiles is a recursive function to read all files in repo
func readFiles(repo *meta.Repo, dir *os.File, prefix string) map[string]int {
    files, err := dirutils.GetDirectoryContents(dir)
    if err != nil { panic(err) }

    filemap := make(map[string]int, len(files))

    // Recursively check directories (can be made concurrent)
    for _, file := range files {
        if file.IsDir() {
            if file.Name() == dirutils.ObjectDir {
                continue
            }
            recurseDir, err := os.Open(prefix + file.Name())
            if err != nil { panic(err) }
            defer recurseDir.Close()
            stats := readFiles(repo, recurseDir, prefix + file.Name() + "/")
            for key, value := range stats {
                filemap[key] = value
            }
        } else {
            // Check if ignored
            // TODO
            filemap[prefix + file.Name()] = NEW
        }
    }
    return filemap
}


// markStatus marks file status
func markStatus(repo *meta.Repo, files map[string]int) {
    branch, _ := repo.Find(repo.Current)

    // Check if initial commit
    numCommits := len(branch.Commit)
    if numCommits == 0 { return }

    for _, cached := range branch.Commit[numCommits-1].File {
        // compute sha and compare to cached
        title := cached.Title
        workingSHA, err := crypto.SHA512(title)
        switch {
        case err != nil:
            files[title] = REMOVED
        case workingSHA != cached.SHA:
            files[title] = MODIFIED
        default:
            files[title] = UNMODIFIED
        }
    }
}


// getStatusBins separates files based on status
func getStatusBins(files map[string]int) [][]string {
    bins := make([][]string, 5)
    for i := range bins {
        bins[i] = make([]string, 0)
    }

    for key, value := range files {
        bins[value] = append(bins[value], key)
    }

    return bins
}


// IsClean returns true if the file status map indicates no changes to commit
func IsClean(files map[string]int, bins [][]string) bool {
    if bins == nil {
        bins = getStatusBins(files)
    }
    if len(bins[UNMODIFIED]) + len(bins[UNTRACKED]) == len(files) {
        return true
    }
    return false
}


// printStatus prints the repo's current status
func printStatus() {
    files, err := GetFileStatus()
    if err != nil {
        fmt.Println(err)
    }
    bins := getStatusBins(files)
    for i := range bins {
        if len(bins[i]) == 0 { continue }
        switch i {
            case MODIFIED:
                fmt.Println("\x1b[31;1mModified:\x1b[0m")
            case NEW:
                fmt.Println("\x1b[31;1mNew:\x1b[0m")
            case UNTRACKED:
                fmt.Println("\x1b[31;1mUntracked:\x1b[0m")
            case REMOVED:
                fmt.Println("\x1b[31;1mRemoved:\x1b[0m")
            case UNMODIFIED:
                continue;
        }
        for j := range bins[i] {
            fmt.Printf("%s\n", bins[i][j])
        }
    }

    if IsClean(files, bins) {
        fmt.Println("\x1b[32;1mNothing to commit\x1b[0m")
    }
}
