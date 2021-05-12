package v1alpha2

type Phase string

const (
	PendingPhase   Phase = "Pending"
	SucceededPhase       = "Succeeded"
	FailedPhase          = "Failed"
	UnknownPhase         = "Unknown"
)

type NamespacedName struct {
	// name
	Name string `json:"name"`

	// namespace
	Namespace string `json:"namespace,omitempty"`
}
