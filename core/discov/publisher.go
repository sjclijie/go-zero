package discov

type Publisher interface {
	Register() error
	Deregister() error
}
