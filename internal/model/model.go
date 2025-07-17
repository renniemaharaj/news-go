package model

// Singleton model instance
var singleton *Instance

func Initialize() {
	singleton = &Instance{}
	singleton.l = createLogger()
}

// Get returns the singleton instance
func Get() *Instance {
	if singleton == nil {
		Initialize()
	}

	return singleton
}
