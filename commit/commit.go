package commit

import (
    "fmt"
    "time"
    "encoding/xml"
    "os"
    "io"
    "compress/zlib"
    "strings"
    "github.com/skirkpatrick/svc/meta"
    "github.com/skirkpatrick/svc/status"
    "github.com/skirkpatrick/svc/crypto"
    "github.com/skirkpatrick/svc/dirutils"
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
    err = stashFiles(branch.Title, commit)
    if err != nil {
        branch.DeleteLastCommit()
        fmt.Println(err)
        return
    }
    fmt.Println(commit.Timestamp)
    fmt.Println(commit.Title + "\n")
    fmt.Println(commit.Message)
}


// stashFiles gzips all commited files and stashes the compressed versions
// in a unique folder based on the commit time
func stashFiles(branch string, commit *meta.Commit) error {
    // Make unique directory .svc/BRANCH/COMMIT_TIME
    repoDir, err := dirutils.OpenRepo()
    if err != nil {
        return err
    }
    stashDir := repoDir.Name() + "/" + dirutils.ObjectDir + "/" + branch + "/" + commit.Timestamp.String()
    err = os.MkdirAll(stashDir, dirutils.Permissions)
    if err != nil {
        return err
    }
    // zlib compress all files
    origDir, err := os.Getwd()
    if err != nil {
        os.RemoveAll(stashDir) // At this point, screw the error
        return err
    }
    defer os.Chdir(origDir)
    buffer := make([]byte, 1024)
    for _, file := range commit.File {
        fileDir := stashDir + "/" + file.Title
        fileDir = fileDir[:strings.LastIndex(fileDir, "/")]
        err = os.MkdirAll(fileDir, dirutils.Permissions)
        if err != nil {
            os.RemoveAll(stashDir)
            return err
        }
        stashFile, err := os.Create(stashDir + "/" + file.Title)
        if err != nil {
            os.RemoveAll(stashDir)
            return err
        }
        defer stashFile.Close()
        origFile, err := os.Open(file.Title)
        if err != nil {
            os.RemoveAll(stashDir)
            return err
        }
        defer origFile.Close()
        compressor := zlib.NewWriter(stashFile)
        defer compressor.Close()
        for n, err := origFile.Read(buffer); n > 0; n, err = origFile.Read(buffer) {
            if err != nil && err != io.EOF {
                os.RemoveAll(stashDir)
                return err
            }
            n, err = compressor.Write(buffer)
            if err != nil {
                os.RemoveAll(stashDir)
                return err
            }
        }
    }
    return nil
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
