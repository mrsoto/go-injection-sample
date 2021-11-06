package router

import (
	"testing"

	"github.com/go-playground/assert/v2"
)

func Test_configUrl(t *testing.T) {

	cfg := Config{BaseUrl: "foo"}

	actual := cfg.GetBaseUrl()

	assert.Equal(t, actual, "foo")
}

func Test_configChield(t *testing.T) {

	cfg := Config{BaseUrl: "foo"}

	actual := cfg.Child("bar")

	expected := Config{BaseUrl: "foo/bar"}
	assert.Equal(t, expected, actual)
}
