package main

// TransmissionRPCServerName Variable for declare rpc host name
var TransmissionRPCServerName string

// TransmissionSessionID valid value for header X-Transmission-Session-Id
var TransmissionSessionID string

const (
	// TransmissionRPCUrl default url for requests to RPC api
	TransmissionRPCUrl = "http://%s:9091/transmission/rpc"
	// RefreshSessionCodeError status code which use for change sessionID
	RefreshSessionCodeError = 409
)
