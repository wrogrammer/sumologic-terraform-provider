---
layout: "sumologic"
page_title: "SumoLogic: sumologic_metrics_alert_monitor"
description: |-
  Provides a Sumologic Metrics Alert Monitor
---

# sumologic_metrics_alert_monitor
Provides a [Sumologic Metrics Alert Monitor][1].

## Example Usage
```hcl
resource "sumologic_metrics_alert_monitor" "high_cpu_monitor" {
  name        = "High CPU monitor"
  description = "Fires when CPU reaches high levels"
  timezone    = "America/Los_Angeles"
  alert_queries {
    row_id = "A"
    query  = "cpu_usage"
  }
  monitor_rules { 
    warning_rule {
      threshold_type = "Above"
      threshold = 80
      duration = "5m"
      notifications {
        email_notifications {
          recipients = ["user@domain.com"]
        }
      }
    }
    critical_rule {
      threshold_type = "Above"
      threshold = 90
      duration = "10m"
      notifications {
        email_notifications {
          recipients = ["user@domain.com"]
        }
        webhook_notifications {
          webhook_id = "0000000131"
          payload = "Critical level of CPU"
        }
      }
    }
    missing_data_rule {
      affected_time_series = "all"
      duration = "5m"
      notifications {
        email_notifications {
          recipients = ["user@domain.com", "other_user@domain.com"]
        }			
      }
    }
  }
}
```

## Argument reference
The following arguments are supported:
- `name` - (Required) Monitor name.
- `description` - (Optional) Monitor description.
- `alert_queries` - (Required) Monitor queries.
- `timezone` - (Required) Monitor time zone in 
[IANA Time Zone Database format](https://en.wikipedia.org/wiki/List_of_tz_database_time_zones#List). Date time ranges 
shown in emails and sent to webhooks are expressed in this time zone.
- `monitor_rules` - (Required) Monitor rules.

### Schemas:  
  
`alert_queries` entry schema:
  - `row_id` - (Required) Row identifier.
  - `query` - (Required) A monitor query.
  
`monitor_rules` object schema:
  - `warning_rule` - (Optional)
  - `critical_rule` - (Optional)
  - `missing_data_rule` - (Optional)
  
`warning_rule` and `critical_rule` objects schema:
  -  `threshold_type` - (Required) One of: `Above`, `Below`
  -  `threshold` - (Required) Threshold for the monitor, data points above or below this threshold are treated as
    outliers.
  -  `duration` - (Required) A period of time, in which the alert condition has to remain true before the 
    notification is triggered. If a single threshold violation is enough to trigger the alert, set this field to 0. 
    Currently, the only accepted values are `0`, `5m`, `10m`, `15m`, `30m` and `60m` (`m` suffix means minutes).
  -  `notifications` - (Optional) Monitor notifications.
  
`missing_data_rule` object schema:
  -  `affected_time_series` - (Required) Defines when an alert should be raised: either when all or any time series are 
  missing data. Accepted values for this field are: `all` and `any`.
  -  `duration` - (Required) A time window. Currently, the only accepted values are `5m`, `10m`, `15m`, `30m` and `60m` 
    (m suffix means minutes). The minimum value is equal to the query quantization.
  -  `notifications` - (Optional) Monitor notifications.
  
`notifications` object schema:
  -  `email_notifications` - (Optional) Monitor email notifications.
  -  `webhook_notifications` - (Optional) Monitor webhook notifications.
  
`email_notifications` object schema:
  -  `recipients`- (Required) List of notification recipients' email addresses.
  
`webhook_notifications` entry schema:
  -  `webhook_id` - (Required) Identifier of the webhook.
  -  `payload` - (Optional) Webhook's payload.
   
## Attributes reference
The following attributes are exported:
- `id` - The internal ID of the monitor.

[Back to Index][0]

[0]: ../README.md
[1]: https://help.sumologic.com/Metrics/Metric-Queries-and-Alerts/Metrics_Monitors_and_Alerts