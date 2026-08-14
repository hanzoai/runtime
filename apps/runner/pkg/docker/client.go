// Copyright 2025 Hanzo Industries Inc.
// SPDX-License-Identifier: AGPL-3.0

package docker

import (
	"io"
	"sync"

	"github.com/hanzoai/runtime/apps/runner/pkg/cache"
	"github.com/docker/docker/client"
)

type DockerClientConfig struct {
	ApiClient             client.APIClient
	Cache                 cache.IRunnerCache
	LogWriter             io.Writer
	AWSRegion             string
	AWSEndpointUrl        string
	AWSAccessKeyId        string
	AWSSecretAccessKey    string
	DaemonPath            string
	ComputerUsePluginPath string
	// Runtimes the Docker daemon offers, which decides which boundaries this
	// node can provide.
	Runtimes map[string]bool
	// RuncOrgs are the organizations allowed to run without a boundary. Empty
	// names nobody.
	RuncOrgs map[string]bool
}

func NewDockerClient(config DockerClientConfig) *DockerClient {
	return &DockerClient{
		apiClient:             config.ApiClient,
		cache:                 config.Cache,
		logWriter:             config.LogWriter,
		awsRegion:             config.AWSRegion,
		awsEndpointUrl:        config.AWSEndpointUrl,
		awsAccessKeyId:        config.AWSAccessKeyId,
		awsSecretAccessKey:    config.AWSSecretAccessKey,
		volumeMutexes:         make(map[string]*sync.Mutex),
		daemonPath:            config.DaemonPath,
		computerUsePluginPath: config.ComputerUsePluginPath,
		runtimes:              config.Runtimes,
		runcOrgs:              config.RuncOrgs,
	}
}

type DockerClient struct {
	apiClient             client.APIClient
	cache                 cache.IRunnerCache
	logWriter             io.Writer
	awsRegion             string
	awsEndpointUrl        string
	awsAccessKeyId        string
	awsSecretAccessKey    string
	volumeMutexes         map[string]*sync.Mutex
	volumeMutexesMutex    sync.Mutex
	daemonPath            string
	computerUsePluginPath string
	runtimes              map[string]bool
	runcOrgs              map[string]bool
}
