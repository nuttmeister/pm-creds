// Package types includes the interfaces that must be satisfied to create a
// provider that can be used by the providers package.
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
