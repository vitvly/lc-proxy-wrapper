package types

const (
	ProxyInitialised int = iota
	OptimisticHeader
	FinalizedHeader
)

type ProxyEvent struct {
	EventType int
	Msg       string
}
