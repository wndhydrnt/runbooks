# runbooks

A Proof of concept to keep alerts and their runbooks next to each other.

Alerts are part of a runbook which is written in Markdown. Each Alert is a [Prometheus alerting rule](https://prometheus.io/docs/prometheus/latest/configuration/alerting_rules/).

## How it works

````markdown
# MyService <-- The h1 heading at the start is required. It will become the name of the runbook.

## Alerts <-- This h2 heading starts the alerts section. Within this section, each h3 heading followed by a code block are parsed as an alert.

### Instance down <-- This h3 heading becaomes the name of the alert.

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

## Actions <-- This h2 heading starts the actions section. Within this section, each h3 heading followed by a code block are parsed as an action. Each action can then be executed via the UI.

### Restart Service <-- This h3 heading becaomes the name of the action.

- Check if the service is running by SSHing into the instance and execute `systemctl status myservice`.
- Check the logs at `/var/log/myservice/out.log`.

```bash <-- This is the action.
systemctl restart myservice
```
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
