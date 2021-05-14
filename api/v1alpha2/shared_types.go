package v1alpha2

type Phase string

const (
	FailedPhase    Phase = "Failed"
	PendingPhase   Phase = "Pending"
	ReadyPhase     Phase = "Ready"
	SucceededPhase Phase = "Succeeded"
	UnknownPhase   Phase = "Unknown"
)

type NamespacedName struct {
	// name
	Name string `json:"name"`

	// namespace
	Namespace string `json:"namespace,omitempty"`
}
