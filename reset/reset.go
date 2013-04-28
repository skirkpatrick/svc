package reset

import (
    "os"
    "strings"
    "fmt"
    "github.com/skirkpatrick/svc/meta"
    "github.com/skirkpatrick/svc/status"
    "github.com/skirkpatrick/svc/dirutils"
    "github.com/skirkpatrick/svc/stash"
)


// Reset resets the state of the working directory to that
// of the last commit.
func Reset() {
    repo, err := meta.Open()
    if err != nil {
        fmt.Println(err)
        return
    }
    branch, _ := repo.Find(repo.Current)
    files, err := status.GetFileStatus()
    if err != nil {
        fmt.Println(err)
        return
    }
    repoDir, err := dirutils.OpenRepo()
    if err != nil {
        fmt.Println(err)
        return
    }
    // Might have to restore some files if not at initial commit
    if len(branch.Commit) > 0 {
        err = stash.Init(branch.Title, branch.Commit[len(branch.Commit)-1].Timestamp.String())
        if err != nil {
            fmt.Println(err)
            return
        }
    }
    for file, stat := range files {
        if stat == status.NEW {
            fileWithPath := repoDir.Name() + "/" + file
            err = os.Remove(fileWithPath)
            if err != nil {
                fmt.Println(err)
                fmt.Println("Repo may be corrupt.")
                return
            }
            err = dirutils.DeleteIfEmpty(fileWithPath[:strings.LastIndex(fileWithPath, "/")])
            if err != nil {
                fmt.Println(err)
                fmt.Println("Repo may be corrupt.")
                return
            }
        } else if stat == status.REMOVED || stat == status.MODIFIED {
            err = stash.Restore(file)
            if err != nil {
                fmt.Println(err)
                fmt.Println("Repo may be corrupt.")
                return
            }
        }
    }
    if len(branch.Commit) == 0 {
        fmt.Println("Repo reset to intial state")
    } else {
        fmt.Println("Repo reset to last commit:")
        fmt.Println("\x1b[32;1m" + branch.Commit[len(branch.Commit)-1].Title + "\x1b[0m")
    }
}
