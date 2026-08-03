// Copyright 2025 Daytona Platforms Inc.
// SPDX-License-Identifier: AGPL-3.0

package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/hanzoai/proxy/cmd/proxy/config"
	"github.com/hanzokv/go/v9"
)

type KVCache[T any] struct {
	kv        *kv.Client
	keyPrefix string
}

type ValueObject[T any] struct {
	Value T `json:"value"`
}

var client *kv.Client

func (c *KVCache[T]) Set(ctx context.Context, key string, value T, expiration time.Duration) error {
	jsonValue, err := json.Marshal(ValueObject[T]{Value: value})
	if err != nil {
		return err
	}
	return c.kv.Set(ctx, c.keyPrefix+key, string(jsonValue), expiration).Err()
}

func (c *KVCache[T]) Has(ctx context.Context, key string) (bool, error) {
	err := c.kv.Get(ctx, c.keyPrefix+key).Err()
	if err == nil {
		return true, nil
	}

	if err == kv.Nil {
		return false, nil
	}

	return false, err
}

func (c *KVCache[T]) Get(ctx context.Context, key string) (*T, error) {
	value, err := c.kv.Get(ctx, c.keyPrefix+key).Result()
	if err != nil {
		return nil, err
	}
	var result ValueObject[T]
	err = json.Unmarshal([]byte(value), &result)
	if err != nil {
		return nil, err
	}
	return &result.Value, nil
}

func (c *KVCache[T]) Delete(ctx context.Context, key string) error {
	return c.kv.Del(ctx, c.keyPrefix+key).Err()
}

func NewKVCache[T any](config *config.KVConfig, keyPrefix string) (*KVCache[T], error) {
	if config.Host == nil || config.Port == nil {
		return nil, errors.New("host and port are required")
	}

	password := ""
	if config.Password != nil {
		password = *config.Password
	}

	if client == nil {
		client = kv.NewClient(&kv.Options{
			Addr:     fmt.Sprintf("%s:%d", *config.Host, *config.Port),
			Password: password,
		})
	}

	return &KVCache[T]{
		kv:        client,
		keyPrefix: keyPrefix,
	}, nil
}
