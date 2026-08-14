// Copyright 2025 Hanzo Industries Inc.
// SPDX-License-Identifier: AGPL-3.0

package common

import (
	metric "github.com/luxfi/metric"
)

type PrometheusOperationStatus string

const (
	PrometheusOperationStatusSuccess PrometheusOperationStatus = "success"
	PrometheusOperationStatusFailure PrometheusOperationStatus = "failure"
)

// Define your metrics
var (
	// Histogram to track duration of container operations, by the boundary the
	// sandbox runs behind, which is what makes one boundary comparable to
	// another.
	//
	// The lowest bucket used to be 0.1s, so every creation that met the
	// sub-90ms claim landed in the same bucket as one that missed it and the
	// claim could not be observed at all. There is a boundary at 0.09 now: the
	// share of creations meeting it is one query.
	ContainerOperationDuration = metric.NewHistogramVec(
		metric.HistogramOpts{
			Name:    "container_operation_duration_seconds",
			Help:    "Time taken for container operations in seconds",
			Buckets: []float64{0.01, 0.025, 0.05, 0.075, 0.09, 0.1, 0.25, 0.5, 1, 2, 5, 10, 30, 60, 120, 300},
		},
		[]string{"operation", "isolation"},
	)

	// Counter to track occurrence of container operations with status
	ContainerOperationCount = metric.NewCounterVec(
		metric.CounterOpts{
			Name: "container_operation_total",
			Help: "Total number of container operations",
		},
		[]string{"operation", "status"},
	)
)
