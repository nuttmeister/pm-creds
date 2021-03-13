// Package types includes the interfaces that must be satisfied to create an
// provider that can be used by the provider package.
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
