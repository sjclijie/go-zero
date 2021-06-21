package metric

import (
	prom "github.com/prometheus/client_golang/prometheus"
	"github.com/sjclijie/go-zero/core/proc"
)

type (
	GaugeVecOpts VectorOpts

	GuageVec interface {
		Set(v float64, labels ...string)
		Inc(labels ...string)
		Add(v float64, labels ...string)
		close() bool
	}

	promGuageVec struct {
		gauge              *prom.GaugeVec
		defaultLabelValues []string
	}
)

func NewGaugeVec(cfg *GaugeVecOpts) GuageVec {
	if cfg == nil {
		return nil
	}

	cfg.Labels = append(cfg.Labels, CommonLabel...)

	vec := prom.NewGaugeVec(
		prom.GaugeOpts{
			Namespace: cfg.Namespace,
			Subsystem: cfg.Subsystem,
			Name:      cfg.Name,
			Help:      cfg.Help,
		}, cfg.Labels)
	prom.MustRegister(vec)

	defaultLabelValues := make([]string, len(CommonLabel))
	defaultLabelValues = append(defaultLabelValues, prom.BuildFQName(cfg.Namespace, cfg.Subsystem, cfg.Name))
	defaultLabelValues = append(defaultLabelValues, "gauge")
	defaultLabelValues = append(defaultLabelValues, cfg.AppId)
	defaultLabelValues = append(defaultLabelValues, cfg.Env)
	defaultLabelValues = append(defaultLabelValues, cfg.Ip)
	defaultLabelValues = append(defaultLabelValues, cfg.DataType)

	gv := &promGuageVec{
		gauge:              vec,
		defaultLabelValues: defaultLabelValues,
	}
	proc.AddShutdownListener(func() {
		gv.close()
	})

	return gv
}

func (gv *promGuageVec) Inc(labels ...string) {
	labels = append(labels, gv.defaultLabelValues...)
	gv.gauge.WithLabelValues(labels...).Inc()
}

func (gv *promGuageVec) Add(v float64, labels ...string) {
	labels = append(labels, gv.defaultLabelValues...)
	gv.gauge.WithLabelValues(labels...).Add(v)
}

func (gv *promGuageVec) Set(v float64, labels ...string) {
	labels = append(labels, gv.defaultLabelValues...)
	gv.gauge.WithLabelValues(labels...).Set(v)
}

func (gv *promGuageVec) close() bool {
	return prom.Unregister(gv.gauge)
}
