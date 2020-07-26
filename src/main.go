/*
	initialise the API
*/

package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/stianeikeland/go-rpio"
	"gitlab.com/bobymcbobs/go-rpi-gpio-api/src/common"
	"gitlab.com/bobymcbobs/go-rpi-gpio-api/src/pin"
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
	apiEndpointPrefix := ""

	for _, endpoint := range routes.GetEndpoints(apiEndpointPrefix) {
	      router.HandleFunc(endpoint.EndpointPath, endpoint.HandlerFunc).Methods(endpoint.HTTPMethod, http.MethodOptions)
	}

	router.HandleFunc(apiEndpointPrefix+"/{.*}", routes.GetUnknownEndpoint)
	router.HandleFunc(apiEndpointPrefix, routes.GetRoot)

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
	log.Printf("launching gpio-api (%v, %v, %v, %v)\n", common.GetAppBuildVersion(), common.GetAppBuildHash(), common.GetAppBuildDate(), common.GetAppBuildMode())

	err := pin.OpenGPIOpins()
	if err != nil {
		log.Fatalln("[error] cannot talk to gpio pins")
		os.Exit(1)
	}
	defer rpio.Close()
	// initialise the app
	handleWebserver()
}
