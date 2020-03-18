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
			HandlerFunc:  common.HTTPuseMiddleware(APIgetPins, HTTPvalidateAuth),
			HttpMethod:   http.MethodGet,
		},
		{
			EndpointPath: endpointPrefix + "/pin/{pin:[0-9]+}",
			HandlerFunc:  common.HTTPuseMiddleware(APIgetPin, HTTPvalidateAuth),
			HttpMethod:   http.MethodGet,
		},
		{
			EndpointPath: endpointPrefix + "/pin/{pin:[0-9]+}/{state:[0-1]+}",
			HandlerFunc:  common.HTTPuseMiddleware(APIpostPin, HTTPvalidateAuth),
			HttpMethod:   http.MethodPost,
		},
	}
}
