package config

// Configer is implemented by configurations.
type Configer interface {

	// IsValid returns an error if the current configuration is invalid for any reason.
	// This function is typically called after configuration settings have been loaded into a Configer.
	IsValid() error
}
