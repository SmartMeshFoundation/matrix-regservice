package rest

import (
	"fmt"
	"io/ioutil"

	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/ant0ine/go-json-rest/rest"
)

func aliasQuery(w rest.ResponseWriter, r *rest.Request) {
	printreq(r)
}
func userQuery(w rest.ResponseWriter, r *rest.Request) {
	printreq(r)
}
func newTransaction(w rest.ResponseWriter, r *rest.Request) {
	printreq(r)
}
func printreq(r *rest.Request) {
	log.Trace(fmt.Sprintf("path=%s,args=%s", r.RequestURI, r.PathParams))
	if r.Body != nil {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Error(fmt.Sprintf("read body err %s", err))
		} else {
			log.Trace(fmt.Sprintf("body=%s", string(body)))
		}
	}
}
