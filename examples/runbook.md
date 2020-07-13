# MyService

## Instance down

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

### Actions

- Check if the service is running by SSHing into the instance and execute `systemctl status myservice`.
- Try restarting the service: `systemctl restart myservice`.
- Check the logs at `/var/log/myservice/out.log`.

## High HTTP Errors

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

### Actions

- Check the logs at `/var/log/myservice/out.log`.
- Check database connection.
