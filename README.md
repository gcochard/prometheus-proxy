# Prometheus Custom Collector Proxy

The Custom Collector proxy is a library for writing proxy collectors for Prometheus. It's inspired by the snmp_exporter, and provides some helper functions for setting up a service to collect metrics from other systems not setup to export metrics.

## Example

```go
import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/nik-johnson-net/prometheus-proxy"
)

// factory is a function which takes a target and returns a Collector.
// NewCustomCollector is your constructor for your custom Collector.
func factory(target string) prometheus.Collector {
	return NewCustomCollector(target)
}

func main() {
	app := proxy.Application{
		CreateFactory: func() proxy.CollectorFactory {
			return factory
		},
	}

	proxy.Main(app)
}
```

In the prometheus.yml config file:
```yaml
scrape_configs:
  - job_name: 'prometheus-custom-collector'
    static_configs:
      - targets:
        - '192.168.1.40'
        - '192.168.1.41'
    relabel_configs:
      - source_labels: [__address__]
        target_label: __param_target
      - source_labels: [__param_target]
        target_label: instance
      - target_label: __address__
        replacement: 127.0.0.1:2112  # The prometheus-proxy's real hostname:port.
```

## License

This library is provided under the [MIT License](LICENSE.md)
