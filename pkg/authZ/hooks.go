package authZ

import (
	"net/http"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/swarm/cluster"
)

//Hooks - Entry point to AuthZ mechanisem
type Hooks struct{}

//TODO  - Hooks Infra for overriding swarm
//TODO  - Take bussiness logic out
//TODO  - Refactor packages
//TODO  - Expand API
//TODO -  Images...
//TODO - https://github.com/docker/docker/pull/15953
//TODO - https://github.com/docker/docker/pull/16331

var authZAPI HandleAuthZAPI
var aclsAPI ACLsAPI

type EventEnum int
type ApprovalEnum int

//PrePostAuthWrapper - Hook point from primary to the authZ mechanisem
func (*Hooks) PrePostAuthWrapper(cluster cluster.Cluster, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		eventType := eventParse(r)
		allowed, containerID := aclsAPI.ValidateRequest(cluster, eventType, w, r)
		//TODO - all kinds of conditionals
		if eventType == passAsIs || allowed == approved || allowed == conditionFilter {
			authZAPI.HandleEvent(eventType, w, r, next, containerID)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Not Authorized!"))
		}
	})
}

func eventParse(r *http.Request) EventEnum {
	log.Debug("Got the uri...", r.RequestURI)

	if strings.Contains(r.RequestURI, "/containers") && (strings.Contains(r.RequestURI, "create")) {
		return containerCreate
	}

	if strings.Contains(r.RequestURI, "/containers/json") {
		return containersList
	}

	if strings.Contains(r.RequestURI, "/containers") &&
		(strings.Contains(r.RequestURI, "logs") || strings.Contains(r.RequestURI, "attach") || strings.Contains(r.RequestURI, "exec")) {
		return streamOrHijack
	}
	if strings.Contains(r.RequestURI, "/containers") && strings.HasSuffix(r.RequestURI, "/json") {
		return containerInspect
	}
	if strings.Contains(r.RequestURI, "/containers") {
		return containerOthers
	}

	if strings.Contains(r.RequestURI, "Will add to here all APIs we explicitly want to block") {
		return notSupported
	}

	return passAsIs
}

//Init - Initialize the Validation and Handling APIs
func (*Hooks) Init() {
	//TODO - should use a map for all the Pre . Post function like in primary.go

	aclsAPI = new(DefaultACLsImpl)
	authZAPI = new(DefaultImp)
	//TODO reflection using configuration file tring for the backend type

	log.Info("Init provision engine OK")
}