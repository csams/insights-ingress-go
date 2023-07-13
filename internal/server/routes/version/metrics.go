package version

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

func ExposeVersion(v *IngressVersion) {
	versionMetric := promauto.NewGauge(prometheus.GaugeOpts{
		Namespace:   "ingress",
		Name:        "version",
		Help:        "Version information for Ingress service.",
		ConstLabels: prometheus.Labels{"commit": v.Commit, "version": v.Version},
	})
	versionMetric.Set(1.0)
}
