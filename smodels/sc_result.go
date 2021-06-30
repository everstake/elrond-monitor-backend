package swagger

type ScResult struct {
	From string `json:"from,omitempty"`

	To string `json:"to,omitempty"`

	Value float64 `json:"value,omitempty"`

	Data string `json:"data,omitempty"`

	Message string `json:"message,omitempty"`
}
