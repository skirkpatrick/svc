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
package stash

import (
    "compress/zlib"
    "os"
    "strings"
    "bytes"
    "github.com/skirkpatrick/svc/dirutils"
)


var (
    repoDir string
    stashDir string
    buffer bytes.Buffer
)


// Init initializes the stash directory
func Init(branch string, timestamp string) error {
    // Make unique directory .svc/BRANCH/COMMIT_TIME if needed
    repoBase, err := dirutils.OpenRepo()
    if err != nil {
        return err
    }
    repoDir = repoBase.Name()
    stashDir = repoDir + "/" + dirutils.ObjectDir + "/" + branch + "/" + timestamp
    err = os.MkdirAll(stashDir, dirutils.Permissions)
    if err != nil {
        return err
    }
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
    workingFile, err := os.Open(repoDir + "/" + file)
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
    _, err := buffer.ReadFrom(source)
    if err != nil {
        return err
    }
    _, err = buffer.WriteTo(compressor)
    return err
}


// Restore restores a stashed file
func Restore(file string) error {
    fileDir := stashDir + "/" + file
    stashFile, err := os.Open(fileDir)
    if err != nil {
        return err
    }
    defer stashFile.Close()
    workingFile, err := os.Create(repoDir + "/" + file)
    if err != nil {
        return err
    }
    defer workingFile.Close()
    return decompress(stashFile, workingFile)
}


// decompress decompresses source into dest
func decompress(source *os.File, dest *os.File) error {
    decompressor, err := zlib.NewReader(source)
    if err != nil {
        return err
    }
    defer decompressor.Close()
    _, err = buffer.ReadFrom(decompressor)
    if err != nil {
        return err
    }
    _, err = buffer.WriteTo(dest)
    return err
}
