# runbooks

A Proof of concept to keep alerts and their runbooks next to each other.

Alerts are part of a runbook which is written in Markdown. Each Alert is a [Prometheus alerting rule](https://prometheus.io/docs/prometheus/latest/configuration/alerting_rules/).

## How it works

````markdown
# MyService <-- The h1 heading at the start is required. It will become the name of the runbook.

## Instance down <-- This h2 heading starts a new rule. Each h2 in the document that is followed by a yaml code block is treated as a rule.

```yaml <-- This is the Alerting rule. It will be parsed and exposed for consumption by a Prometheus server.
alert: InstanceDown
expr: up{job="MyService"} == 0
for: 1m
labels:
  severity: page
annotations:
  summary: Instance not responding
  description: "Instance: {{$labels.instance}}"
```

### Actions <-- This heading and the following list will not be parsed but still displayed in the UI. You can put what ever you want here.

- Check if the service is running by SSHing into the instance and execute `systemctl status myservice`.
- Try restarting the service: `systemctl restart myservice`.
- Check the logs at `/var/log/myservice/out.log`.
````

See also [examples/runbook.md](examples/runbook.md).

## Try it

```
# Build the Docker image
docker build -t runbooks .
# Start a Docker container
docker run -d -p "8090:8090" runbooks
# Create a new runbook at the API
curl -i -XPOST --data-binary @examples/runbook.md http://localhost:8090/api/v0/runbooks
# Visit http://localhost:8090/runbooks to browse available runbooks
# Download alerting rules ready for consumption by Prometheus
curl http://localhost:8090/api/v0/prometheus/rules
```
