package persistence

// LogRecord represents an entry
type LogRecord struct {
	Namespace string
	Container string
	Pod       string
	Image     string
	Timestamp string
	Message   string
}
