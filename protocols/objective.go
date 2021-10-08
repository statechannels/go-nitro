package protocols

import "errors"

type Objective struct {
	Id     string
	Status string
	Scope  []string
}

var ErrNotApproved = errors.New("objective not approved")
var ErrNotInScope = errors.New("channel not in scope of objective")
