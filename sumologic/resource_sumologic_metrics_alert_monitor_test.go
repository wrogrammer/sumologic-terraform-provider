package sumologic

import (
	"reflect"
	"testing"
)

func TestExpandAlertQueries(t *testing.T) {
	alertQueriesConfig := getTestAlertQueriesConfig()
	expanded := expandAlertQueries(alertQueriesConfig)
	expected := getTestAlertQueries()
	if !reflect.DeepEqual(expanded, expected) {
		t.Fatalf("Got:\n\n%#v\n\nExpected:\n\n%#v", expanded, expected)
	}
}

func TestFlattenAlertQueries(t *testing.T) {
	alertQueries := getTestAlertQueries()
	flattened := flattenAlertQueries(alertQueries)
	expected := getTestAlertQueriesConfig()
	if !reflect.DeepEqual(flattened, expected) {
		t.Fatalf("Got:\n\n%#v\n\nExpected:\n\n%#v", flattened, expected)
	}
}

func TestExpandMonitorRules(t *testing.T) {
	monitorRulesConfig := getTestMonitorRulesConfig()
	expanded := expandMonitorRules(monitorRulesConfig)
	expected := getTestMonitorRules()
	if !reflect.DeepEqual(expanded, expected) {
		t.Fatalf("Got:\n\n%#v\n\nExpected:\n\n%#v", expanded, expected)
	}
}

func TestFlattenMonitorRules(t *testing.T) {
	monitorRules := getTestMonitorRules()
	flattened := flattenMonitorRules([1]MonitorRules{monitorRules})[0]
	expected := getTestMonitorRulesConfig()
	if !reflect.DeepEqual(flattened, expected) {
		t.Fatalf("Got:\n\n%#v\n\nExpected:\n\n%#v", flattened, expected)
	}
}

func getTestAlertQueries() []AlertQuery {
	return []AlertQuery{
		{RowId: "A", Query: "cpu_usage"},
		{RowId: "B", Query: "#A | avg"},
	}
}

func getTestAlertQueriesConfig() []interface{} {
	return []interface{}{
		map[string]interface{}{
			"row_id": "A",
			"query":  "cpu_usage",
		},
		map[string]interface{}{
			"row_id": "B",
			"query":  "#A | avg",
		},
	}
}

func getTestMonitorRules() MonitorRules {
	warningRuleEmailNotifications := &EmailNotifications{
		Recipients: []string{"user@domain.com"},
	}
	warningRuleNotifications := &Notifications{
		EmailNotifications: warningRuleEmailNotifications,
	}
	warningRule := &Rule{
		ThresholdType: "Above",
		Threshold:     80,
		Duration:      "5m",
		Notifications: warningRuleNotifications,
	}

	criticalRuleEmailNotifications := &EmailNotifications{
		Recipients: []string{"user@domain.com"},
	}
	criticalRuleWebhookNotifications := []WebhookNotification{
		{
			WebhookId: "0000000131", Payload: "Critical level of CPU",
		},
	}
	criticalRuleNotifications := &Notifications{
		EmailNotifications:   criticalRuleEmailNotifications,
		WebhookNotifications: criticalRuleWebhookNotifications,
	}
	criticalRule := &Rule{
		ThresholdType: "Above",
		Threshold:     90,
		Duration:      "10m",
		Notifications: criticalRuleNotifications,
	}

	missingDataRuleEmailNotifications := &EmailNotifications{
		Recipients: []string{"user@domain.com", "other_user@domain.com"},
	}
	missingDataRuleNotifications := &Notifications{
		EmailNotifications: missingDataRuleEmailNotifications,
	}
	missingDataRule := &MissingDataRule{
		AffectedTimeSeries: "all",
		Duration:           "5m",
		Notifications:      missingDataRuleNotifications,
	}

	return MonitorRules{
		WarningRule:     warningRule,
		CriticalRule:    criticalRule,
		MissingDataRule: missingDataRule,
	}
}

func getTestMonitorRulesConfig() map[string]interface{} {
	warningRuleEmailNotificationsConfig := []interface{}{
		map[string]interface{}{
			"recipients": []interface{}{"user@domain.com"},
		},
	}
	warningRuleNotificationsConfig := []interface{}{
		map[string]interface{}{
			"email_notifications": warningRuleEmailNotificationsConfig,
		},
	}
	warningRuleConfig := []interface{}{
		map[string]interface{}{
			"threshold_type": "Above",
			"threshold":      float64(80),
			"duration":       "5m",
			"notifications":  warningRuleNotificationsConfig,
		},
	}

	criticalRuleEmailNotificationsConfig := []interface{}{
		map[string]interface{}{
			"recipients": []interface{}{"user@domain.com"},
		},
	}
	criticalRuleWebhookNotificationsConfig := []interface{}{
		map[string]interface{}{
			"webhook_id": "0000000131",
			"payload":    "Critical level of CPU",
		},
	}
	criticalRuleNotificationsConfig := []interface{}{
		map[string]interface{}{
			"email_notifications":   criticalRuleEmailNotificationsConfig,
			"webhook_notifications": criticalRuleWebhookNotificationsConfig,
		},
	}
	criticalRuleConfig := []interface{}{
		map[string]interface{}{
			"threshold_type": "Above",
			"threshold":      float64(90),
			"duration":       "10m",
			"notifications":  criticalRuleNotificationsConfig,
		},
	}

	missingDataRuleEmailNotificationsConfig := []interface{}{
		map[string]interface{}{
			"recipients": []interface{}{"user@domain.com", "other_user@domain.com"},
		},
	}
	missingDataRuleNotificationsConfig := []interface{}{
		map[string]interface{}{
			"email_notifications": missingDataRuleEmailNotificationsConfig,
		},
	}
	missingDataRuleConfig := []interface{}{
		map[string]interface{}{
			"affected_time_series": "all",
			"duration":             "5m",
			"notifications":        missingDataRuleNotificationsConfig,
		},
	}

	monitorRulesConfig := map[string]interface{}{
		"warning_rule":      warningRuleConfig,
		"critical_rule":     criticalRuleConfig,
		"missing_data_rule": missingDataRuleConfig,
	}
	return monitorRulesConfig
}
