/*
    Types and functions for handling metadata file
*/
package meta

import (
    "encoding/xml"
    "time"
    "os"
    "io/ioutil"
    "fmt"
)

type Repo struct {
    XMLName xml.Name `xml:"REPO"`
    Current string `xml:"current"`
    Branch []Branch `xml:"branch"`
}

type Branch struct {
    XMLName xml.Name `xml:"branch"`
    Title string `xml:"title"`
    Commit []Commit `xml:"commit"`
}

type Commit struct {
    XMLName xml.Name `xml:"commit"`
    SHA string `xml:"sha512"`
    Title string `xml:"title"`
    Message string `xml:"message"`
    Timestamp time.Time `xml:"timestamp"`
    File []File `xml:"file"`
}

type File struct {
    XMLName xml.Name `xml:"file"`
    SHA string `xml:"sha512"`
    Title string `xml:"title"`
}


func ReadMetadata(filename string) *Repo {
    // Reading data this way because it's easier
    // for reading a []byte of unkown size
    data, err := ioutil.ReadFile(filename)
    if err != nil { panic(err) }

    repo := new(Repo)
    err = xml.Unmarshal(data, repo)
    if err != nil { panic(err) }
    return repo
}

func InitializeMetafile() {
}

func WriteMetadata(file *os.File, repo Repo) {
}


/// Begin repo functions ///

// SetCurrent sets current branch on repo.
// Returns error if current is not a valid branch in repo.
func (repo *Repo) SetCurrent(current string) error {
    return nil
}

// Find finds a Branch in a repo, if it exists.
// If the branch is found, it is returned along with its
// position in the []Branch.
// If the branch is not found, nil and -1 are returned.
func (repo *Repo) Find(branchname string) (branch *Branch, pos int) {
    for i, b := range repo.Branch {
        if b.Title == branchname {
            return &b, i
        }
    }
    return nil, -1
}

// Print outputs content of repo to standard out.
// Useful for debugging.
func (repo *Repo) Print() {
    fmt.Printf("Current: %q\n", repo.Current)
    for _, branch := range repo.Branch {
        fmt.Printf("Branch: %q\n", branch.Title)
        for _, commit := range branch.Commit {
            fmt.Printf("\tCommit: %q\n", commit.Title)
            fmt.Printf("\tSHA512: %v\n", commit.SHA)
            fmt.Printf("\tMessage: %q\n", commit.Message)
            fmt.Printf("\tTimestamp: %v\n", commit.Timestamp)
            for _, file := range commit.File {
                fmt.Printf("\t\tFile: %q\n", file.Title)
                fmt.Printf("\t\t\t%q SHA: %v\n", file.Title, file.SHA)
            }
        }
    }
}
