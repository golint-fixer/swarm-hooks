package authZ

import (
	"net/http"
)

//HandleAuthZAPI - API for backend ACLs services - for now only tenant seperation - finer grained later
type HandleAuthZAPI interface {

	//The Admin should first provision itself before starting to servce
	Init() error

	HandleEvent(eventType EventEnum, w http.ResponseWriter, r *http.Request, next http.Handler, containerID string)
}