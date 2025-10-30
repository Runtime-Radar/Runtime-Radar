package metrics

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	client "github.com/prometheus/client_model/go"
	"gorm.io/gorm"
	gorm_prometheus "gorm.io/plugin/prometheus"
)

var _ prometheus.Gatherer = (*Registry)(nil)

type Registry struct {
	prometheus.Registerer

	gatherer prometheus.Gatherer
}

func (w *Registry) Gather() ([]*client.MetricFamily, error) {
	return w.gatherer.Gather()
}

// NewRegistry creates and registers service global metrics
// Custom collectors could be added via extraCollectors
func NewRegistry(service, cluster string, extraCollectors ...prometheus.Collector) (*Registry, error) {
	registry := prometheus.NewRegistry()
	registerer := prometheus.WrapRegistererWith(map[string]string{
		"service": service,
		"cluster": cluster,
	}, registry)

	// Default system collectors
	if err := registerer.Register(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{})); err != nil {
		return nil, fmt.Errorf("can't register process metrics: %w", err)
	}

	if err := registerer.Register(collectors.NewGoCollector()); err != nil {
		return nil, fmt.Errorf("can't register go metrics: %w", err)
	}

	for _, collector := range extraCollectors {
		if err := registerer.Register(collector); err != nil {
			return nil, fmt.Errorf("can't register extra metrics: %w", err)
		}
	}

	return &Registry{Registerer: registerer, gatherer: registry}, nil
}

func GormPGMetrics(dbName string, db *gorm.DB) ([]prometheus.Collector, error) {
	pm := gorm_prometheus.New(gorm_prometheus.Config{
		DBName: db.Name() + "_" + dbName,
		MetricsCollector: []gorm_prometheus.MetricsCollector{
			&gorm_prometheus.Postgres{
				VariableNames: []string{"Threads_running"},
			},
		},
	})

	return pm.Collectors, db.Use(pm)
}
