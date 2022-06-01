package engine

import (
	"fmt"
	"time"

	"github.com/statechannels/go-nitro/client/engine/store/safesync"
	"github.com/statechannels/go-nitro/types"
)

// MetricsApi defines the API for recording metrics the engine expects.
// This is used by the engine to record various metrics, like the duration of an objective.
type MetricsApi interface {
	// RecordMetric records metric information in some data store.
	// The metricName is the name of the metric. IE: "message-count"
	// The value is the value to record for the metric IE: the count of messages
	// additionalData is any additional data to record with the metric. IE: The address of the engine
	RecordMetric(metricName string, value float64, additionalData map[string]string)
}

// MetricRecorder is a wrapper around a MetricsApi that aids in recording duration metrics.
type MetricRecorder struct {
	metricsApi MetricsApi
	operations safesync.Map[operation]
}

// operation is used to keep track of when an operation started internally in the MetricRecorder
type operation struct {
	Id             string
	Name           string
	StartTime      time.Time
	AdditionalData map[string]string
}

// NewMetricRecorder creates a new MetricRecorder given a MetricsApi.
func NewMetricRecorder(metricsApi MetricsApi) *MetricRecorder {
	return &MetricRecorder{
		metricsApi: metricsApi,
		operations: safesync.Map[operation]{},
	}
}

// RecordOperationDuration records the duration of some operation.
// It does this by calling the provided operationFunc and recording the duration.
// The metric information recorded will also contain the wallet address and start/stop timestamps
func (r *MetricRecorder) RecordOperationDuration(metricName string, operationFunc func(), walletAddress types.Address) {

	startTime := time.Now()
	operationFunc()
	stopTime := time.Now()
	additionalData := map[string]string{"wallet": walletAddress.String()}
	additionalData["startTime"] = fmt.Sprint(startTime.Nanosecond())
	additionalData["stopTime"] = fmt.Sprint(stopTime.Nanosecond())

	r.metricsApi.RecordMetric(metricName, float64(stopTime.Sub(startTime).Nanoseconds()), additionalData)

}

// MarkOperationStart is used to record that some operation has started.
// MarkOperationStop can then be called when the operation is complete to record the duration of the operation.
// The operationId is some unique identifier for the operation that can be later be passed into MarkOperationStop
func (r *MetricRecorder) MarkOperationStart(metricName string, operationId string, walletAddress types.Address) {
	additionalData := map[string]string{"wallet": walletAddress.String()}
	r.operations.Store(operationId, operation{Name: metricName, Id: operationId, StartTime: time.Now(), AdditionalData: additionalData})
}

// MarkOperationStop is used to indicate the operation has now completed.
// This records the elapsed time for the operation using the metrics API
func (r *MetricRecorder) MarkOperationStop(operationId string) {
	operation, ok := r.operations.Load(operationId)
	if !ok {
		return
	}
	completeTime := time.Now()
	operation.AdditionalData["startTime"] = fmt.Sprint(operation.StartTime.Nanosecond())
	operation.AdditionalData["stopTime"] = fmt.Sprint(completeTime.Nanosecond())
	operation.AdditionalData["id"] = operation.Id
	r.metricsApi.RecordMetric(operation.Name, float64(completeTime.Sub(operation.StartTime).Nanoseconds()), operation.AdditionalData)
	r.operations.Delete(operationId)
}

//  RecordMetric records a float64 value for the metric specified by metricName.
// The metric information recorded will contain a timestamp of the current time and the wallet address.
func (r *MetricRecorder) RecordMetric(metricName string, value float64, walletAddress types.Address) {
	additionalData := map[string]string{"wallet": walletAddress.String()}
	r.metricsApi.RecordMetric(metricName, value, additionalData)
}

// NoOpMetricsApi implements the MetricsAPI interface and does nothing.
// It is used when a metric API is not provided to the client.
type NoOpMetricsApi struct{}

func (r *NoOpMetricsApi) RecordMetric(metricName string, value float64, additionalData map[string]string) {
}
