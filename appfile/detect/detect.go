package detect

// App will detect the application type for the given directory.
func App(dir string, c *Config) (string, error) {
	for _, d := range c.Detectors {
		check, err := d.Detect(dir)
		if err != nil {
			return "", err
		}

		if check {
			return d.Type, nil
		}
	}

	return "", nil
}
