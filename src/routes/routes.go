/*
	route related
*/

package routes

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gitlab.com/bobymcbobs/go-rpi-gpio-api/src/common"
	"gitlab.com/bobymcbobs/go-rpi-gpio-api/src/pin"
	"gitlab.com/bobymcbobs/go-rpi-gpio-api/src/types"
)

// GetRoot ...
// webserver root
func GetRoot(w http.ResponseWriter, r *http.Request) {
	// root of API
	JSONresp := types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: "Hit root of webserver",
		},
	}
	common.JSONResponse(r, w, 200, JSONresp)
}

// GetVersion ...
// returns version information
func GetVersion(w http.ResponseWriter, r *http.Request) {
	// root of API
	version := common.GetAppBuildVersion()
	commitHash := common.GetAppBuildHash()
	date := common.GetAppBuildDate()
	mode := common.GetAppBuildMode()

	JSONresp := types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: "Responded with version information",
		},
		Spec: types.VersionInformation{
			Version:    version,
			CommitHash: commitHash,
			Date:       date,
			Mode:       mode,
		},
	}
	common.JSONResponse(r, w, 200, JSONresp)
}

// GetPins ...
// returns all pins and states
func GetPins(w http.ResponseWriter, r *http.Request) {
	responseMsg := "Failed to fetch pins"
	responseCode := 500
	pinsList, err := pin.ListPins()
	if err == nil {
		responseMsg = "Fetched all pins"
		responseCode = 200
	}
	JSONresp := types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: responseMsg,
		},
		Spec: pinsList,
	}
	common.JSONResponse(r, w, responseCode, JSONresp)
}

// GetPin ...
// returns a given pin and it's state
func GetPin(w http.ResponseWriter, r *http.Request) {
	responseMsg := "Failed to fetch pin state"
	responseCode := 500
	vars := mux.Vars(r)
	pinIdstr := vars["pin"]
	if pin.CheckGPIOpinForGround(r) == true {
		responseMsg = "Requested pin is a ground pin."
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: responseMsg,
			},
			Spec: nil,
		}
		common.JSONResponse(r, w, responseCode, JSONresp)
		return
	}
	pinID, err := strconv.Atoi(pinIdstr)
	pinState, err1 := pin.GetPin(pinID)
	if err == nil && err1 == nil {
		responseMsg = "Fetched pin state"
		responseCode = 200
	}
	JSONresp := types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: responseMsg,
		},
		Spec: pinState,
	}
	common.JSONResponse(r, w, responseCode, JSONresp)
}

// PostPin ...
// writes state to a given pin
func PostPin(w http.ResponseWriter, r *http.Request) {
	responseMsg := "Failed to update pin state"
	responseCode := 500
	vars := mux.Vars(r)
	pinIDstr := vars["pin"]
	if pin.CheckGPIOpinForGround(r) == true {
		responseMsg = "Requested pin is a ground pin."
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: responseMsg,
			},
			Spec: nil,
		}
		common.JSONResponse(r, w, responseCode, JSONresp)
		return
	}
	pinID, err := strconv.Atoi(pinIDstr)
	pinStateStr := vars["state"]
	if pin.CheckForValidState(r) == false {
		responseMsg = "Invalid pin state."
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: responseMsg,
			},
			Spec: nil,
		}
		common.JSONResponse(r, w, responseCode, JSONresp)
		return
	}
	pinState, err := strconv.Atoi(pinStateStr)
	pin, err := pin.WritePin(pinID, pinState, 1)
	if err == nil {
		responseMsg = "Updated pin state"
		responseCode = 200
	}
	JSONresp := types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: responseMsg,
		},
		Spec: pin,
	}
	common.JSONResponse(r, w, responseCode, JSONresp)
}

// GetGroundPinEndpoint ...
// response for if pin is a ground pin
func GetGroundPinEndpoint(w http.ResponseWriter, r *http.Request) {
	common.JSONResponse(r, w, 400, types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: "This pin is a ground pin and is not able to be toggled",
		},
	})
}

// GetUnknownEndpoint ...
// response for if endpoint doesn't exist
func GetUnknownEndpoint(w http.ResponseWriter, r *http.Request) {
	common.JSONResponse(r, w, 404, types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: "This endpoint doesn't seem to exist.",
		},
	})
}

// APIfailedAuth ...
// response if authorization doesn't work
func APIfailedAuth(w http.ResponseWriter, r *http.Request) {
	common.JSONResponse(r, w, 401, types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: "Unauthorized",
		},
	})
}

// HTTPvalidateAuth ...
// checks the auth for a request
func HTTPvalidateAuth(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if valid, err := common.CheckAuth(r, w); valid == true && err == nil {
			h.ServeHTTP(w, r)
		} else {
			APIfailedAuth(w, r)
		}
	}
}
