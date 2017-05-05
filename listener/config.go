package listener

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/rewardStyle/kinetic/config"
	"github.com/rewardStyle/kinetic/logging"
)

// Config is used to configure a Listener instance
type Config struct {
	*config.AwsOptions
	*listenerOptions
	LogLevel aws.LogLevelType
}

// NewConfig creates a new instance of Config
func NewConfig(stream, shard string) *Config {
	return &Config{
		AwsOptions: config.DefaultAwsOptions(),
		listenerOptions: &listenerOptions{
			stream:                stream,
			shard:                 shard,
			batchSize:             10000,
			concurrency:           10000,
			shardIterator:         NewShardIterator(),
			getRecordsReadTimeout: 1 * time.Second,
			Stats: &NilStatsCollector{},
		},
		LogLevel: logging.LogOff,
	}
}

// SetAwsConfig configures the AWS Config used to create Sessions (and therefore
// kinesis clients).
func (c *Config) SetAwsConfig(config *aws.Config) {
	c.AwsConfig = config
}

// SetBatchSize configures the batch size of the GetRecords call.
func (c *Config) SetBatchSize(batchSize int) {
	c.batchSize = batchSize
}

// SetConcurrency controls the number of goroutines the Listener will spawn to
// process messages.
func (c *Config) SetConcurrency(concurrency int) {
	c.concurrency = concurrency
}

// SetInitialShardIterator configures the settings used to retrieve initial
// shard iterator via the GetShardIterator call.
func (c *Config) SetInitialShardIterator(shardIterator *ShardIterator) {
	c.shardIterator = shardIterator
}

// SetGetRecordsReadTimeout configures the time to wait for each successive
// Read operation on the GetRecords response payload.
func (c *Config) SetGetRecordsReadTimeout(timouet time.Duration) {
	c.getRecordsReadTimeout = timouet
}

// SetLogLevel configures both the SDK and Kinetic log levels.
func (c *Config) SetLogLevel(logLevel aws.LogLevelType) {
	c.AwsOptions.SetLogLevel(logLevel)
	c.LogLevel = logLevel & 0xffff0000
}

// SetStatsCollector configures a listener to handle listener metrics.
func (c *Config) SetStatsCollector(stats StatsCollector) {
	c.Stats = stats
}
