package http

import (
	"fmt"
)

const (
	ContentTypeProblemJSON = "application/problem+json"
	ProblemTypeBaseURL     = "https://weiss.example.com/problems"
)

type ProblemDetails struct {
	Type       string                 `json:"type"`
	Title      string                 `json:"title"`
	Status     int                    `json:"status"`
	Detail     string                 `json:"detail,omitempty"`
	Instance   string                 `json:"instance,omitempty"`
	Extensions map[string]interface{} `json:",omitempty"`
}

func NewProblemDetails(status int, title string, detail string, instance string) *ProblemDetails {
	problemType := fmt.Sprintf("%s/%d", ProblemTypeBaseURL, status)

	return &ProblemDetails{
		Type:       problemType,
		Title:      title,
		Status:     status,
		Detail:     detail,
		Instance:   instance,
		Extensions: make(map[string]interface{}),
	}
}

func (p *ProblemDetails) WithExtension(key string, value interface{}) *ProblemDetails {
	if p.Extensions == nil {
		p.Extensions = make(map[string]interface{})
	}
	p.Extensions[key] = value
	return p
}

func GetProblemTypeForStatus(status int) string {
	return fmt.Sprintf("%s/%d", ProblemTypeBaseURL, status)
}

func GetTitleForStatus(status int) string {
	switch status {
	case 400:
		return "Bad Request"
	case 401:
		return "Unauthorized"
	case 403:
		return "Forbidden"
	case 404:
		return "Not Found"
	case 409:
		return "Conflict"
	case 422:
		return "Unprocessable Entity"
	case 429:
		return "Too Many Requests"
	case 500:
		return "Internal Server Error"
	case 502:
		return "Bad Gateway"
	case 503:
		return "Service Unavailable"
	case 504:
		return "Gateway Timeout"
	default:
		if status >= 400 && status < 500 {
			return "Client Error"
		}
		if status >= 500 && status < 600 {
			return "Server Error"
		}
		return "Error"
	}
}

func GetInstanceFromPath(path string) string {
	return path
}
