package engine

import (
	"fmt"
	"strings"
	"time"

	"github.com/rcrowley/go-metrics"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

// MetricsAPI is an interface for recording metrics
// It is heavily based on https://github.com/testground/sdk-go/blob/master/runtime/metrics_api.go
// It exposes some basic functionality that is useful for recording engine metrics
type MetricsApi interface {
	// RecordPoint records a float64 point under the provided metric name + tags.
	//
	// The format of the metric name is a comma-separated list, where the first
	// element is the metric name, and optionally, an unbounded list of
	// key-value pairs. Example:
	//
	//   requests_received,tag1=value1,tag2=value2,tag3=value3
	RecordPoint(name string, value float64)

	// Counter creates a measurement of counter type. The returned type is an
	// alias of go-metrics' Counter type. Refer to godocs there for details.
	//
	// The format of the metric name is a comma-separated list, where the first
	// element is the metric name, and optionally, an unbounded list of
	// key-value pairs. Example:
	//
	//   requests_received,tag1=value1,tag2=value2,tag3=value3
	Counter(name string) metrics.Counter

	// Timer creates a measurement of timer type.
	// The returned type is an alias of go-metrics' Timer type. Refer to
	// godocs there for details.
	//
	// The format of the metric name is a comma-separated list, where the first
	// element is the metric name, and optionally, an unbounded list of
	// key-value pairs. Example:
	//
	//   requests_received,tag1=value1,tag2=value2,tag3=value3
	Timer(name string) metrics.Timer

	// Gauge creates a measurement of gauge type (float64).
	// The returned type is an alias of go-metrics' GaugeFloat64 type. Refer to
	// godocs there for details.
	//
	// The format of the metric name is a comma-separated list, where the first
	// element is the metric name, and optionally, an unbounded list of
	// key-value pairs. Example:
	//
	//   requests_received,tag1=value1,tag2=value2,tag3=value3
	Gauge(name string) metrics.GaugeFloat64
}

// NewNoOpMetrics returns a MetricsApi that does nothing.
type NoOpMetrics struct{}

func (n *NoOpMetrics) Counter(name string) metrics.Counter {
	return metrics.NilCounter{}
}

func (n *NoOpMetrics) Timer(name string) metrics.Timer {
	return metrics.NilTimer{}
}

func (n *NoOpMetrics) RecordPoint(name string, value float64) {

}
func (n *NoOpMetrics) Gauge(name string) metrics.GaugeFloat64 {
	return metrics.NilGaugeFloat64{}
}

// MetricsRecorder is used to record metrics about the engine
type MetricsRecorder struct {
	startTimes map[protocols.ObjectiveId]time.Time
	me         types.Address
	metrics    MetricsApi
}

// NewMetricsRecorder returns a new MetricsRecorder that uses the metricsApi to record metrics
func NewMetricsRecorder(me types.Address, metrics MetricsApi) *MetricsRecorder {
	return &MetricsRecorder{
		startTimes: make(map[protocols.ObjectiveId]time.Time),
		me:         me,
		metrics:    metrics,
	}
}

// RecordDuration records the duration of the given function for the metric specified by name
func (o *MetricsRecorder) RecordDuration(name string, funcToTime func()) {

	timer := o.metrics.Timer(o.addMyAddress(name))
	// A nil timer's Time function does nothing so we need to manually call funcToTime
	if _, isNilTimer := timer.(metrics.NilTimer); isNilTimer {
		funcToTime()
	} else {
		timer.Time(funcToTime)
	}

}

// RecordObjectiveStarted records metrics about the start of an objective
// This should be called when an objective is first created
func (o *MetricsRecorder) RecordObjectiveStarted(id protocols.ObjectiveId) {
	o.metrics.Counter(o.addMyAddress("active_objective_count")).Inc(1)
	o.startTimes[id] = time.Now()
}

// RecordObjectiveCompleted records metrics about the completion of an objective
// This should be called when an objective is completed
func (o *MetricsRecorder) RecordObjectiveCompleted(id protocols.ObjectiveId) {
	start := o.startTimes[id]

	elapsed := time.Since(start)
	oType := strings.Split(string(id), "-")[0]

	timer := o.metrics.Timer(o.addMyAddress("objective_complete_time") + fmt.Sprintf(",type=%s", oType))
	timer.Update(elapsed)
	o.metrics.Counter(o.addMyAddress("objective_complete_count")).Inc(1)
	o.metrics.Counter(o.addMyAddress("active_objective_count")).Dec(1)

	delete(o.startTimes, id)

}

// RecordOutgoingMessage records metrics about the the outgoing message
func (o *MetricsRecorder) RecordOutgoingMessage(msg protocols.Message) {
	proposalCount := len(msg.SignedProposals())
	o.metrics.RecordPoint(o.addMyAddress("message_proposals")+fmt.Sprintf(",to=%s", msg.To), float64(proposalCount))
	stateCount := len(msg.SignedProposals())
	o.metrics.RecordPoint(o.addMyAddress("message_states")+fmt.Sprintf(",to=%s", msg.To), float64(stateCount))
}

// RecordQueueLength records metrics about the length of some queue
func (o *MetricsRecorder) RecordQueueLength(name string, queueLength int) {
	o.metrics.Gauge(o.addMyAddress(name)).Update(float64(queueLength))
}

func (o *MetricsRecorder) addMyAddress(name string) string {
	return fmt.Sprintf("%s,wallet=%s", name, o.me.String())
}
