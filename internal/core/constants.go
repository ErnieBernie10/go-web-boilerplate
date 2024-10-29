package core

type Environment string

const (
	Development Environment = "local"
	Staging     Environment = "staging"
	Production  Environment = "production"
)

type ContextKey string

const UserContextKey ContextKey = "User"
const TokenContextKey ContextKey = "Token"
const RefreshContextKey ContextKey = "Refresh"