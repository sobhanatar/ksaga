package exceptions

const (
	ClientReadBodyError    = "Reading body encountered an error: %s"
	ClientStatusCodeError  = "Call to \"%s\" Failed with status #%d.\nError: %s\nRollback transaction will be called..."
	ClientBackendCallError = "Call to \"%s\" has failed\nError: %s\nRollback transaction will be called..."
)
