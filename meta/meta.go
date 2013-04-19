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
    "github.com/skirkpatrick/svc/dirutils"
)

const (
    metafileName = "metadata"
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

// WriteMetadata writes repo to metadata file "file"
func WriteMetadata(file *os.File, repo *Repo) error {
    xml_raw, err := xml.Marshal(repo)
    if err != nil { return err }

    // This looks too much like C...must be a more elegant way
    _, err = file.Write(append([]byte(xml.Header), xml_raw...))
    return err
}

// InitializeRepo initializes metadata file for a new repo
func InitializeRepo() {
    // Create metadata file
    file, err := os.Create(dirutils.ObjectDir + "/" + metafileName)
    if err != nil { panic(err) }
    defer file.Close()

    // Initialize metadata file
    InitializeMetafile(file)
}

// InitializeMetafile writes initial metadata to file
func InitializeMetafile(file *os.File) {
    repo := new(Repo)
    branch := new(Branch)
    branch.Title = "master"
    repo.AddBranch(branch)
    err := repo.SetCurrent("master")
    if err != nil { panic(err) }
    err = WriteMetadata(file, repo)
    if err != nil { panic(err) }
}

// Open returns the metadata file for the current repo (if it exists)
func Open() (*Repo, error) {
    // Find file location
    curDir, file, err := dirutils.GetCurrentDirectory()
    if err != nil { return nil, err }
    defer file.Close()

    exists, dir := dirutils.RecursivelyCheckForRepo(file)
    if !exists {
        err = fmt.Errorf("No existing SVC repo found in %s", curDir)
        return nil, err
    }

    filename := dir + "/" + dirutils.ObjectDir + "/" + metafileName
    repo := ReadMetadata(filename)
    return repo, nil
}


/// Begin repo functions ///

// Write writes a repo to the metadata file
func (repo *Repo) Write() error {
    dir, err := dirutils.OpenRepo()
    if err != nil { return err }
    defer dir.Close()
    filename := dir.Name() + "/" + dirutils.ObjectDir + "/" + metafileName
    file, err := os.Create(filename)
    if err != nil { return err }
    defer file.Close()
    err = WriteMetadata(file, repo)
    return err
}

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


// DeleteBranch removes an entire branch from a repo
func (repo *Repo) DeleteBranch(branchname string) error {
    if len(repo.Branch) == 1 {
        return fmt.Errorf("cannot remove the last branch of a repo")
    }
    _, i := repo.Find(branchname)
    if i == -1 {
        return fmt.Errorf("branch %q does not exist", branchname)
    }
    copy(repo.Branch[i:], repo.Branch[i+1:])
    repo.Branch[len(repo.Branch)-1] = Branch{}
    repo.Branch = repo.Branch[:len(repo.Branch)-1]
    return nil
}

// Find finds a Branch in a repo, if it exists.
// If the branch is found, it is returned along with its
// position in the []Branch.
// If the branch is not found, nil and -1 are returned.
func (repo *Repo) Find(branchname string) (branch *Branch, pos int) {
    for i, b := range repo.Branch {
        if b.Title == branchname {
            return &repo.Branch[i], i
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


// Branch.DeleteCommit deletes a commit from the branch
func (branch *Branch) DeleteCommit(i int) error {
    if i >= len(branch.Commit) {
        return fmt.Errorf("commit %d does not exist in branch %q", i, branch.Title)
    }
    copy(branch.Commit[i:], branch.Commit[i+1:])
    branch.Commit[len(branch.Commit)-1] = Commit{}
    branch.Commit = branch.Commit[:len(branch.Commit)-1]
    return nil
}


// Branch.DeleteLastCommit deletes the last commit from the branch
func (branch *Branch) DeleteLastCommit() error {
    return branch.DeleteCommit(len(branch.Commit)-1)
}


/// Begin Commit functions ///

// Commit.Copy() returns an exact copy of a Commit.
func (commit *Commit) Copy() *Commit {
    newCommit := new(Commit)
    newCommit.XMLName = commit.XMLName
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
