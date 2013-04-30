package reset

import (
    "github.com/skirkpatrick/svc/revert"
)


// Reset resets the state of the working directory to that
// of the last commit.
func Reset() {
    revert.Revert(0)
}
