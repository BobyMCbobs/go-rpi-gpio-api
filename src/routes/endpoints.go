package routes

import (
	"net/http"

	"gitlab.com/bobymcbobs/go-rpi-gpio-api/src/types"
	"gitlab.com/bobymcbobs/go-rpi-gpio-api/src/common"
)

func GetEndpoints(endpointPrefix string) types.Endpoints {
	return types.Endpoints{
		{
			EndpointPath: endpointPrefix + "/pin",
			HandlerFunc:  common.HTTPuseMiddleware(APIgetPins, HTTPvalidateAuth),
			HttpMethod:   http.MethodGet,
		},
		{
			EndpointPath: endpointPrefix + "/pin/{pin}",
			HandlerFunc:  common.HTTPuseMiddleware(APIgetPin, HTTPvalidateAuth),
			HttpMethod:   http.MethodGet,
		},
		{
			EndpointPath: endpointPrefix + "/pin/{pin}/{state}",
			HandlerFunc:  common.HTTPuseMiddleware(APIpostPin, HTTPvalidateAuth),
			HttpMethod:   http.MethodPost,
		},
	}
}
