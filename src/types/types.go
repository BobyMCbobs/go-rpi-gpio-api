/*
	handle all types used by API
*/

package types

import (
	"net/http"
)

type JSONResponseMetadata struct {
	URL       string `json:"selfLink"`
	Version   string `json:"version"`
	Timestamp int64  `json:"timestamp"`
	Response  string `json:"response"`
}

type JSONMessageResponse struct {
	Metadata JSONResponseMetadata `json:"metadata"`
	Spec     interface{}          `json:"spec"`
}

type Endpoints []struct {
	EndpointPath string
	HandlerFunc  http.HandlerFunc
	HttpMethod   string
}

type Pin struct {
	Number int `json:"number"`
	State int `json:"state"`
}

type PinList []Pin

