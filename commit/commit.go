/*
    Simple Version Control
    Copyright (C) 2013  Scott Kirkpatrick

    This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License
    along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/
package commit

import (
    "fmt"
    "time"
    "encoding/xml"
    "os"
    "io/ioutil"
    "strings"
    "github.com/skirkpatrick/svc/meta"
    "github.com/skirkpatrick/svc/status"
    "github.com/skirkpatrick/svc/crypto"
    "github.com/skirkpatrick/svc/stash"
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
    err := stash.Init(branch, commit.Timestamp.String())
    if err != nil {
        return err
    }
    for _, file := range commit.File {
        err = stash.Stash(file.Title)
        if err != nil {
            return err
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
    fmt.Println("To cancel the commit, leave the title empty and press ^D")
    fmt.Println("When finished, press <ENTER>+^D")
    err := promptTitle(commit)
    if err != nil {
        return err
    }
    if commit.Title == "" {
        return fmt.Errorf("Commit aborted due to empty title")
    }
    err = promptMessage(commit)
    return err
}


// promptTitle prompts the user for the commit title
func promptTitle(commit *meta.Commit) error {
    fmt.Println("Enter a title for this commit:")
    buf, err := ioutil.ReadAll(os.Stdin)
    if err != nil {
        return err
    }
    commit.Title = strings.TrimSpace(string(buf))
    return nil
}


// promptMessage prompts the user for the commit message
func promptMessage(commit *meta.Commit) error {
    fmt.Println("Enter a message for this commit:")
    buf, err := ioutil.ReadAll(os.Stdin)
    if err != nil {
        return err
    }
    commit.Message = strings.TrimSpace(string(buf))
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
