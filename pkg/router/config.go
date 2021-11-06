package router

import "fmt"

type Config struct {
	BaseUrl string
}

func (c Config) GetBaseUrl() string {
	return c.BaseUrl
}

func (c Config) Child(p string) Config {
	c.BaseUrl = fmt.Sprintf("%s/%s", c.BaseUrl, p)
	return c
}
