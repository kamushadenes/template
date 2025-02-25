package bootstrap

type BootstrapType uint64

const (
	BootstrapTypeEmbedded BootstrapType = iota
	BootstrapTypeExternalDir
	BootstrapTypeGitHub
)
