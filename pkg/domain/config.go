package domain

type Config struct {
	BaseURL string
	Headers map[string]string
}

func NewConfig(baseURL string) *Config {
	return &Config{
		BaseURL: baseURL,
		Headers: make(map[string]string),
	}
}

func (c *Config) WithHeader(key, value string) *Config {
	c.Headers[key] = value
	return c
}
