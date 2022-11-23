// SPDX-FileCopyrightText: 2020-2022 Open Networking Foundation <info@opennetworking.org>
// SPDX-FileCopyrightText: 2022-present Sander Tolsma. All rights reserved
// SPDX-License-Identifier: Apache-2.0

package logging

// SinkType is the type of a sink
type SinkType string

func (t SinkType) String() string {
	return string(t)
}

const (
	// StdoutSinkType is the sink type for stdout
	StdoutSinkType SinkType = "stdout"
	// StderrSinkType is the sink type for stderr
	StderrSinkType SinkType = "stderr"
	// FileSinkType is the type for a file sink
	FileSinkType SinkType = "file"
	// KafkaSinkType is the sink type for the Kafka sink
	KafkaSinkType SinkType = "kafka"
)

// SinkEncoding is the encoding for a sink
type SinkEncoding string

func (e SinkEncoding) String() string {
	return string(e)
}

const (
	// ConsoleEncoding is an encoding for outputs to the console
	ConsoleEncoding SinkEncoding = "console"
	// JSONEncoding is an encoding for JSON outputs
	JSONEncoding SinkEncoding = "json"
)

const (
	rootLoggerName = "root"
)

// Config logging configuration
type Config struct {
	Loggers map[string]LoggerConfig `json:"loggers" yaml:"loggers"`
	Sinks   map[string]SinkConfig   `json:"sinks" yaml:"sinks"`
}

// GetRootLogger returns the root logger configuration
func (c Config) GetRootLogger() LoggerConfig {
	root := c.Loggers[rootLoggerName]
	if root.Output == nil {
		defaultSink := ""
		root.Output = map[string]OutputConfig{
			"": {
				Sink: &defaultSink,
			},
		}
	}
	return root
}

// GetLoggers returns the configured loggers
func (c Config) GetLoggers() []LoggerConfig {
	loggers := make([]LoggerConfig, 0, len(c.Loggers))
	for name, logger := range c.Loggers {
		if name != rootLoggerName {
			logger.Name = name
			if logger.Output == nil {
				logger.Output = map[string]OutputConfig{}
			}
			loggers = append(loggers, logger)
		}
	}
	return loggers
}

// GetLogger returns a logger by name
func (c Config) GetLogger(name string) (LoggerConfig, bool) {
	if name == rootLoggerName {
		return LoggerConfig{}, false
	}

	logger, ok := c.Loggers[name]
	if ok {
		logger.Name = name
		if logger.Output == nil {
			logger.Output = map[string]OutputConfig{}
		}
		return logger, true
	}
	return LoggerConfig{Output: map[string]OutputConfig{}}, false
}

// GetSinks returns the configured sinks
func (c Config) GetSinks() []SinkConfig {
	sinks := []SinkConfig{
		{
			Name: "",
		},
	}
	for name, sink := range c.Sinks {
		sink.Name = name
		sinks = append(sinks, sink)
	}
	return sinks
}

// GetSink returns a sink by name
func (c Config) GetSink(name string) (SinkConfig, bool) {
	if name == "" {
		return SinkConfig{
			Name: "",
		}, true
	}
	sink, ok := c.Sinks[name]
	if ok {
		sink.Name = name
		return sink, true
	}
	return SinkConfig{}, false
}

// LoggerConfig is the configuration for a logger
type LoggerConfig struct {
	Name   string                  `json:"name" yaml:"name"`
	Level  *string                 `json:"level" yaml:"level,omitempty"`
	Output map[string]OutputConfig `json:"output" yaml:"output"`
}

// GetLevel returns the logger level
func (c LoggerConfig) GetLevel() Level {
	level := c.Level
	if level != nil {
		return levelStringToLevel(*level)
	}
	return ErrorLevel
}

// GetOutputs returns the logger outputs
func (c LoggerConfig) GetOutputs() []OutputConfig {
	outputs := c.Output
	outputsList := make([]OutputConfig, 0, len(outputs))
	for name, output := range outputs {
		output.Name = name
		outputsList = append(outputsList, output)
	}
	return outputsList
}

// GetOutput returns an output by name
func (c LoggerConfig) GetOutput(name string) (OutputConfig, bool) {
	output, ok := c.Output[name]
	if ok {
		output.Name = name
		return output, true
	}
	return OutputConfig{}, false
}

// OutputConfig is the configuration for a sink output
type OutputConfig struct {
	Name  string  `json:"name" yaml:"name"`
	Sink  *string `json:"sink" yaml:"sink"`
	Level *string `json:"level" yaml:"level,omitempty"`
}

// GetSink returns the output sink
func (c OutputConfig) GetSink() string {
	sink := c.Sink
	if sink != nil {
		return *sink
	}
	return ""
}

// GetLevel returns the output level
func (c OutputConfig) GetLevel() Level {
	level := c.Level
	if level != nil {
		return levelStringToLevel(*level)
	}
	return DebugLevel
}

// SinkConfig is the configuration for a sink
type SinkConfig struct {
	Name     string            `json:"name" yaml:"name"`
	Type     *SinkType         `json:"type" yaml:"type,omitempty"`
	Encoding *SinkEncoding     `json:"encoding" yaml:"encoding,omitempty"`
	Stdout   *StdoutSinkConfig `json:"stdout" yaml:"stdout,omitempty"`
	Stderr   *StderrSinkConfig `json:"stderr" yaml:"stderr,omitempty"`
	File     *FileSinkConfig   `json:"file" yaml:"file,omitempty"`
	Kafka    *KafkaSinkConfig  `json:"kafka" yaml:"kafka,omitempty"`
}

// GetType returns the sink type
func (c SinkConfig) GetType() SinkType {
	sinkType := c.Type
	if sinkType != nil {
		return *sinkType
	}
	return StdoutSinkType
}

// GetEncoding returns the sink encoding
func (c SinkConfig) GetEncoding() SinkEncoding {
	encoding := c.Encoding
	if encoding != nil {
		return *encoding
	}
	return ConsoleEncoding
}

// GetStdoutSinkConfig returns the stdout sink configuration
func (c SinkConfig) GetStdoutSinkConfig() StdoutSinkConfig {
	config := c.Stdout
	if config != nil {
		return *config
	}
	return StdoutSinkConfig{}
}

// GetStderrSinkConfig returns the stderr sink configuration
func (c SinkConfig) GetStderrSinkConfig() StderrSinkConfig {
	config := c.Stderr
	if config != nil {
		return *config
	}
	return StderrSinkConfig{}
}

// GetFileSinkConfig returns the file sink configuration
func (c SinkConfig) GetFileSinkConfig() FileSinkConfig {
	config := c.File
	if config != nil {
		return *config
	}
	return FileSinkConfig{}
}

// GetKafkaSinkConfig returns the Kafka sink configuration
func (c SinkConfig) GetKafkaSinkConfig() KafkaSinkConfig {
	config := c.Kafka
	if config != nil {
		return *config
	}
	return KafkaSinkConfig{}
}

// StdoutSinkConfig is the configuration for an stdout sink
type StdoutSinkConfig struct {
}

// StderrSinkConfig is the configuration for an stderr sink
type StderrSinkConfig struct {
}

// FileSinkConfig is the configuration for a file sink
type FileSinkConfig struct {
	Path string `json:"path" yaml:"path"`
}

// KafkaSinkConfig is the configuration for a Kafka sink
type KafkaSinkConfig struct {
	Topic   string   `json:"topic" yaml:"topic"`
	Key     string   `json:"key" yaml:"key"`
	Brokers []string `json:"brokers" yaml:"brokers"`
}
