package kinetic

import (
	"math/rand"
	"testing"
	"time"

	metrics "github.com/jasonyurs/go-metrics"
	. "github.com/smartystreets/goconvey/convey"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func TestConsumerStatsCollector(t *testing.T) {
	Convey("given a NilStatsCollector", t, func() {
		var sc ConsumerStatsCollector = &NilConsumerStatsCollector{}
		So(sc, ShouldNotBeNil)

		Convey("check that AddConsumed does not error", func() {
			sc.AddConsumed(1)
		})

		Convey("check that AddDelivered does not error", func() {
			sc.AddDelivered(1)
		})

		Convey("check that AddProcessed does not error", func() {
			sc.AddProcessed(1)
		})

		Convey("check that AddBatchSize does not error", func() {
			sc.AddBatchSize(1)
		})

		Convey("check that AddGetRecordsCalled does not error", func() {
			sc.AddGetRecordsCalled(1)
		})

		Convey("check that AddGetRecordsProvisionedThroughputExceeded does not error", func() {
			sc.AddGetRecordsProvisionedThroughputExceeded(1)
		})

		Convey("check that AddGetRecordsTimeout does not error", func() {
			sc.AddGetRecordsTimeout(1)
		})

		Convey("check that AddGetRecordsReadTimeout does not error", func() {
			sc.AddGetRecordsReadTimeout(1)
		})

		Convey("check that AddProcessedDuration does not error", func() {
			sc.AddProcessedDuration(1)
		})

		Convey("check that AddGetRecordsDuration does not error", func() {
			sc.AddGetRecordsDuration(1)
		})

		Convey("check that AddGetRecordsReadResponseDuration does not error", func() {
			sc.AddGetRecordsReadResponseDuration(1)
		})

		Convey("check that AddGetRecordsUnmarshalDuration does not error", func() {
			sc.AddGetRecordsUnmarshalDuration(1)
		})
	})

	Convey("given a DefaultStatsCollector", t, func() {
		r := metrics.NewRegistry()
		var sc ConsumerStatsCollector = NewDefaultConsumerStatsCollector(r)
		So(sc, ShouldNotBeNil)

		Convey("check that AddConsumed does not error", func() {
			count := rand.Int()
			sc.AddConsumed(count)
			So(sc.(*DefaultConsumerStatsCollector).Consumed.Count(), ShouldEqual, int64(count))
		})

		Convey("check that AddDelivered does not error", func() {
			count := rand.Int()
			sc.AddDelivered(count)
			So(sc.(*DefaultConsumerStatsCollector).Delivered.Count(), ShouldEqual, int64(count))
		})

		Convey("check that AddProcessed does not error", func() {
			count := rand.Int()
			sc.AddProcessed(count)
			So(sc.(*DefaultConsumerStatsCollector).Processed.Count(), ShouldEqual, int64(count))
		})

		Convey("check that AddBatchSize does not error", func() {
			count := rand.Int()
			sc.AddBatchSize(count)
			So(sc.(*DefaultConsumerStatsCollector).BatchSize.Count(), ShouldEqual, int64(count))
		})

		Convey("check that AddGetRecordsCalled does not error", func() {
			count := rand.Int()
			sc.AddGetRecordsCalled(count)
			So(sc.(*DefaultConsumerStatsCollector).GetRecordsCalled.Count(), ShouldEqual, int64(count))
		})

		Convey("check that AddGetRecordsProvisionedThroughputExceeded does not error", func() {
			count := rand.Int()
			sc.AddGetRecordsProvisionedThroughputExceeded(count)
			So(sc.(*DefaultConsumerStatsCollector).GetRecordsProvisionedThroughputExceeded.Count(), ShouldEqual, int64(count))
		})

		Convey("check that AddGetRecordsTimeout does not error", func() {
			count := rand.Int()
			sc.AddGetRecordsTimeout(count)
			So(sc.(*DefaultConsumerStatsCollector).GetRecordsTimeout.Count(), ShouldEqual, int64(count))
		})

		Convey("check that AddGetRecordsReadTimeout does not error", func() {
			count := rand.Int()
			sc.AddGetRecordsReadTimeout(count)
			So(sc.(*DefaultConsumerStatsCollector).GetRecordsReadTimeout.Count(), ShouldEqual, int64(count))
		})

		Convey("check that AddProcessedDuration does not error", func() {
			sc.AddProcessedDuration(time.Second)
			So(sc.(*DefaultConsumerStatsCollector).ProcessedDuration.Value(), ShouldEqual, 1000000000)
		})

		Convey("check that AddGetRecordsDuration does not error", func() {
			sc.AddGetRecordsDuration(time.Second)
			So(sc.(*DefaultConsumerStatsCollector).GetRecordsDuration.Value(), ShouldEqual, 1000000000)
		})

		Convey("check that AddGetRecordsReadResponseDuration does not error", func() {
			sc.AddGetRecordsReadResponseDuration(time.Second)
			So(sc.(*DefaultConsumerStatsCollector).GetRecordsReadResponseDuration.Value(), ShouldEqual, 1000000000)
		})

		Convey("check that AddGetRecordsUnmarshalDuration does not error", func() {
			sc.AddGetRecordsUnmarshalDuration(time.Second)
			So(sc.(*DefaultConsumerStatsCollector).GetRecordsUnmarshalDuration.Value(), ShouldEqual, 1000000000)
		})
	})
}