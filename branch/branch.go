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
package branch

import (
    "fmt"
    "os/exec"
    "github.com/skirkpatrick/svc/meta"
    "github.com/skirkpatrick/svc/status"
    "github.com/skirkpatrick/svc/dirutils"
    "github.com/skirkpatrick/svc/revert"
)


// Branch creates a new branch off of the current branch
// and switches the working branch to the new one.
func Branch(title string) {
    clean, err := isClean()
    if err != nil {
        fmt.Println(err)
        return
    }
    if !clean {
        status.Status()
        fmt.Println("\x1b[31;1mCommit changes before branching\x1b[0m")
        return
    }
    _, err = newBranch(title)
    if err != nil {
        fmt.Println(err)
        return
    }
}


// Display displays a list of all branches
func Display() {
    repo, err := meta.Open()
    if err != nil {
        fmt.Println(err)
        return
    }
    for _, branch := range repo.Branch {
        if branch.Title == repo.Current {
            fmt.Print("*")
        } else {
            fmt.Print(" ")
        }
        fmt.Println(branch.Title)
    }
}


// isClean returns true if the current working directory
// is clean and false otherwise.
func isClean() (bool, error) {
    files, err := status.GetFileStatus()
    if err != nil {
        return false, err
    }
    return status.IsClean(files, nil), nil
}


// newBranch creates a new branch and returns a pointer to it
func newBranch(title string) (*meta.Branch, error) {
    repo, err := meta.Open()
    if err != nil {
        return nil, err
    }
    if repo.SetCurrent(title) == nil {
        err = repo.Write()
        if err != nil {
            return nil, err
        }
        fmt.Println("Switched to branch: \x1b[32;1m" + title + "\x1b[0m")
        revert.Revert(0)
        return nil, nil
    }
    current, _ := repo.Find(repo.Current)
    if current == nil { panic(fmt.Errorf("Repo is corrupt")) }
    branch := current.Copy()
    branch.Title = title
    branch.Parent = current.Title
    current.Child = append(current.Child, branch.Title)
    err = repo.AddBranch(branch)
    if err != nil {
        return nil, err
    }
    err = repo.SetCurrent(branch.Title)
    if err != nil {
        return nil, err
    }
    err = repo.Write()
    if err != nil {
        return nil, err
    }
    err = copyStash(current.Title, branch.Title)
    if err != nil {
        return nil, err
    }
    fmt.Println("Created new branch: \x1b[32;1m" + title + "\x1b[0m")
    return branch, nil
}


// copyStash copies parent's stash directory to child's
func copyStash(parent string, child string) error {
    repoDir, err := dirutils.OpenRepo()
    if err != nil {
        return err
    }
    metaDir := repoDir.Name() + "/" + dirutils.ObjectDir + "/"
    return exec.Command("cp", "-r", metaDir + parent, metaDir + child).Run()
}
