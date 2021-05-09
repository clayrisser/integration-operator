package v1alpha2

type NamespacedName struct {
	// name
	Name string `json:"name"`

	// namespace
	Namespace string `json:"namespace,omitempty"`
}
