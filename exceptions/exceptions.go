package exceptions

const (
	ClientReadBodyError    = "Reading body encountered an error: %s"
	ClientStatusCodeError  = "Call \"%s\" Failed with status #%d.\nError: %s\nRollback transaction will be called..."
	ClientBackendCallError = "Call \"%s\" has failed\nError: %s\nRollback transaction will be called..."
	ClientRollbackError    = "Call \"%s\" rollback transaction..."
)
