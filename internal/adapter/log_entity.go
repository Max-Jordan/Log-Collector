package adapter

import "time"

type LogEntity struct {
	Timestamp time.Time `json:"timestamp"`
	SessionID string    `json:"session_id"`
	Hook      string    `json:"hook"`
	Project   string    `json:"project"`
	User      string    `json:"user"`
	Branch    string    `json:"branch"`
	Commit    string    `json:"commit"`
	Status    string    `json:"status"`
	Message   string    `json:"message"`
	Details   Details   `json:"details"`
}

type Details struct {
	Total_findings int    `json:"total_findings"`
	Total_commits  int    `json:"total_commits"`
	Failed_commits string `json:"failed_commits"`
}
