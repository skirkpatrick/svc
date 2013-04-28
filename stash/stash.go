package stash

import (
    "compress/zlib"
    "os"
    "io"
    "strings"
    "github.com/skirkpatrick/svc/dirutils"
)


var (
    stashDir string
    buffer []byte
)


// Init initializes the stash directory
func Init(branch string, timestamp string) error {
    // Make unique directory .svc/BRANCH/COMMIT_TIME
    repoDir, err := dirutils.OpenRepo()
    if err != nil {
        return err
    }
    stashDir = repoDir.Name() + "/" + dirutils.ObjectDir + "/" + branch + "/" + timestamp
    err = os.MkdirAll(stashDir, dirutils.Permissions)
    if err != nil {
        return err
    }
    buffer = make([]byte, 1024)
    return nil
}


// Stash stashes a compressed version of a committed file
func Stash(file string) error {
    fileDir := stashDir + "/" + file
    fileDir = fileDir[:strings.LastIndex(fileDir, "/")]
    err := os.MkdirAll(fileDir, dirutils.Permissions)
    if err != nil {
        os.RemoveAll(stashDir)
        return err
    }
    stashFile, err := os.Create(stashDir + "/" + file)
    if err != nil {
        os.RemoveAll(stashDir)
        return err
    }
    defer stashFile.Close()
    workingFile, err := os.Open(file)
    if err != nil {
        os.RemoveAll(stashDir)
        return err
    }
    defer workingFile.Close()
    return compress(workingFile, stashFile)
}


// compress compresses source into dest
func compress(source *os.File, dest *os.File) error {
    compressor := zlib.NewWriter(dest)
    defer compressor.Close()
    for n, err := source.Read(buffer); n > 0; n, err = source.Read(buffer) {
        if err != nil && err != io.EOF {
            os.RemoveAll(stashDir)
            return err
        }
        n, err = compressor.Write(buffer)
        if err != nil {
            os.RemoveAll(stashDir)
            return err
        }
    }
    return nil
}
