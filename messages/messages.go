package messages

const (
	ClientPluginLoad               = "sagaClient plugin loaded"
	ClientServiceCall              = "Call \"%s\" transaction @%s"
	CallNumberOfBackendService     = "Number of services to call: %d"
	CallServiceGlobalTransactionID = "Universal Transaction ID: %s"

	ClientConfigFileError          = "Error reading \"client.json\" file: %s"
	ClientEndpointNotFoundError    = "No matching endpoint found as %s in \"config.json\""
	ClientConfigFIleUnmarshalError = "Error unmarshalling \"client.json\" file: %s"
	ClientConfigEndpointEmptyError = "Endpoint name can't be empty (%d)"
	ClientConfigMessagesEmptyError = "Endpoint messages can't be empty (%d)"
	ClientConfigAliasEmptyError    = "Endpoint alias can't be empty (%d, %d)"
	ClientConfigUrlEmptyError      = "Endpoint url can't be empty (%d, %d)"
	ClientConfigMethodError        = "Endpoint method can only be GET, POST, PUT, PATCH, and (%d, %d)"

	ClientRollbackError    = "Call \"%s\" rollback transaction @%s..."
	ClientCloseBodyError   = "Closing response body of \"%s\" encountered a problem.\nError: %s"
	ClientStatusCodeError  = "Call \"%s\" failed with status #%d. Rollback transaction will be called..."
	ClientServiceCallError = "Call \"%s\" has failed\nError: %s\nRollback transaction will be called..."
)
