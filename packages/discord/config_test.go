package discord

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDefaultConfig(t *testing.T) {
	t.Parallel()
	assert.NotNil(t, DefaultConfig())
}

func TestNewConfig(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name   string
		opts   []Option
		expect *Config
	}{
		{
			name:   "should return default config without options",
			opts:   []Option{},
			expect: DefaultConfig(),
		},
		{
			name: "should return config with API timeout set",
			opts: []Option{
				WithAPITimeout(10 * time.Second),
			},
			expect: &Config{
				APITimeout: 10 * time.Second,
			},
		},
		{
			name: "should return config with API timeout set to minimum",
			opts: []Option{
				WithAPITimeout(1),
			},
			expect: &Config{
				APITimeout: minAPITimeout,
			},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.expect, NewConfig(tc.opts...))
		})
	}
}
