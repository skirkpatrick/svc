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
package crypto

import (
    "fmt"
    "io/ioutil"
    "crypto/sha512"
)

//SHA512 computes a SHA512 of the contents of the file with the given name.
func SHA512 (filename string) (string, error) {
    // Read file contents
    contents, err := ioutil.ReadFile(filename)
    if err != nil {
        return "", err
    }

    // Compute SHA512
    c := sha512.New()
    _, err = c.Write(contents)
    if err != nil {
        return "", err
    }
    result := fmt.Sprintf("%x", c.Sum(nil))
    return result, nil
}
