package mcprotocol

const (
	IntentStatus = iota + 1
	IntentLogin
	IntentTransfer // added in 1.20.5 (24w03a)
)

// Deprecated: Next State has been renamed to Intent.
const (
	NextStateStatus   = IntentStatus
	NextStateLogin    = IntentLogin
	NextStateTransfer = IntentTransfer
)
