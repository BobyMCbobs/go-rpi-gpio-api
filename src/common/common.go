/*
	common function calls
*/

package common

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"gitlab.com/bobymcbobs/go-rpi-gpio-api/src/types"
)

// AppVars ...
// defaults which are overridden with build
var (
	AppBuildVersion = "0.0.0"
	AppBuildHash    = "???"
	AppBuildDate    = "???"
	AppBuildMode    = "development"
)

// GetHostname ...
// returns the hostname of the node or pod
func GetHostname() (hostname string) {
	return os.Getenv("HOSTNAME")
}

// GetAppBuildVersion ...
// return the version of the current FlatTrack instance
func GetAppBuildVersion() string {
	return AppBuildVersion
}

// GetAppBuildHash ...
// return the commit which the current FlatTrack binary was built from
func GetAppBuildHash() string {
	return AppBuildHash
}

// GetAppBuildDate ...
// return the build date of FlatTrack
func GetAppBuildDate() string {
	return AppBuildDate
}

// GetAppBuildMode ...
// return the mode that the app is built in
func GetAppBuildMode() string {
	return AppBuildMode
}

// GetEnvOrDefault ...
// return a default value if an environment variable's value is empty
func GetEnvOrDefault(envName string, defaultValue string) (output string) {
	output = os.Getenv(envName)
	if output == "" {
		output = defaultValue
	}
	return output
}

// GetAppPort ...
// determine the port for the app to run on
func GetAppPort() (output string) {
	return GetEnvOrDefault("APP_PORT", ":8080")
}

// GetAppPortTLS ..
// determine the tls port for the app to run on
func GetAppPortTLS() (output string) {
	return GetEnvOrDefault("APP_PORT_TLS", ":4433")
}

// GetAppUseTLS ...
// determine if the app should host with TLS
func GetAppUseTLS() (output string) {
	return GetEnvOrDefault("APP_USE_TLS", "false")
}

// GetAppTLSpublicCert ...
// determine path to the public SSL cert
func GetAppTLSpublicCert() (output string) {
	return GetEnvOrDefault("APP_TLS_PUBLIC_CERT", "server.crt")
}

// GetAppTLSprivateCert ...
// determine path to the private SSL cert
func GetAppTLSprivateCert() (output string) {
	return GetEnvOrDefault("APP_TLS_PRIVATE_CERT", "server.key")
}

// GetAuthSecretFromEnv ...
// return the HTTP auth secret
func GetAuthSecretFromEnv() string {
	return GetEnvOrDefault("APP_AUTH_SECRET", "")
}

// Logging ...
// print request information as it comes through
func Logging(next http.Handler) http.Handler {
	// log all requests
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%v %v %v %v %v", r.Method, r.URL, r.Proto, r.Response, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}

// JSONResponse ...
// handle generic responses
func JSONResponse(r *http.Request, w http.ResponseWriter, code int, output types.JSONMessageResponse) {
	// simpilify sending a JSON response
	output.Metadata.URL = r.RequestURI
	output.Metadata.Timestamp = time.Now().Unix()
	output.Metadata.Version = AppBuildVersion
	output.Metadata.Hostname = GetHostname()
	response, _ := json.Marshal(output)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// HTTPuseMiddleware ...
// handle middleware for requests
func HTTPuseMiddleware(handler http.HandlerFunc, middlewares ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	for _, middleware := range middlewares {
		handler = middleware(handler)
	}

	return handler
}

// CheckAuth ...
// checks if the auth secret is valid
func CheckAuth(r *http.Request, w http.ResponseWriter) (valid bool, err error) {
	authSecret := GetAuthSecretFromEnv()
	if authSecret == "" {
		return true, err
	}

	requestAuthHeader := r.Header.Get("Authorization")
	requestAuthHeaderValueArray := strings.Split(requestAuthHeader, " ")
	if len(requestAuthHeaderValueArray) < 2 {
		return false, err
	}
	requestAuthHeaderValue := requestAuthHeaderValueArray[1]

	if requestAuthHeader == "" || requestAuthHeaderValue == "" || authSecret != requestAuthHeaderValue {
		return false, err
	}

	return true, err
}
