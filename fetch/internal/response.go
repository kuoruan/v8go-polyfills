package internal

type Response struct {
	Headers    map[string]string `json:"headers"`
	Status     int               `json:"status"`
	StatusText string            `json:"statusText"`
	OK         bool              `json:"ok"`
	Redirected bool              `json:"redirected"`
	URL        string            `json:"url"`
	Body       string            `json:"body"`
}
