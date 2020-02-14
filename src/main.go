/*
	initialise the API
*/

package main

import (
	"os"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/stianeikeland/go-rpio"
	"gitlab.com/bobymcbobs/go-rpi-gpio-api/src/common"
	"gitlab.com/bobymcbobs/go-rpi-gpio-api/src/routes"
)

func handleWebserver() {
	// bring up the API
	port := common.GetAppPort()
	appUseTLS := common.GetAppUseTLS()
	appPortTLS := common.GetAppPortTLS()
	appTLSpublicCert := common.GetAppTLSpublicCert()
	appTLSprivateCert := common.GetAppTLSprivateCert()

	router := mux.NewRouter().StrictSlash(true)
	apiEndpointPrefix := "/api"

	for _, endpoint := range routes.GetEndpoints(apiEndpointPrefix) {
		router.HandleFunc(endpoint.EndpointPath, endpoint.HandlerFunc).Methods(endpoint.HttpMethod, http.MethodOptions)
	}

	router.HandleFunc(apiEndpointPrefix+"/{.*}", routes.APIUnknownEndpoint)
	router.HandleFunc(apiEndpointPrefix, routes.APIroot)

	router.Use(common.Logging)

	srv := &http.Server{
		Handler:      router,
		Addr:         port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	if appUseTLS == "true" {
		log.Println("Listening on", appPortTLS)
		srv.Addr = appPortTLS
		log.Fatal(srv.ListenAndServeTLS(appTLSpublicCert, appTLSprivateCert))
	} else {
		log.Println("Listening on", port)
		log.Fatal(srv.ListenAndServe())
	}
}

func main() {
	err := common.OpenGPIOpins()
	if err != nil {
		log.Fatalln("[error] cannot talk to gpio pins")
		os.Exit(1)
	}
	defer rpio.Close()
	// initialise the app
	handleWebserver()
}
