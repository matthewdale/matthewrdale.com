package metrics

import (
	"net/http"

	gometrics "github.com/rcrowley/go-metrics"
)

type Service struct {
	registry gometrics.Registry
}

func New(registry gometrics.Registry) *Service {
	return &Service{
		registry: registry,
	}
}

type MetricsArgs struct{}

type MetricsReply struct {
	Metrics map[string]interface{}
}

func (svc *Service) Metrics(r *http.Request, args *MetricsArgs, reply *MetricsReply) error {
	metrics := make(map[string]interface{})
	svc.registry.Each(func(name string, metric interface{}) {
		switch m := metric.(type) {
		case gometrics.Counter:
			metrics[name] = m.Count()
		case gometrics.Gauge:
			metrics[name] = m.Value()
		case gometrics.GaugeFloat64:
			metrics[name] = m.Value()
		case gometrics.Histogram:
			metrics[name+".mean"] = m.Mean()
			metrics[name+".max"] = m.Max()
			metrics[name+".50percentile"] = m.Percentile(0.50)
			metrics[name+".99percentile"] = m.Percentile(0.99)
			metrics[name+".count"] = m.Count()
		case gometrics.Meter:
			metrics[name+".rate5"] = m.Rate5()
			metrics[name+"rate15"] = m.Rate15()
			metrics[name+".count"] = m.Count()
		case gometrics.Timer:
			metrics[name+".mean"] = m.Mean()
			metrics[name+".max"] = m.Max()
			metrics[name+".50percentile"] = m.Percentile(0.50)
			metrics[name+".99percentile"] = m.Percentile(0.99)
			metrics[name+".count"] = m.Count()
		}
	})
	reply.Metrics = metrics
	return nil
}
