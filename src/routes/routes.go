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
	pinId, err := strconv.Atoi(pinIdstr)
	pinStateStr := vars["state"]
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
		Spec: pin.State,
	}
	common.JSONResponse(r, w, responseCode, JSONresp)
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

