package models

type HTTPResponse struct {
	Status     string                 `json:"status"`
	StatusCode int                    `json:"statusCode"`
	Message    string                 `json:"message"`
	Data       map[string]interface{} `json:"data"`
}
