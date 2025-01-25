// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package taliexporter

import "go.opentelemetry.io/collector/component"

// Config is the struct for configuring the tali exporter
type Config struct {
	S3AccessKey string `mapstructure:"s3_access_key"`
	S3SecretKey string `mapstructure:"s3_secret_key"`
}

func createDefaultConfig() component.Config {
	return Config{}
}
