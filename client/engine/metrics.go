package engine

import (
	"fmt"
	"time"

	"github.com/statechannels/go-nitro/client/engine/store/safesync"
)

// The API for recording metrics.
type MetricsApi interface {
	// RecordMetric records metric information in some data store.
	// The metricName is the name of the metric. IE: "message-count"
	// The value is the value to record for the metric IE: the count of messages
	// additionalData is any additional data to record with the metric. IE: The address of the engine
	RecordMetric(metricName string, value float64, additionalData map[string]string)
}

// MetricRecorder is a wrapper around a MetricsApi that helps in recording elapsed time for operations.
type MetricRecorder struct {
	metricsApi MetricsApi
	operations safesync.Map[operation]
}

// operation is used to keep track of when an operation started.
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
func (r *MetricRecorder) RecordOperationDuration(operationName string, operationFunc func(), additionalData map[string]string) {

	startTime := time.Now()
	operationFunc()
	stopTime := time.Now()

	additionalData["startTime"] = fmt.Sprint(startTime.Nanosecond())
	additionalData["stopTime"] = fmt.Sprint(stopTime.Nanosecond())

	r.metricsApi.RecordMetric(operationName, float64(stopTime.Sub(startTime).Nanoseconds()), additionalData)

}

// MarkOperationStart is used to mark the start of some operation.
func (r *MetricRecorder) MarkOperationStart(operationName string, operationId string, additionalData map[string]string) {
	r.operations.Store(operationId, operation{Name: operationName, Id: operationId, StartTime: time.Now(), AdditionalData: additionalData})
}

// MarkOperationStop is used to indicate the operation has now completed. This records the elapsed time for the operation using the metrics API
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

//  RecordMetric records any arbitrary metric information.
func (r *MetricRecorder) RecordMetric(metricName string, value float64, additionalData map[string]string) {
	r.RecordMetric(metricName, value, additionalData)
}

// NoOpMetricsApi implements the MetricsAPI interface and does nothing.
// It is used when a metric API is not provided to the client.
type NoOpMetricsApi struct{}

func (r *NoOpMetricsApi) RecordMetric(metricName string, value float64, additionalData map[string]string) {
}
