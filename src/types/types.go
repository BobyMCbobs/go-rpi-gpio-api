/*
	handle all types used by API
*/

package types

import (
	"net/http"
)

// JSONResponseMetadata ...
// generic metadata for each request
type JSONResponseMetadata struct {
	URL       string `json:"selfLink"`
	Version   string `json:"version"`
	Timestamp int64  `json:"timestamp"`
	Response  string `json:"response"`
	Hostname  string `json:"hostname"`
}

// JSONMessageResponse ...
// generic JSON response
type JSONMessageResponse struct {
	Metadata JSONResponseMetadata `json:"metadata"`
	Spec     interface{}          `json:"spec"`
}

// Endpoints ...
// http route endpoints
type Endpoints []struct {
	EndpointPath string
	HandlerFunc  http.HandlerFunc
	HTTPMethod   string
}

// Pin ...
// pin information
type Pin struct {
	Number int `json:"number"`
	State  int `json:"state"`
}

// PinList ...
// many pins
type PinList []Pin

// VersionInformation ...
// information about an instance
type VersionInformation struct {
	Version    string `json:"version"`
	CommitHash string `json:"commitHash"`
	Mode       string `json:"mode"`
	Date       string `json:"date"`
}
