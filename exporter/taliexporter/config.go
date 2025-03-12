// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package taliexporter

import (
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
)

// Config is the struct for configuring the tali exporter
type Config struct {
	// These entries are only for headless mode.
	S3AccessKey string `mapstructure:"s3_access_key"`
	S3SecretKey string `mapstructure:"s3_secret_key"`
	S3Endpoint  string `mapstructure:"s3_endpoint"`
	S3UseSSL    bool   `mapstructure:"s3_use_ssl"`
	// Timeout controls the timeout for sending the data to Tali+S3 and dictates
	// the max amount of wait time.
	Timeout time.Duration `mapstructure:"timeout"`

	// Set this for testing (data is /dev/null'ed)
	S3DevNullMode bool `mapstructure:"s3_dev_null_mode"`

	// Enable new string-interned index marshalling
	SIIMarshal bool `mapstructure:"sii_marshal"`

	QueueSettings exporterhelper.QueueConfig `mapstructure:"sending_queue"`
}

func createDefaultConfig() component.Config {
	return Config{
		QueueSettings: exporterhelper.NewDefaultQueueConfig(),
		Timeout:       20 * time.Second,
	}
}
