// auth provides authentication procedure
package auth

type Provider uint8

const (
	ProviderGithub Provider = iota
	ProviderGitlab
)
