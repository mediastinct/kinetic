package producer

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/kinesis"
	"github.com/aws/aws-sdk-go/service/kinesis/kinesisiface"

	"github.com/rewardStyle/kinetic/errs"
	"github.com/rewardStyle/kinetic/logging"
	"github.com/rewardStyle/kinetic/message"
)

type kinesisWriterOptions struct {
	Stats  StatsCollector
}

// KinesisWriter handles the API to send records to Kinesis.
type KinesisWriter struct {
	*kinesisWriterOptions
	*logging.LogHelper

	stream string
	client  kinesisiface.KinesisAPI
}

// NewKinesisWriter creates a new stream writer to write records to a Kinesis.
func NewKinesisWriter(c *aws.Config, stream string, fn ...func(*KinesisWriterConfig)) (*KinesisWriter, error) {
	cfg := NewKinesisWriterConfig(c)
	for _, f := range fn {
		f(cfg)
	}
	sess, err := cfg.GetSession()
	if err != nil {
		return nil, err
	}
	return &KinesisWriter{
		kinesisWriterOptions: cfg.kinesisWriterOptions,
		LogHelper: &logging.LogHelper{
			LogLevel: cfg.LogLevel,
			Logger:  cfg.AwsConfig.Logger,
		},
		stream: stream,
		client: kinesis.New(sess),
	}, nil
}

// PutRecords sends a batch of records to Kinesis and returns a list of records that need to be retried.
func (w *KinesisWriter) PutRecords(messages []*message.Message) ([]*message.Message, error) {
	var startSendTime time.Time
	var startBuildTime time.Time

	start := time.Now()
	var records []*kinesis.PutRecordsRequestEntry
	for _, msg := range messages {
		records = append(records, msg.MakeRequestEntry())
	}
	req, resp := w.client.PutRecordsRequest(&kinesis.PutRecordsInput{
		StreamName: aws.String(w.stream),
		Records:    records,
	})

	req.Handlers.Build.PushFront(func(r *request.Request) {
		startBuildTime = time.Now()
		w.LogDebug("Start PutRecords Build, took", time.Since(start))
	})

	req.Handlers.Build.PushBack(func(r *request.Request) {
		w.Stats.AddPutRecordsBuildDuration(time.Since(startBuildTime))
		w.LogDebug("Finished PutRecords Build, took", time.Since(start))
	})

	req.Handlers.Send.PushFront(func(r *request.Request) {
		startSendTime = time.Now()
		w.LogDebug("Start PutRecords Send took", time.Since(start))
	})

	req.Handlers.Send.PushBack(func(r *request.Request) {
		w.Stats.AddPutRecordsSendDuration(time.Since(startSendTime))
		w.LogDebug("Finished PutRecords Send, took", time.Since(start))
	})

	w.LogDebug("Starting PutRecords Build/Sign request, took", time.Since(start))
	w.Stats.AddPutRecordsCalled(1)
	if err := req.Send(); err != nil {
		w.LogError("Error putting records:", err.Error())
		return nil, err
	}
	w.Stats.AddPutRecordsDuration(time.Since(start))

	if resp == nil {
		return nil, errs.ErrNilPutRecordsResponse
	}
	if resp.FailedRecordCount == nil {
		return nil, errs.ErrNilFailedRecordCount
	}
	attempted := len(messages)
	failed := int(aws.Int64Value(resp.FailedRecordCount))
	sent := attempted - failed
	w.LogDebug(fmt.Sprintf("Finished PutRecords request, %d records attempted, %d records successful, %d records failed, took %v\n", attempted, sent, failed, time.Since(start)))

	var retries []*message.Message
	var err error
	for idx, record := range resp.Records {
		if record.SequenceNumber != nil && record.ShardId != nil {
			// TODO: per-shard metrics
			messages[idx].SequenceNumber = record.SequenceNumber
			messages[idx].ShardID = record.ShardId
		} else {
			switch aws.StringValue(record.ErrorCode) {
			case kinesis.ErrCodeProvisionedThroughputExceededException:
				w.Stats.AddProvisionedThroughputExceeded(1)
			default:
				w.LogDebug("PutRecords record failed with error:", aws.StringValue(record.ErrorCode), aws.StringValue(record.ErrorMessage))
			}
			messages[idx].ErrorCode = record.ErrorCode
			messages[idx].ErrorMessage = record.ErrorMessage
			messages[idx].FailCount++
			retries = append(retries, messages[idx])
		}
	}
	if len(retries) > 0 {
		err = errs.ErrRetryRecords
	}
	return retries, err
}
