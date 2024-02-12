package metric

import (
	"context"
	"fmt"
	"testing"

	mock_statsd "github.com/coopnorge/go-datadog-lib/v2/internal/generated/mocks/DataDog/datadog-go/v5/statsd"
	mock_metrics "github.com/coopnorge/go-datadog-lib/v2/internal/generated/mocks/metric"

	gomock "go.uber.org/mock/gomock"
)

func TestAddMetric(t *testing.T) {
	testCases := []struct {
		name        string
		metricType  Type
		isWithError bool
	}{
		{
			name:       "MetricTypeEvent",
			metricType: MetricTypeEvent,
		},
		{
			name:       "MetricTypeMeasurement",
			metricType: MetricTypeMeasurement,
		},
		{
			name:       "MetricTypeCountEvents",
			metricType: MetricTypeCountEvents,
		},
		{
			name:        "metricCollectionErr",
			metricType:  MetricTypeEvent,
			isWithError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockDatadogClient := mock_metrics.NewMockDatadogMetricsClient(ctrl)
			mockDatadogStatsd := mock_statsd.NewMockClientInterface(ctrl)
			ctrl.Finish()
			tMetricData := Data{
				Name:  "RuntimeTest",
				Type:  tc.metricType,
				Value: float64(42),
				MetricTags: []Tag{
					{Name: "Unit", Value: "Test"},
				},
			}

			mockDatadogClient.
				EXPECT().
				GetClient().
				Return(mockDatadogStatsd).
				MaxTimes(2)

			mockDatadogClient.
				EXPECT().
				GetServiceNamePrefix().
				Return("metrics").
				MaxTimes(1)

			if tc.isWithError {
				mockDatadogStatsd.
					EXPECT().
					Incr(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(fmt.Errorf("datadog statsd have error here")).
					MaxTimes(1)
			} else if tc.metricType == MetricTypeEvent {
				mockDatadogStatsd.
					EXPECT().
					Incr(
						"metrics.RuntimeTest",
						[]string{"unit:Test"},
						float64(1),
					).
					Return(nil).
					MaxTimes(1)
			} else if tc.metricType == MetricTypeMeasurement {
				mockDatadogStatsd.
					EXPECT().
					Gauge(
						"metrics.RuntimeTest",
						tMetricData.Value,
						[]string{"unit:Test"},
						float64(1),
					).
					Return(nil).
					MaxTimes(1)
			} else if tc.metricType == MetricTypeCountEvents {
				mockDatadogStatsd.
					EXPECT().
					Count(
						"metrics.RuntimeTest",
						int64(tMetricData.Value),
						[]string{"unit:Test"},
						float64(1),
					).
					Return(nil).
					MaxTimes(1)
			}

			bmc := &BaseMetricCollector{DatadogMetrics: mockDatadogClient}
			bmc.AddMetric(context.Background(), tMetricData)
		})
	}
}

func TestAddMetricNoClient(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockDatadogClient := mock_metrics.NewMockDatadogMetricsClient(ctrl)
	ctrl.Finish()

	mockDatadogClient.EXPECT().GetClient().Return(nil).MaxTimes(1)

	bmc := &BaseMetricCollector{DatadogMetrics: mockDatadogClient}
	bmc.AddMetric(context.Background(), Data{})
}
