package metric

import (
	prom "github.com/prometheus/client_golang/prometheus"
	"github.com/sjclijie/go-zero/core/proc"
)

type (
	CounterVecOpts VectorOpts

	CounterVec interface {
		Inc(lables ...string)
		Add(v float64, labels ...string)
		close() bool
	}

	promCounterVec struct {
		counter            *prom.CounterVec
		defaultLabelValues []string
	}
)

func NewCounterVec(cfg *CounterVecOpts) CounterVec {
	if cfg == nil {
		return nil
	}

	cfg.Labels = append(cfg.Labels, CommonLabel...)

	vec := prom.NewCounterVec(prom.CounterOpts{
		Namespace: cfg.Namespace,
		Subsystem: cfg.Subsystem,
		Name:      cfg.Name,
		Help:      cfg.Help,
	}, cfg.Labels)
	prom.MustRegister(vec)

	defaultLabelValues := make([]string, 0)
	defaultLabelValues = append(defaultLabelValues, prom.BuildFQName(cfg.Namespace, cfg.Subsystem, cfg.Name))
	defaultLabelValues = append(defaultLabelValues, "counter")
	defaultLabelValues = append(defaultLabelValues, cfg.AppId)
	defaultLabelValues = append(defaultLabelValues, cfg.Env)
	defaultLabelValues = append(defaultLabelValues, cfg.Ip)
	defaultLabelValues = append(defaultLabelValues, cfg.DataType)

	cv := &promCounterVec{
		counter:            vec,
		defaultLabelValues: defaultLabelValues,
	}

	proc.AddShutdownListener(func() {
		cv.close()
	})

	return cv
}

func (cv *promCounterVec) Inc(labels ...string) {
	labels = append(labels, cv.defaultLabelValues...)
	cv.counter.WithLabelValues(labels...).Inc()
}

func (cv *promCounterVec) Add(v float64, labels ...string) {
	labels = append(labels, cv.defaultLabelValues...)
	cv.counter.WithLabelValues(labels...).Add(v)
}

func (cv *promCounterVec) close() bool {
	return prom.Unregister(cv.counter)
}
