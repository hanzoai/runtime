// Copyright 2025 Hanzo Industries Inc.
// SPDX-License-Identifier: AGPL-3.0

package tools

import "github.com/hanzoai/runtime/cli/apiclient"

var runtimeMCPHeaders map[string]string = map[string]string{
	apiclient.RuntimeSourceHeader: "runtime-mcp",
}
