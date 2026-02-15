package validator

// Exit codes as defined in DESIGN.md
const (
	ExitCodeSuccess            = 0
	ExitCodeTestFailed         = 1
	ExitCodePreconditionFailed = 2
	ExitCodeTeardownFailed     = 3
	ExitCodeSetupFailed        = 4
)
