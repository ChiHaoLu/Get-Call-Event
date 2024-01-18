package tracer

type TraceConfig struct {
	Tracer string `json:"tracer"`
}
type CallTracerResult struct {
	From    string `json:"from"`
	Gas     string `json:"gas"`
	GasUsed string `json:"gasUsed"`
	To      string `json:"to"`
	Input   string `json:"input"`
	Output  string `json:"output"`
	Value   string `json:"value"`
	Type    string `json:"type"`
}

type EventTracerResult struct {
	Data    []any `json:"data"`
}