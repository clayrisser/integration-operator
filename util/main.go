package util

type Condition struct {
	Message string
	Reason  string
	Status  bool
	Type    string
}

type StatusCondition string

const (
	CouplingInProcessStatusCondition StatusCondition = "CouplingInProcess"
	CouplingSucceededStatusCondition StatusCondition = "CouplingSucceeded"
	ErrorStatusCondition             StatusCondition = "Error"
	PlugCreatedStatusCondition       StatusCondition = "PlugCreated"
	SocketCreatedStatusCondition     StatusCondition = "SocketCreated"
	SocketNotCreatedStatusCondition  StatusCondition = "SocketNotCreated"
	SocketNotReadyStatusCondition    StatusCondition = "SocketNotReady"
	SocketReadyStatusCondition       StatusCondition = "SocketReady"
)
