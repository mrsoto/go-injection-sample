package router

import "fmt"

type Config struct {
	BaseUrl string
}

func (c Config) GetBaseUrl() string {
	return c.BaseUrl
}

func (c Config) Child(p string) Config {
	return Config{BaseUrl: fmt.Sprintf("%s/%s", c.BaseUrl, p)}
}
