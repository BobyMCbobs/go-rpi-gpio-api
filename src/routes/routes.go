/*
	route related
*/

package routes

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gitlab.com/bobymcbobs/go-rpi-gpio-api/src/common"
	"gitlab.com/bobymcbobs/go-rpi-gpio-api/src/types"
)

func GetRoot(w http.ResponseWriter, r *http.Request) {
	// root of API
	JSONresp := types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: "Hit root of webserver",
		},
	}
	common.JSONResponse(r, w, 200, JSONresp)
}

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
			Version: version,
			CommitHash: commitHash,
			Date: date,
			Mode: mode,
		},
	}
	common.JSONResponse(r, w, 200, JSONresp)
}

func GetPins(w http.ResponseWriter, r *http.Request) {
	responseMsg := "Failed to fetch pins"
	responseCode := 500
	pinsList, err := common.ListPins()
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

func GetPin(w http.ResponseWriter, r *http.Request) {
	responseMsg := "Failed to fetch pin state"
	responseCode := 500
	vars := mux.Vars(r)
	pinIdstr := vars["pin"]
	if common.CheckGPIOpinForGround(r) == true {
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
	pinId, err := strconv.Atoi(pinIdstr)
	pinState, err1 := common.GetPin(pinId)
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

func PostPin(w http.ResponseWriter, r *http.Request) {
	responseMsg := "Failed to update pin state"
	responseCode := 500
	vars := mux.Vars(r)
	pinIdstr := vars["pin"]
	if common.CheckGPIOpinForGround(r) == true {
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
	pinId, err := strconv.Atoi(pinIdstr)
	pinStateStr := vars["state"]
	if common.CheckForValidState(r) == false {
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
	pin, err := common.WritePin(pinId, pinState, 1)
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

func GetGroundPinEndpoint(w http.ResponseWriter, r *http.Request) {
	common.JSONResponse(r, w, 400, types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: "This pin is a ground pin and is not able to be toggled",
		},
	})
}

func GetUnknownEndpoint(w http.ResponseWriter, r *http.Request) {
	common.JSONResponse(r, w, 404, types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: "This endpoint doesn't seem to exist.",
		},
	})
}

func APIfailedAuth(w http.ResponseWriter, r *http.Request) {
	common.JSONResponse(r, w, 401, types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: "Unauthorized",
		},
	})
}

func HTTPvalidateAuth(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if valid, err := common.CheckAuth(r, w); valid == true && err == nil {
			h.ServeHTTP(w, r)
		} else {
			APIfailedAuth(w, r)
		}
	}
}
