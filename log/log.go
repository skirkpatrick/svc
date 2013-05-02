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
