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
    }
    fmt.Println("Repo successfully removed.")
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
