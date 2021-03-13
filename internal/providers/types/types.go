package types

type Provider interface {
	Name() string
	Type() string
	Get(name string) (Profile, error)
}

type Profile interface {
	Name() string
	Payload() []byte
}
