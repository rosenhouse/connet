package models

import "errors"

type Rule struct {
	Group1 string `json:"group1"`
	Group2 string `json:"group2"`
}

func (r Rule) Equals(otherRule Rule) bool {
	return r.Group1 == otherRule.Group1 &&
		r.Group2 == otherRule.Group2
}

func (r Rule) Validate() error {
	ok := (r.Group1 != "" && r.Group2 != "")
	if !ok {
		return errors.New("missing required field(s)")
	}
	return nil
}
