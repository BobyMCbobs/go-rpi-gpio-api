package routes

import (
	"net/http"

	"gitlab.com/bobymcbobs/go-rpi-gpio-api/src/common"
	"gitlab.com/bobymcbobs/go-rpi-gpio-api/src/types"
)

func GetEndpoints(endpointPrefix string) types.Endpoints {
	return types.Endpoints{
		{
			EndpointPath: endpointPrefix + "/pin",
			HandlerFunc:  common.HTTPuseMiddleware(GetPins, HTTPvalidateAuth),
			HttpMethod:   http.MethodGet,
		},
		{
			EndpointPath: endpointPrefix + "/pin/{pin:[0-9]+}",
			HandlerFunc:  common.HTTPuseMiddleware(GetPin, HTTPvalidateAuth),
			HttpMethod:   http.MethodGet,
		},
		{
			EndpointPath: endpointPrefix + "/pin/{pin:[0-9]+}/{state:[0-1]+}",
			HandlerFunc:  common.HTTPuseMiddleware(PostPin, HTTPvalidateAuth),
			HttpMethod:   http.MethodPost,
		},
		{
			EndpointPath: endpointPrefix + "/version",
			HandlerFunc:  common.HTTPuseMiddleware(GetVersion, HTTPvalidateAuth),
			HttpMethod:   http.MethodGet,
		},
	}
}
