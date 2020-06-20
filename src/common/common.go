/*
	common function calls
*/

package common

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"regexp"
	"time"

	"github.com/gorilla/mux"
	"github.com/stianeikeland/go-rpio"
	"gitlab.com/bobymcbobs/go-rpi-gpio-api/src/types"
)

var (
	AppBuildVersion = "0.0.0"
	AppBuildHash = "???"
	AppBuildDate = "???"
	AppBuildMode = "development"
)

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

func GetEnvOrDefault(envName string, defaultValue string) (output string) {
	output = os.Getenv(envName)
	if output == "" {
		output = defaultValue
	}
	return output
}

// determine the port for the app to run on
func GetAppPort() (output string) {
	return GetEnvOrDefault("APP_PORT", ":8080")
}

// determine the tls port for the app to run on
func GetAppPortTLS() (output string) {
	return GetEnvOrDefault("APP_PORT_TLS", ":4433")
}

// determine if the app should host with TLS
func GetAppUseTLS() (output string) {
	return GetEnvOrDefault("APP_USE_TLS", "false")
}

// determine path to the public SSL cert
func GetAppTLSpublicCert() (output string) {
	return GetEnvOrDefault("APP_TLS_PUBLIC_CERT", "server.crt")
}

// determine path to the private SSL cert
func GetAppTLSprivateCert() (output string) {
	return GetEnvOrDefault("APP_TLS_PRIVATE_CERT", "server.key")
}

func GetAuthSecretFromEnv() string {
	return GetEnvOrDefault("APP_AUTH_SECRET", "")
}

func Logging(next http.Handler) http.Handler {
	// log all requests
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%v %v %v %v %v", r.Method, r.URL, r.Proto, r.Response, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}

func GetAppDistFolder() string {
	appDistFolder := GetEnvOrDefault("APP_DIST_FOLDER", "./dist")
	return appDistFolder
}

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

func GetHTTPresponseBodyContents(response *http.Response) (output types.JSONMessageResponse) {
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal(responseData, &output)
	return output
}

func HTTPuseMiddleware(handler http.HandlerFunc, middlewares ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	for _, middleware := range middlewares {
		handler = middleware(handler)
	}

	return handler
}

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

func OpenGPIOpins() (err error) {
	if err := rpio.Open(); err != nil {
		return err
	}
	return err
}

func ValidatePinNumber(num int) (err error) {
	if !(num <= 40 && num >= 0) {
		return errors.New("Invalid pin number")
	}
	return err
}

func GetPin(num int) (pin types.Pin, err error) {
	err = ValidatePinNumber(num)
	if err != nil {
		return pin, err
	}
	err = OpenGPIOpins()
	if err != nil {
		return pin, err
	}
	defer rpio.Close()
	pinSelect := rpio.Pin(num)
	state := pinSelect.Read()

	pin = types.Pin{
		Number: num,
		State:  int(state),
	}
	return pin, err
}

func WritePin(num int, state int, mode int) (pin types.Pin, err error) {
	err = ValidatePinNumber(num)
	if err != nil {
		return pin, err
	}
	err = OpenGPIOpins()
	if err != nil {
		return pin, err
	}
	defer rpio.Close()
	pinSelect := rpio.Pin(num)
	switch mode {
	case 0:
		pinSelect.Mode(rpio.Input)
	case 1:
		pinSelect.Mode(rpio.Output)
	default:
		return types.Pin{}, errors.New("Invalid mode - valid options: 0 (input), 1 (output)")
	}
	switch state {
	case 0:
		pinSelect.Write(rpio.Low)
		// pinSelect.Low()
	case 1:
		pinSelect.Write(rpio.High)
		// pinSelect.High()
	default:
		return types.Pin{}, errors.New("Invalid pin number - valid options: 0 (low), 1 (high)")
	}

	pin = types.Pin{
		Number: num,
		State:  state,
	}
	return pin, err
}

func ListPins() (pinList types.PinList, err error) {
	for num := 1; num <= 40; num++ {
		pinState, err := GetPin(num)
		if err != nil {
			return types.PinList{}, err
		}
		pinList = append(pinList, pinState)
	}
	return pinList, err
}

func CheckGPIOpinForGround(r *http.Request) (matches bool) {
	vars := mux.Vars(r)
	pinIdstr := vars["pin"]
	matches, _ = regexp.MatchString(`\b(([1-5])|(([7-8]))|(([0-1]|1[0-3]))|(([0-1]|1[5-9]))|(([2]|2[1-4]))|(([2]|2[6-9]))|(([2]|2[7-8]))|(([3]|3[1-3]))|(([3]|3[5-8]))|40)\b`, pinIdstr)
	return !matches
}

func CheckForValidState(r *http.Request) (matches bool) {
	vars := mux.Vars(r)
	pinIdstr := vars["state"]
	matches, _ = regexp.MatchString(`\b(0|1)\b`, pinIdstr)
	return matches
}
