package methods

type CanReuestJson struct {
	Metrics metricsList `json:"metrics"`
}

type metricsList map[string]interface{}
