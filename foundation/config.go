package foundation

// Config is the configuration of the various foundational domains
// that Otto supports.
type Config struct {
	// Service settings are used to configure the service discovery
	// (if there is any).
	//
	// ServiceName is the name of the service. If this is blank,
	// service discovery won't be configured for this application.
	//
	// ServicePort is the port this service is running on.
	//
	// ServiceTags is a list of tags associated with the service.
	ServiceName string
	ServicePort int
	ServiceTags []string
}
