package engine

import (
	"fmt"
	"regexp"
	"runtime"
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

// RecordFunctionDuration records the duration of the function
// It should be called at the start of the function like so  `defer e.metrics.RecordFunctionDuration()()`
func (o *MetricsRecorder) RecordFunctionDuration() func() {
	start := time.Now()
	return func() {
		elapsed := time.Since(start)

		// Skip this function, and fetch the PC for its parent.
		pc, _, _, _ := runtime.Caller(1)

		// Retrieve a function object this function's parent.
		funcObj := runtime.FuncForPC(pc)

		// Use a regex to strip out the module path
		funcNameRegex := regexp.MustCompile(`^.*\.(.*)$`)
		name := funcNameRegex.ReplaceAllString(funcObj.Name(), "$1")

		timer := o.metrics.Timer(o.addMyAddress(name))
		timer.Update(elapsed)
	}
}

// RecordObjectiveStarted records metrics about the start of an objective
// This should be called when an objective is first created
func (o *MetricsRecorder) RecordObjectiveStarted(id protocols.ObjectiveId) {
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

	delete(o.startTimes, id)
}

// RecordQueueLength records metrics about the length of some queue
func (o *MetricsRecorder) RecordQueueLength(name string, queueLength int) {
	o.metrics.Gauge(o.addMyAddress(name)).Update(float64(queueLength))
}

func (o *MetricsRecorder) addMyAddress(name string) string {
	return fmt.Sprintf("%s,wallet=%s", name, o.me.String())
}
