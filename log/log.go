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
package log


import (
    "fmt"
    "github.com/skirkpatrick/svc/meta"
)


// DisplayLog will display a log of past commits
func Display() {
    repo, err := meta.Open()
    if err != nil {
        fmt.Println(err)
        return
    }
    branch, _ := repo.Find(repo.Current)
    displayLog(branch)
}


// displayLog displays past commits on a branch
func displayLog(branch *meta.Branch) {
    for _, commit := range branch.Commit {
        displayCommit(&commit)
    }
}


// displayCommit display's a commit's timestamp and title
func displayCommit(commit *meta.Commit) {
    fmt.Printf("%v\n", commit.Timestamp)
    fmt.Println(commit.Title + "\n")
}
