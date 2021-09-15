package messages

const (
	CallService                    = "Transaction completed successfully"
	ClientPluginLoad               = "fidiboSagaClient plugin loaded"
	ClientServiceCall              = "Call \"%s\" transaction @%s"
	CallServiceRollback            = "Rollback completed successfully"
	CallNumberOfBackendService     = "Number of services to call: %d"
	CallServiceGlobalTransactionID = "Transaction Id of the call: %s"

	ClientConfigFileError          = "Error reading \"client.json\" file: %s"
	ClientEndpointNotFoundError    = "No matching endpoint found as %s in \"config.json\""
	ClientConfigFIleUnmarshalError = "Error unmarshalling \"client.json\" file: %s"

	ClientRollbackError    = "Call \"%s\" rollback transaction @%s..."
	ClientCloseBodyError   = "Closing response body of \"%s\" encountered a problem.\nError: %s"
	ClientStatusCodeError  = "Call \"%s\" Failed with status #%d. Rollback transaction will be called..."
	ClientServiceCallError = "Call \"%s\" has failed\nError: %s\nRollback transaction will be called..."

	CallServiceRollbackError = "Rollback process failed"
)
