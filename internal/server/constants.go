package server

type Environment string

const (
	Development Environment = "local"
	Staging     Environment = "staging"
	Production  Environment = "production"
)
