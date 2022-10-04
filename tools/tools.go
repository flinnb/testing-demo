//go:build tools
// +build tools

// See https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module

package tools

import (
	_ "github.com/stretchr/testify/assert"
	_ "github.com/stretchr/testify/require"
	_ "gotest.tools/gotestsum"
)
