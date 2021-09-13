package messages

const (
	ClientPluginLoad  = "fidiboSagaClient plugin loaded"
	ClientServiceCall = "Call \"%s\" transaction @%s"

	ClientConfigFileError          = "Error reading \"client.json\" file: %s"
	ClientConfigFIleUnmarshalError = "Error unmarshalling \"client.json\" file: %s"
	ClientEndpointNotFoundError    = "No matching endpoint found as %s in \"config.json\""

	ClientStatusCodeError  = "Call \"%s\" Failed with status #%d. Rollback transaction will be called..."
	ClientServiceCallError = "Call \"%s\" has failed\nError: %s\nRollback transaction will be called..."
	ClientCloseBodyError   = "Closing response body of \"%s\" encountered a problem.\nError: %s"
	ClientRollbackError    = "Call \"%s\" rollback transaction @ %s..."
)
