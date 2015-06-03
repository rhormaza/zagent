package zutil


// Version information
const (
	ZAGT_VERSION = "zagent-v1.0.0"
	ZAGT_MAJOR   = 1 
	ZAGT_MINOR   = 0
	ZAGT_BUILD   = 0

    DEFAULT_LOG = 4
)

const (
	FINEST int = iota // 0
	FINE                // 1
	DEBUG               // 2
	TRACE               // 3
	INFO                // 4
	WARNING             // 5
	ERROR               // 6
	CRITICAL            // 7
)

//// Logging level strings
//var (
//	levelStrings = [...]string{"FNST", "FINE", "DEBG", "TRAC", "INFO", "WARN", "EROR", "CRIT"}
//)
