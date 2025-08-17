// Copyright 2025 Hanzo Industries Inc.
// SPDX-License-Identifier: AGPL-3.0

package port

type PortList struct {
	Ports []uint `json:"ports"`
} // @name PortList

type IsPortInUseResponse struct {
	IsInUse bool `json:"isInUse"`
} // @name IsPortInUseResponse
