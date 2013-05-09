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

/*
    Functions to initialize empty SVC repository in current directory
*/
package initialize

import (
    "fmt"
    "github.com/skirkpatrick/svc/dirutils"
    "github.com/skirkpatrick/svc/meta"
)

// Initializes an empty repo in current directory.
func Initialize() {
    curDir, file, err := dirutils.GetCurrentDirectory()
    if err != nil { panic(err) }

    // Check if already in repo
    exists, dir := dirutils.RecursivelyCheckForRepo(file)
    if exists {
        fmt.Printf("Found existing SVC repo in %s\n", dir)
        return
    }

    fmt.Printf("Initializing empty repo in %s\n", curDir)

    // Make .svc directory and create metadata file
    dirutils.InitializeRepo()
    meta.InitializeRepo()
}
