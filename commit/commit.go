package commit

import (
    "fmt"
    "time"
    "encoding/xml"
    "github.com/skirkpatrick/svc/meta"
    "github.com/skirkpatrick/svc/status"
    "github.com/skirkpatrick/svc/crypto"
)


// Commit commits the current changes to the repo
func Commit() {
    repo, err := meta.Open()
    if err != nil {
        fmt.Println(err)
        return
    }
    commit := new(meta.Commit)
    err = fillCommit(commit, repo)
    if err != nil {
        fmt.Println(err)
        return
    }
    branch, _ := repo.Find(repo.Current)
    branch.Commit = append(branch.Commit, *commit)
    err = repo.Write()
    if err != nil { panic(err) }
    //TODO commit contents
    fmt.Println(commit.Timestamp)
    fmt.Println(commit.Title + "\n")
    fmt.Println(commit.Message)
}


// fillCommit fills a commit with all needed information
func fillCommit(commit *meta.Commit, repo *meta.Repo) error {
    files, err := status.GetFileStatus()
    if err != nil {
        return err
    }
    if status.IsClean(files, nil) {
        return fmt.Errorf("\x1b[32;1mNothing to commit\x1b[0m")
    }
    err = commitPrompt(commit)
    if err != nil {
        return err
    }
    //TODO figure out or get rid of commit sha
    commit.SHA = ""
    commit.Timestamp = time.Now()
    commitFiles(commit, repo, files)
    return nil
}


// commitPrompt prompts the user for the commit title and message
func commitPrompt(commit *meta.Commit) error {
    err := promptTitle(commit)
    if err != nil {
        return err
    }
    err = promptMessage(commit)
    return err
}


// promptTitle prompts the user for the commit title
func promptTitle(commit *meta.Commit) error {
    //TODO
    commit.Title = "Temporary title"
    return nil
}


// promptMessage prompts the user for the commit message
func promptMessage(commit *meta.Commit) error {
    //TODO
    commit.Message = "Temporary commit message"
    return nil
}


// commitFiles adds file information to commit
func commitFiles(commit *meta.Commit, repo *meta.Repo, files map[string]int) {
    // Preallocation is the way to go
    commit.File = make([]meta.File, 0, len(files))
    for key, value := range files {
        if value == status.UNTRACKED || value == status.REMOVED { continue }
        s, err := crypto.SHA512(key)
        if err != nil { panic(err) }
        commit.File = append(commit.File, meta.File{xml.Name{"", "FILE"}, s, key})
    }
}
