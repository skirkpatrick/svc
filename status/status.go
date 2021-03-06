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

    branch, _ := repo.Find(repo.Current)
    if branch == nil { panic(fmt.Errorf("Repo is corrupt")) }

    // Check if initial commit
    commit := (*meta.Commit) (nil)
    numCommits := len(branch.Commit)
    if numCommits > 0 {
        commit = &branch.Commit[numCommits-1]
    }
    return CompareStatusToCommit(commit)
}


// CompareStatusToCommit maps file names to their status with respect
// to a specific commit.
func CompareStatusToCommit(commit *meta.Commit) (map[string]int, error) {
    // Open base directory
    repoDir, err := dirutils.OpenRepo()
    if err != nil {
        return nil, err
    }

    files := readFiles(repoDir, "")

    // Check if modified, unmodified, or new
    markStatus(commit, files)

    return files, nil
}


// readFiles is a recursive function to read all files in repo
func readFiles(dir *os.File, prefix string) map[string]int {
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
            stats := readFiles(recurseDir, prefix + file.Name() + "/")
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
func markStatus(commit *meta.Commit, files map[string]int) {
    if commit == nil { return }
    for _, cached := range commit.File {
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
