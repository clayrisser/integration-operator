package v1alpha2

type Phase string

const (
	PendingPhase   Phase = "Pending"
	SucceededPhase Phase = "Succeeded"
	FailedPhase    Phase = "Failed"
	UnknownPhase   Phase = "Unknown"
)

type NamespacedName struct {
	// name
	Name string `json:"name"`

	// namespace
	Namespace string `json:"namespace,omitempty"`
}
