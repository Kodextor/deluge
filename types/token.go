// Steve Phillips / elimisteve
// 2013.04.28

package types

import (
	"github.com/nu7hatch/gouuid"
)

var (
	tokens chan string
)

func init() {
	if tokens == nil {
		tokens = make(chan string)
		// var err error
		go func() {
			var t *uuid.UUID
			var err error
			for {
				t, err = uuid.NewV4()
				if err != nil {
					// TODO: Find a better solution
					tokens <- "SOMETHING_BAD_HAPPENED"
					continue
				}
				tokens <- t.String()
			}
		}()
	}
}

func NewToken() string {
    return <-tokens
}
