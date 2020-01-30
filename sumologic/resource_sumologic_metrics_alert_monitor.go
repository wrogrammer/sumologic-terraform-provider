package sumologic

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"time"
)

func resourceSumologicMetricsAlertMonitor() *schema.Resource {
	return &schema.Resource{
		Create: resourceSumologicMetricsAlertMonitorCreate,
		Read:   resourceSumologicMetricsAlertMonitorRead,
		Delete: resourceSumologicMetricsAlertMonitorDelete,
		Update: resourceSumologicMetricsAlertMonitorUpdate,
		Exists: resourceSumologicMetricsAlertMonitorExists,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 255),
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Default: 	  "",
				ValidateFunc: validation.StringLenBetween(0, 4095),
			},
			"alert_queries": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				MaxItems: 6,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"row_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"query": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringLenBetween(0, 4095),
						},
					},
				},
			},
			"timezone": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: ianaTimeZoneDatabaseFormat(),
			},
			"monitor_rules": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"warning_rule": getRuleSchema(),
						"critical_rule": getRuleSchema(),
						"missing_data_rule": getMissingDataRuleSchema(),
					},
				},
			},
		},
	}
}

func getRuleSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MinItems: 1,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"threshold_type": {
					Type:     	  schema.TypeString,
					Required: 	  true,
					ValidateFunc: validation.StringInSlice([]string{"Above", "Below"}, false),
				},
				"threshold": {
					Type:     schema.TypeFloat,
					Required: true,
				},
				"duration": {
					Type:     	  schema.TypeString,
					Required: 	  true,
					ValidateFunc: validation.StringInSlice([]string{"0", "5m", "10m", "15m", "30m", "60m"}, false),
				},
				"notifications": getNotificationsSchema(),
			},
		},
	}
}

func getMissingDataRuleSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MinItems: 1,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"affected_time_series": {
					Type:     	  schema.TypeString,
					Required: 	  true,
					ValidateFunc: validation.StringInSlice([]string{"all", "any"}, false),
				},
				"duration": {
					Type:     	  schema.TypeInt,
					Required: 	  true,
					ValidateFunc: validation.IntBetween(60000, 3600000),
				},
				"notifications": getNotificationsSchema(),
			},
		},
	}
}

func getNotificationsSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MinItems: 1,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"email_notifications": {
					Type:     schema.TypeList,
					Optional: true,
					MinItems: 1,
					MaxItems: 1,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"recipients": {
								Type: 	  schema.TypeList,
								Required: true,
								MinItems: 1,
								Elem: &schema.Schema{
									Type: 		  schema.TypeString,
								},
							},
						},
					},
				},
				"webhook_notifications": {
					Type:     schema.TypeList,
					Optional: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"webhook_id": {
								Type:     schema.TypeString,
								Required: true,
							},
							"payload": {
								Type:     schema.TypeString,
								Optional: true,
							},
						},
					},
				},
			},
		},
	}
}

func resourceSumologicMetricsAlertMonitorCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	monitor := resourceToMetricsAlertMonitor(d)

	id, err := c.CreateMetricsAlertMonitor(monitor)
	if err != nil {
		return err
	}

	d.SetId(id)

	return resourceSumologicMetricsAlertMonitorRead(d, meta)
}

func resourceSumologicMetricsAlertMonitorRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	monitor, err := c.GetMetricsAlertMonitor(d.Id())
	if err != nil {
		return err
	}

	if err = d.Set("name", monitor.Name); err != nil {
		return err
	}
	if err = d.Set("description", monitor.Description); err != nil {
		return err
	}
	if err = d.Set("alert_queries", flattenAlertQueries(monitor.AlertQueries)); err != nil {
		return err
	}
	if err = d.Set("timezone", monitor.Timezone); err != nil {
		return err
	}
	if err = d.Set("monitor_rules", flattenMonitorRules([1]MonitorRules{monitor.MonitorRules})); err != nil {
		return err
	}

	return nil
}

func resourceSumologicMetricsAlertMonitorDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	return c.DeleteMetricsAlertMonitor(d.Id())
}

func resourceSumologicMetricsAlertMonitorUpdate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	monitor := resourceToMetricsAlertMonitor(d)

	err := c.UpdateMetricsAlertMonitor(d.Id(), monitor)
	if err != nil {
		return err
	}

	return resourceSumologicMetricsAlertMonitorRead(d, meta)
}

func resourceSumologicMetricsAlertMonitorExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	c := meta.(*Client)

	monitor, err := c.GetMetricsAlertMonitor(d.Id())
	if err != nil {
		return false, err
	}

	return monitor != nil, nil
}

func resourceToMetricsAlertMonitor(d *schema.ResourceData) MetricsAlertMonitor {
	alertQueries := getAlertQueries(d.Get("alert_queries").([]interface{}))
	monitorRules := getMonitorRules(d.Get("monitor_rules").([]interface{})[0])

	metricsAlertMonitor := MetricsAlertMonitor{
		Name:		  d.Get("name").(string),
		Description:  d.Get("description").(string),
		AlertQueries: alertQueries,
		Timezone: 	  d.Get("timezone").(string),
		MonitorRules: monitorRules,
	}

	return metricsAlertMonitor
}

func getAlertQueries(alertQueries []interface{}) []AlertQuery {
	var result []AlertQuery
	for _, alertQueryRaw := range alertQueries {
		alertQueryMap := alertQueryRaw.(map[string]interface{})
		alertQuery := AlertQuery{
			RowId: alertQueryMap["row_id"].(string),
			Query: alertQueryMap["query"].(string),
		}
		result = append(result, alertQuery)
	}
	return result
}

func getMonitorRules(monitorRules interface{}) MonitorRules {
	if monitorRules == nil {
		return MonitorRules{}
	}
	var warningRule, criticalRule *Rule
	var missingDataRule *MissingDataRule

	var rulesMap = monitorRules.(map[string]interface{})
	warningRule = getRule("warning_rule", rulesMap)
	criticalRule = getRule("critical_rule", rulesMap)
	missingDataRule = getMissingDataRule(rulesMap)

	return MonitorRules{
		WarningRule: warningRule,
		CriticalRule: criticalRule,
		MissingDataRule: missingDataRule,
	}
}

func getRule(key string, rulesMap map[string]interface{}) *Rule {
	var rule *Rule
	if ruleRaw := rulesMap[key].([]interface{}); len(ruleRaw) == 1 {
		ruleMap := ruleRaw[0].(map[string]interface{})
		thresholdType := ruleMap["threshold_type"].(string)
		threshold := ruleMap["threshold"].(float64)
		duration := ruleMap["duration"].(string)
		notifications := getNotifications(ruleMap)
		rule = &Rule{
			ThresholdType:	thresholdType,
			Threshold:		threshold,
			Duration:		duration,
			Notifications:  notifications,
		}
	}
	return rule
}

func getMissingDataRule(rulesMap map[string]interface{}) *MissingDataRule {
	var missingDataRule *MissingDataRule
	if missingDataRuleRaw := rulesMap["missing_data_rule"].([]interface{}); len(missingDataRuleRaw) == 1 {
		missingDataRuleMap := missingDataRuleRaw[0].(map[string]interface{})
		affectedTimeSeries := missingDataRuleMap["affected_time_series"].(string)
		duration := missingDataRuleMap["duration"].(int)
		notifications := getNotifications(missingDataRuleMap)
		missingDataRule = &MissingDataRule{
			AffectedTimeSeries: affectedTimeSeries,
			Duration:           duration,
			Notifications:      notifications,
		}
	}
	return missingDataRule
}

func getNotifications(ruleMap map[string]interface{}) *Notifications {
	var notifications *Notifications
	if notificationsRaw := ruleMap["notifications"].([]interface{}); len(notificationsRaw) == 1 {
		if notificationsRaw[0] != nil {
			notificationsMap := notificationsRaw[0].(map[string]interface{})
			emailNotifications := getEmailNotifications(notificationsMap)
			webhookNotifications := getWebhookNotifications(notificationsMap)
			notifications = &Notifications{
				EmailNotifications:   emailNotifications,
				WebhookNotifications: webhookNotifications,
			}
		} else {
			notifications = &Notifications{}
		}
	}
	return notifications
}

func getEmailNotifications(notificationsMap map[string]interface{}) *EmailNotifications {
	var emailNotifications *EmailNotifications
	if emailNotificationsRaw := notificationsMap["email_notifications"].([]interface{}); len(emailNotificationsRaw) == 1 {
		if emailNotificationsRaw[0] != nil {
			emailNotificationsMap := emailNotificationsRaw[0].(map[string]interface{})
			recipientsRaw := emailNotificationsMap["recipients"].([]interface{})
			var recipients []string
			for _, recipientRaw := range recipientsRaw {
				recipients = append(recipients, recipientRaw.(string))
			}
			emailNotifications = &EmailNotifications{
				Recipients:	recipients,
			}
		} else {
			emailNotifications = &EmailNotifications{
				Recipients: []string{},
			}
		}
	}
	return emailNotifications
}

func getWebhookNotifications(notificationsMap map[string]interface{}) []WebhookNotification {
	var webhookNotifications []WebhookNotification
	if webhookNotificationsRawOrNil, ok := notificationsMap["webhook_notifications"]; ok {
		webhookNotificationsRaw := webhookNotificationsRawOrNil.([]interface{})
		for _, webhookNotificationRaw := range webhookNotificationsRaw {
			webhookNotificationMap := webhookNotificationRaw.(map[string]interface{})
			var payload string
			if payloadRaw, ok := webhookNotificationMap["payload"]; ok {
				payload = payloadRaw.(string)
			}
			webhookNotification := WebhookNotification{
				WebhookId: webhookNotificationMap["webhook_id"].(string),
				Payload:   payload,
			}
			webhookNotifications = append(webhookNotifications, webhookNotification)
		}
	}
	return webhookNotifications
}

func ianaTimeZoneDatabaseFormat() schema.SchemaValidateFunc {
	return func(i interface{}, k string) (s []string, es []error) {
		v, ok := i.(string)
		if !ok {
			es = append(es, fmt.Errorf("expected type of %s to be string", k))
			return
		}
		if v == "" || v == "UTC" || v == "Local" {
			es = append(es, fmt.Errorf("timezone must be explicitly named"))
		}
		_, err := time.LoadLocation(v)
		if err != nil {
			es = append(es, fmt.Errorf("%s is not a correct timezone in IANA Time Zone Database format", v))
		}
		return
	}
}

func flattenAlertQueries(in []AlertQuery) []map[string]interface{} {
	var out = make([]map[string]interface{}, len(in), len(in))
	for i, v := range in {
		m := make(map[string]interface{})
		m["row_id"] = v.RowId
		m["query"] = v.Query
		out[i] = m
	}
	return out
}

func flattenMonitorRules(in [1]MonitorRules) []interface{} {
	var out = make([]interface{}, 1, 1)
	m := make(map[string]interface{})
	monitorRules := in[0]
	if monitorRules.WarningRule != nil {
		m["warning_rule"] = flattenRule([1]Rule{*monitorRules.WarningRule})
	}
	if monitorRules.CriticalRule != nil {
		m["critical_rule"] = flattenRule([1]Rule{*monitorRules.CriticalRule})
	}
	if monitorRules.MissingDataRule != nil {
		m["missing_data_rule"] = flattenMissingDataRule([1]MissingDataRule{*monitorRules.MissingDataRule})
	}
	out[0] = m
	return out
}

func flattenRule(in [1]Rule) []interface{} {
	var out = make([]interface{}, 1, 1)
	m := make(map[string]interface{})
	rule := in[0]
	m["threshold_type"] = rule.ThresholdType
	m["threshold"] = rule.Threshold
	m["duration"] = rule.Duration
	if rule.Notifications != nil {
		m["notifications"] = flattenNotifications([1]Notifications{*rule.Notifications})
	}
	out[0] = m
	return out
}

func flattenMissingDataRule(in [1]MissingDataRule) []interface{} {
	var out = make([]interface{}, 1, 1)
	m := make(map[string]interface{})
	missingDataRule := in[0]
	m["affected_time_series"] = missingDataRule.AffectedTimeSeries
	m["duration"] = missingDataRule.Duration
	if missingDataRule.Notifications != nil {
		m["notifications"] = flattenNotifications([1]Notifications{*missingDataRule.Notifications})
	}
	out[0] = m
	return out
}

func flattenNotifications(in [1]Notifications) []interface{} {
	var out = make([]interface{}, 1, 1)
	m := make(map[string]interface{})
	notifications := in[0]
	if notifications.EmailNotifications != nil {
		m["email_notifications"] = flattenEmailNotifications([1]EmailNotifications{*notifications.EmailNotifications})
	}
	if notifications.WebhookNotifications != nil {
		m["webhook_notifications"] = flattenWebhookNotifications(notifications.WebhookNotifications)
	}
	out[0] = m
	return out
}

func flattenEmailNotifications(in [1]EmailNotifications) []interface{} {
	var out = make([]interface{}, 1, 1)
	m := make(map[string]interface{})
	emailNotifications := in[0]
	m["recipients"] = emailNotifications.Recipients
	out[0] = m
	return out
}

func flattenWebhookNotifications(in []WebhookNotification) []map[string]interface{} {
	var out = make([]map[string]interface{}, len(in), len(in))
	for i, v := range in {
		m := make(map[string]interface{})
		m["webhook_id"] = v.WebhookId
		if v.Payload != "" {
			m["payload"] = v.Payload
		}
		out[i] = m
	}
	return out
}
