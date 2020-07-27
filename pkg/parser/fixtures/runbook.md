# MyService

## Alerts

### Instance down

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

## Actions

### Restart Service

```bash
systemctl restart myservice
```

- Check if the service is running by SSHing into the instance and executing `systemctl status myservice`.
- Check the logs.
