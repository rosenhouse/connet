package netapi

import "fmt"

type Rule struct {
	AppGuid1 string
	AppGuid2 string
}

func (r Rule) String() string {
	return fmt.Sprintf("%s <--> %s", r.AppGuid1, r.AppGuid2)
}
