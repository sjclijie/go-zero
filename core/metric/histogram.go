package metric

import (
	prom "github.com/prometheus/client_golang/prometheus"
	"github.com/sjclijie/go-zero/core/proc"
)

type (
	HistogramVecOpts struct {
		Namespace string
		Subsystem string
		Name      string
		Help      string
		Labels    []string
		Buckets   []float64
		AppId     string
		Env       string
		Ip        string
		DataType  string
	}

	HistogramVec interface {
		Observe(v float64, labels ...string)
		close() bool
	}

	promHistogramVec struct {
		histogram          *prom.HistogramVec
		defaultLabelValues []string
	}
)

func NewHistogramVec(cfg *HistogramVecOpts) HistogramVec {
	if cfg == nil {
		return nil
	}

	cfg.Labels = append(cfg.Labels, CommonLabel...)

	vec := prom.NewHistogramVec(prom.HistogramOpts{
		Namespace: cfg.Namespace,
		Subsystem: cfg.Subsystem,
		Name:      cfg.Name,
		Help:      cfg.Help,
		Buckets:   cfg.Buckets,
	}, cfg.Labels)
	prom.MustRegister(vec)

	defaultLabelValues := make([]string, 0)
	defaultLabelValues = append(defaultLabelValues, prom.BuildFQName(cfg.Namespace, cfg.Subsystem, cfg.Name))
	defaultLabelValues = append(defaultLabelValues, "histogram")
	defaultLabelValues = append(defaultLabelValues, cfg.AppId)
	defaultLabelValues = append(defaultLabelValues, cfg.Env)
	defaultLabelValues = append(defaultLabelValues, cfg.Ip)
	defaultLabelValues = append(defaultLabelValues, cfg.DataType)

	hv := &promHistogramVec{
		histogram:          vec,
		defaultLabelValues: defaultLabelValues,
	}
	proc.AddShutdownListener(func() {
		hv.close()
	})

	return hv
}

func (hv *promHistogramVec) Observe(v float64, labels ...string) {
	labels = append(labels, hv.defaultLabelValues...)
	hv.histogram.WithLabelValues(labels...).Observe(v)
}

func (hv *promHistogramVec) close() bool {
	return prom.Unregister(hv.histogram)
}
