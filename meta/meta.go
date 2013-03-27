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


// ReadMetadata reads metadata into repo struct
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

// InitializeMetafile creates initial metadata file
func InitializeMetafile(file *os.File) {
    repo := new(Repo)
    branch := new(Branch)
    branch.Title = "master"
    repo.AddBranch(branch)
    err := repo.SetCurrent("master")
    if err != nil { panic(err) }
    WriteMetadata(file, repo)
}

// WriteMetadata writes repo to metadata file "file"
func WriteMetadata(file *os.File, repo *Repo) {
    xml_raw, err := xml.Marshal(repo)
    if err != nil { panic(err) }

    // This looks too much like C...must be a more elegant way
    _, err = file.Write(append([]byte(xml.Header), xml_raw...))
    if err != nil { panic(err) }
}


/// Begin repo functions ///

// SetCurrent sets current branch on repo.
// Returns error if current is not a valid branch in repo.
func (repo *Repo) SetCurrent(current string) error {
    b,_ := repo.Find(current)
    if b == nil {
        return fmt.Errorf("branch %q not found", current)
    }
    repo.Current = b.Title
    return nil
}

// AddBranch adds a new, empty branch to an existing repo.
// If a branch by the same name already exists in the repo,
// an error is returned.
func (repo *Repo) AddBranch(branch *Branch) error {
    b,_ := repo.Find(branch.Title)
    if b != nil {
        return fmt.Errorf("branch %q already exists", branch.Title)
    }
    repo.Branch = append(repo.Branch, *branch)
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


/// Begin Branch functions ///

// Branch.Copy() returns an exact copy of a Branch.
func (branch *Branch) Copy() *Branch {
    newBranch := new(Branch)
    newBranch.XMLName = branch.XMLName
    newBranch.Title = branch.Title
    newBranch.Commit = make([]Commit, len(branch.Commit))
    for i := range newBranch.Commit {
        newBranch.Commit[i] = *branch.Commit[i].Copy()
    }
    return newBranch
}


/// Begin Commit functions ///

// Commit.Copy() returns an exact copy of a Commit.
func (commit *Commit) Copy() *Commit {
    newCommit := new(Commit)
    newCommit.XMLName = commit.XMLName
    newCommit.SHA = commit.SHA
    newCommit.Title = commit.Title
    newCommit.Message = commit.Message
    newCommit.Timestamp = commit.Timestamp
    newCommit.File = make([]File, len(commit.File))
    for i := range newCommit.File {
        newCommit.File[i] = *commit.File[i].Copy()
    }
    return newCommit
}


/// Begin File Functions ///

// File.Copy() returns an exact copy of a meta.File.
func (file *File) Copy() *File {
    newFile := new(File)
    newFile.XMLName = file.XMLName
    newFile.SHA = file.SHA
    newFile.Title = file.Title
    return newFile
}
