# MyService

## Alerts

### Instance Down

```yaml
alert: InstanceDown
expr: up{job="MyService"} == 0
for: 1m
labels:
  severity: page
annotations:
  summary: Instance not responding
  description: "Instance: {{$labels.instance}}"
```

### High HTTP Errors

```yaml
alert: HighHTTPErrors
expr: sum(rate(http_requests_total{code=~"5.*",job="MyService"}[5m])) > 0
for: 2m
labels:
  severity: page
annotations:
  summary: Service returns 5xx status codes
  description: "RPS: {{$value}}"
```

## Actions

### Restart Service

```bash
ssh 192.168.1.1 systemctl restart myservice
```

- Check the logs at `/var/log/myservice/out.log`.
