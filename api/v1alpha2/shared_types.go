package v1alpha2

type Phase string

const (
	FailedPhase    Phase = "Failed"
	PendingPhase   Phase = "Pending"
	ReadyPhase     Phase = "Ready"
	SucceededPhase Phase = "Succeeded"
	UnknownPhase   Phase = "Unknown"
)

type When string

const (
	BrokenWhen    When = "broken"
	CoupledWhen   When = "coupled"
	CreatedWhen   When = "created"
	DecoupledWhen When = "decoupled"
	DeletedWhen   When = "deleted"
	UpdatedWhen   When = "updated"
)

type Do string

const (
	ApplyDo  Do = "apply"
	DeleteDo Do = "delete"
)

type Resource struct {
	Do       Do     `json:"do,omitempty"`
	Resource string `json:"resource,omitempty"`
	When     When   `json:"when,omitempty"`
}

type NamespacedName struct {
	// name
	Name string `json:"name"`

	// namespace
	Namespace string `json:"namespace,omitempty"`
}
