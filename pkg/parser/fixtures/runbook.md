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

- Check if the service is running by SSHing into the instance and executing `systemctl status myservice`.
- Try restarting the service: `systemctl restart myservice`.
- Check the logs.
