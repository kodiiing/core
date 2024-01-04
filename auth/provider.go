package auth

type Provider uint8

const (
	ProviderGithub Provider = iota
	ProviderGitlab
)

func (p Provider) ToUint8() uint8 {
	return uint8(p)
}
