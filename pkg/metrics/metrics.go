package metrics

import (
	"github.com/operator-framework/operator-lifecycle-manager/pkg/api/client/clientset/versioned"
	"github.com/prometheus/client_golang/prometheus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// TODO(alecmerdler): Can we use this to emit Kubernetes events?
type MetricsProvider interface {
	HandleMetrics() error
}

type metricsCSV struct {
	client versioned.Interface
}

func NewMetricsCSV(client versioned.Interface) MetricsProvider {
	return &metricsCSV{client}
}

func (m *metricsCSV) HandleMetrics() error {
	cList, err := m.client.OperatorsV1alpha1().ClusterServiceVersions(metav1.NamespaceAll).List(metav1.ListOptions{})
	if err != nil {
		return err
	}
	csvCount.Set(float64(len(cList.Items)))
	return nil
}

type metricsInstallPlan struct {
	client versioned.Interface
}

func NewMetricsInstallPlan(client versioned.Interface) MetricsProvider {
	return &metricsInstallPlan{client}
}

func (m *metricsInstallPlan) HandleMetrics() error {
	cList, err := m.client.OperatorsV1alpha1().InstallPlans(metav1.NamespaceAll).List(metav1.ListOptions{})
	if err != nil {
		return err
	}
	installPlanCount.Set(float64(len(cList.Items)))
	return nil
}

type metricsSubscription struct {
	client versioned.Interface
}

func NewMetricsSubscription(client versioned.Interface) MetricsProvider {
	return &metricsSubscription{client}
}

func (m *metricsSubscription) HandleMetrics() error {
	cList, err := m.client.OperatorsV1alpha1().Subscriptions(metav1.NamespaceAll).List(metav1.ListOptions{})
	if err != nil {
		return err
	}
	subscriptionCount.Set(float64(len(cList.Items)))
	return nil
}

type metricsCatalogSource struct {
	client versioned.Interface
}

func NewMetricsCatalogSource(client versioned.Interface) MetricsProvider {
	return &metricsCatalogSource{client}

}

func (m *metricsCatalogSource) HandleMetrics() error {
	cList, err := m.client.OperatorsV1alpha1().CatalogSources(metav1.NamespaceAll).List(metav1.ListOptions{})
	if err != nil {
		return err
	}
	catalogSourceCount.Set(float64(len(cList.Items)))
	return nil
}

type MetricsNil struct{}

func NewMetricsNil() MetricsProvider {
	return &MetricsNil{}
}

func (*MetricsNil) HandleMetrics() error {
	return nil
}

// To add new metrics:
// 1. Register new metrics in Register() below.
// 2. Add appropriate metric updates in HandleMetrics (or elsewhere instead).
var (
	csvCount = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "csv_count",
			Help: "Number of CSVs successfully registered",
		},
	)

	installPlanCount = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "install_plan_count",
			Help: "Number of install plans",
		},
	)

	subscriptionCount = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "subscription_count",
			Help: "Number of subscriptions",
		},
	)

	catalogSourceCount = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "catalog_source_count",
			Help: "Number of catalog sources",
		},
	)

	// exported since it's not handled by HandleMetrics
	CSVUpgradeCount = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "csv_upgrade_count",
			Help: "Monotonic count of catalog sources",
		},
	)
)

func Register() {
	prometheus.MustRegister(csvCount)
	prometheus.MustRegister(installPlanCount)
	prometheus.MustRegister(subscriptionCount)
	prometheus.MustRegister(catalogSourceCount)
	prometheus.MustRegister(CSVUpgradeCount)
}
