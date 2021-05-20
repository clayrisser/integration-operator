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
	SocketCoupledStatusCondition     StatusCondition = "SocketCoupled"
	SocketCreatedStatusCondition     StatusCondition = "SocketCreated"
	SocketEmptyStatusCondition       StatusCondition = "SocketEmpty"
	SocketNotCreatedStatusCondition  StatusCondition = "SocketNotCreated"
	SocketNotReadyStatusCondition    StatusCondition = "SocketNotReady"
)
