/*
    Types and functions for handling metadata file
*/
package meta

import (
    "encoding/xml"
    "time"
    "os"
)

type Repo struct {
    Current Current
    Branch []Branch
}

type Current struct {
    XMLName xml.Name `xml:"current"`
    Branch string `xml:"chardata"`
}

type Branch struct {
    XMLName xml.Name `xml:"branch"`
    Title string `xml:"title"`
    Commit []Commit
}

type Commit struct {
    XMLName xml.Name `xml:"commit"`
    SHA string `xml:"sha512"` //[]byte ?
    Title string `xml:"title"`
    Message string `xml:"message"`
    Timestamp time.Time `xml:"timestamp"`
    File []File
}

type File struct {
    XMLName xml.Name `xml:"file"`
    SHA string `xml:"sha512"` //[]byte ?
    Title string `xml:"title"`
}


func ReadMetadata(file *os.File) Repo {
    // TODO
    return new(Repo)
}

func InitializeMetafile() {
}

func WriteMetadata(file *os.File, repo Repo) {
}
