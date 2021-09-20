package messages

import "encoding/json"

const (
	ClientPluginLoad               = "%s plugin has loaded"
	ClientServiceCall              = "Call \"%s\" transaction @ %s"
	CallNumberOfBackendService     = "Number of services to call: %d"
	CallServiceGlobalTransactionID = "Universal Transaction ID: %s"

	ClientConfigFileError          = "Error reading config file: %s"
	ClientEndpointNotFoundError    = "No matching endpoint found as %s in config file"
	ClientConfigFIleUnmarshalError = "Error unmarshalling config file: %s"

	ClientPluginLoadError          = "Config file has loaded with following errors. Resolve the errors and run again!"
	ClientConfigEndpointEmptyError = "\t- Endpoint name can't be empty (Endpoint: %d)"
	ClientConfigMessagesEmptyError = "\t- Endpoint messages can't be empty (Endpoint: %d)"
	ClientConfigAliasEmptyError    = "\t- Endpoint alias can't be empty (Endpoint: %d, Step: %d)"
	ClientConfigTimeoutError       = "\t- Endpoint timeout can't be equal or less than 0 (Endpoint: %d, Step: %d)"
	ClientConfigRetryError         = "\t- Endpoint max retry can't be less than 0 (Endpoint: %d, Step: %d)"
	ClientConfigRetryWaitError     = "\t- Endpoint max nor min wait retry can't be equal or less than 0 (Endpoint: %d, Step: %d)"
	ClientConfigStatusError        = "\t- Endpoint status can't be other than 2xx statuses (Endpoint: %d, Step: %d)"
	ClientConfigUrlEmptyError      = "\t- Endpoint url can't be empty (Endpoint: %d, Step: %d)"
	ClientConfigMethodError        = "\t- Endpoint method can only be GET, POST, PUT, PATCH, and (Endpoint: %d, Step: %d)"

	ClientRollbackError    = "Call \"%s\" rollback transaction @%s..."
	ClientCloseBodyError   = "Closing response body of \"%s\" encountered a problem.\nError: %s"
	ClientStatusCodeError  = "Call \"%s\" failed with status #%d. Rollback transaction will be called..."
	ClientServiceCallError = "Call \"%s\" has failed\nError: %s\nRollback transaction will be called..."

	ClientResponseWriterError = "Writing response encountered a problem: %s"
)

// Generate generates a message based on input
func Generate(m map[string]interface{}) (resp []byte) {
	resp, _ = json.Marshal(m)

	return
}
