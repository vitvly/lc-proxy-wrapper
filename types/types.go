package types

const (
	ProxyInitialised int = iota
	OptimisticHeader
	FinalizedHeader
	Stopped
	Error
)

type ProxyEvent struct {
	EventType int
	Msg       string
}
