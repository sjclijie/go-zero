package metric

var (
	CommonLabel = []string{"hll_metric_name", "hll_metric_type", "hll_appid", "hll_env", "hll_ip", "hll_data_type"}
)

// VectorOpts general configuration
type VectorOpts struct {
	Namespace string
	Subsystem string
	Name      string
	Help      string
	Labels    []string
	AppId     string
	Env       string
	Ip        string
	DataType  string
}
