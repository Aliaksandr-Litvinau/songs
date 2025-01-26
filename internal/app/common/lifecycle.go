package common

import "context"

// Lifecycle defines a common interface for managing the lifecycle of components
type Lifecycle interface {
	Start(ctx context.Context) error
	Shutdown(ctx context.Context) error
}
