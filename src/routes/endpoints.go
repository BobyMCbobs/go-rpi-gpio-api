package routes

import (
	"net/http"

	"gitlab.com/bobymcbobs/go-rpi-gpio-api/src/common"
	"gitlab.com/bobymcbobs/go-rpi-gpio-api/src/types"
)

// GetEndpoints ...
// returns endpoints
func GetEndpoints(endpointPrefix string) types.Endpoints {
	return types.Endpoints{
		{
			EndpointPath: endpointPrefix + "/pin",
			HandlerFunc:  common.HTTPuseMiddleware(GetPins, HTTPvalidateAuth),
			HTTPMethod:   http.MethodGet,
		},
		{
			EndpointPath: endpointPrefix + "/pin/{pin:[0-9]+}",
			HandlerFunc:  common.HTTPuseMiddleware(GetPin, HTTPvalidateAuth),
			HTTPMethod:   http.MethodGet,
		},
		{
			EndpointPath: endpointPrefix + "/pin/{pin:[0-9]+}/{state:[0-1]+}",
			HandlerFunc:  common.HTTPuseMiddleware(PostPin, HTTPvalidateAuth),
			HTTPMethod:   http.MethodPost,
		},
		{
			EndpointPath: endpointPrefix + "/version",
			HandlerFunc:  common.HTTPuseMiddleware(GetVersion, HTTPvalidateAuth),
			HTTPMethod:   http.MethodGet,
		},
	}
}
