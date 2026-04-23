package vault

type Provider interface {
	Resolve(path string) (string, error)
}
