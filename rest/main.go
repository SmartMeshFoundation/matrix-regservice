package rest

import (
	"net/http"

	"github.com/SmartMeshFoundation/matrix-regservice/params"

	"fmt"

	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/ant0ine/go-json-rest/rest"
)

/*
Start the restful server
*/
func Start() {

	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	router, err := rest.MakeRouter(

		/*
			register user by application service
		*/
		rest.Post("/regapp/1/register", RegisterUser),
		/*
			request from home server
		*/
		rest.Get("/regapp/1/rooms/:alias", aliasQuery),
		rest.Get("/regapp/1/users/:user", userQuery),
		rest.Put("/regapp/1/transactions/:txnId", newTransaction),
	)
	if err != nil {
		log.Crit(fmt.Sprintf("maker router :%s", err))
	}
	api.SetApp(router)
	listen := fmt.Sprintf("%s:%d", params.APIHost, params.APIPort)
	log.Crit(fmt.Sprintf("http listen and serve :%s", http.ListenAndServe(listen, api.MakeHandler())))
}
