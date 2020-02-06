package sumologic

import (
	"encoding/json"
	"fmt"
)

func (s *Client) GetMetricsAlertMonitor(id string) (*MetricsAlertMonitor, error) {
	responseBody, _, err := s.Get(fmt.Sprintf("v1/metricsAlertMonitors/%s", id))
	if err != nil {
		return nil, err
	}

	var monitorInfo MetricsAlertMonitorInfo
	err = json.Unmarshal(responseBody, &monitorInfo)
	if err != nil {
		return nil, err
	}

	return &(monitorInfo.MonitorDefinition), nil
}

func (s *Client) CreateMetricsAlertMonitor(monitor MetricsAlertMonitor) (string, error) {
	responseBody, err := s.Post("v1/metricsAlertMonitors", monitor)
	if err != nil {
		return "", err
	}

	var createdMonitorInfo MetricsAlertMonitorInfo
	err = json.Unmarshal(responseBody, &createdMonitorInfo)
	if err != nil {
		return "", err
	}

	return createdMonitorInfo.Id, nil
}

func (s *Client) DeleteMetricsAlertMonitor(id string) error {
	_, err := s.Delete(fmt.Sprintf("v1/metricsAlertMonitors/%s", id))

	return err
}

func (s *Client) UpdateMetricsAlertMonitor(id string, monitor MetricsAlertMonitor) error {
	_, err := s.Put(fmt.Sprintf("v1/metricsAlertMonitors/%s", id), monitor)

	return err
}

type MetricsAlertMonitorInfo struct {
	MonitorDefinition MetricsAlertMonitor `json:"monitorDefinition"`
	Id                string              `json:"id"`
}

type MetricsAlertMonitor struct {
	Name         string       `json:"name"`
	Description  string       `json:"description,omitempty"`
	AlertQueries []AlertQuery `json:"alertQueries"`
	Timezone     string       `json:"timezone"`
	MonitorRules MonitorRules `json:"monitorRules"`
}

type AlertQuery struct {
	RowId string `json:"rowId"`
	Query string `json:"query"`
}

type MonitorRules struct {
	WarningRule     *Rule            `json:"warningRule,omitempty"`
	CriticalRule    *Rule            `json:"criticalRule,omitempty"`
	MissingDataRule *MissingDataRule `json:"missingDataRule,omitempty"`
}

type Rule struct {
	ThresholdType string         `json:"thresholdType"`
	Threshold     float64        `json:"threshold"`
	Duration      string         `json:"duration"`
	Notifications *Notifications `json:"notifications,omitempty"`
}

type MissingDataRule struct {
	AffectedTimeSeries string         `json:"affectedTimeSeries"`
	Duration           int            `json:"duration"`
	Notifications      *Notifications `json:"notifications,omitempty"`
}

type Notifications struct {
	EmailNotifications   *EmailNotifications   `json:"emailNotifications,omitempty"`
	WebhookNotifications []WebhookNotification `json:"webhookNotifications,omitempty"`
}

type EmailNotifications struct {
	Recipients []string `json:"recipients"`
}

type WebhookNotification struct {
	WebhookId string `json:"webhookId"`
	Payload   string `json:"payload,omitempty"`
}
