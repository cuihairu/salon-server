package config

const (
	DEV  = "dev"
	PROD = "prod"
	TEST = "test"
)

func (c *Config) IsDev() bool {
	return c.env == DEV
}

func (c *Config) IsTest() bool {
	return c.env == TEST
}

func (c *Config) IsProd() bool {
	return c.env == PROD
}

func (c *Config) GetEnv() string {
	return c.env
}

func (c *Config) setEnv(env string) {
	switch env {
	case DEV:
	case TEST:
		env = TEST
	case PROD:
		env = PROD
	default:
		env = DEV
	}
	c.env = env
}
