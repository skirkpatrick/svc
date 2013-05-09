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
package remove


import (
    "os"
    "fmt"
    "bufio"
    "github.com/skirkpatrick/svc/dirutils"
)


// Remove deletes the current repo metadata and stash
func Remove() {
    if prompt() {
        err := dirutils.RemoveRepo()
        if err != nil {
            fmt.Println(err)
            return
        }
        fmt.Println("Repo successfully removed.")
    } else {
        fmt.Println("Repo removal aborted.")
    }
}


// prompt prompts the user for approval to delete the repo
func prompt() bool {
    c := make([]byte, 1)
    stdin := bufio.NewReader(os.Stdin)
    fmt.Println("Working files will not be changed or removed.")
    fmt.Println("Are you sure you wish to delete this repo?")
    fmt.Print("[y/n]: ")
    stdin.Read(c)
    return c[0] == 'y'
}
