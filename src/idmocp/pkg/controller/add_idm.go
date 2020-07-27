package controller

import (
	"idmocp/pkg/controller/idm"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, idm.Add)
}
