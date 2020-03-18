/*
	route related
*/

package routes

import (
	"net/http"
	"regexp"
	"strconv"

	"github.com/gorilla/mux"
	"gitlab.com/bobymcbobs/go-rpi-gpio-api/src/common"
	"gitlab.com/bobymcbobs/go-rpi-gpio-api/src/types"
)

func APIroot(w http.ResponseWriter, r *http.Request) {
	// root of API
	JSONresp := types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: "Hit root of webserver",
		},
	}
	common.JSONResponse(r, w, 200, JSONresp)
}

func APIgetPins(w http.ResponseWriter, r *http.Request) {
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

func APIgetPin(w http.ResponseWriter, r *http.Request) {
	responseMsg := "Failed to fetch pin state"
	responseCode := 500
	vars := mux.Vars(r)
	pinIdstr := vars["pin"]
	if CheckGPIOpinForGround(r) == true {
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

func APIpostPin(w http.ResponseWriter, r *http.Request) {
	responseMsg := "Failed to update pin state"
	responseCode := 500
	vars := mux.Vars(r)
	pinIdstr := vars["pin"]
	if CheckGPIOpinForGround(r) == true {
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
	if CheckForValidState(r) == false {
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

func APIgroundPinEndpoint(w http.ResponseWriter, r *http.Request) {
	common.JSONResponse(r, w, 400, types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: "This pin is a ground pin and is not able to be toggled",
		},
	})
}

func APIUnknownEndpoint(w http.ResponseWriter, r *http.Request) {
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
