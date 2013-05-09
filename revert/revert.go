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
package revert

import (
    "os"
    "strings"
    "fmt"
    "github.com/skirkpatrick/svc/meta"
    "github.com/skirkpatrick/svc/status"
    "github.com/skirkpatrick/svc/dirutils"
    "github.com/skirkpatrick/svc/stash"
)


// Revert reverts the state of the working directory to that
// of n commits ago.
func Revert(n int) {
    repo, err := meta.Open()
    if err != nil {
        fmt.Println(err)
        return
    }
    branch, _ := repo.Find(repo.Current)
    if len(branch.Commit) < n {
        fmt.Printf("Only %d commits exist!\n", len(branch.Commit))
        return
    }
    var commit *meta.Commit
    if len(branch.Commit) == n {
        commit = nil
    } else {
        commit = &branch.Commit[(len(branch.Commit) - 1) - n]
    }
    files, err := status.CompareStatusToCommit(commit)
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
    if commit != nil {
        err = stash.Init(branch.Title, commit.Timestamp.String())
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
    err = cleanup(branch, n)
    if err != nil {
        fmt.Println(err)
        fmt.Println("Repo may be corrupt.")
        return
    }
    err = repo.Write()
    if err != nil {
        fmt.Println(err)
        fmt.Println("Repo may be corrupt.")
        return
    }
    if len(branch.Commit) == 0 {
        fmt.Println("Repo reset to intial state")
    } else {
        fmt.Println("Repo reset to commit:")
        fmt.Println("\x1b[32;1m" + branch.Commit[len(branch.Commit)-1].Title + "\x1b[0m")
        fmt.Printf("%v\n", branch.Commit[len(branch.Commit)-1].Timestamp)
    }
}


// cleanup removes the commits and stash directories after a revert
func cleanup(branch *meta.Branch, n int) error {
    repoBase, err := dirutils.OpenRepo()
    if err != nil {
        return err
    }
    stashBase := repoBase.Name() + "/" + dirutils.ObjectDir + "/" + branch.Title + "/"
    m := len(branch.Commit) - n
    for _, commit := range branch.Commit[m:] {
        err = os.RemoveAll(stashBase + commit.Timestamp.String())
        if err != nil {
            return err
        }
    }
    branch.Commit = branch.Commit[:m]
    return nil
}
