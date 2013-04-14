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
