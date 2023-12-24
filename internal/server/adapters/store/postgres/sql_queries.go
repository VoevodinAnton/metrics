package postgres

const (
	getCounterMetricQuery = `SELECT name, sum(value) FROM counter_metrics WHERE name = $1
		GROUP BY name;`
	getGaugeMetricQuery = `SELECT name, value FROM gauge_metrics WHERE name = $1
		ORDER BY updated_at DESC LIMIT 1;`
	insertGaugeMetricQuery   = `INSERT INTO gauge_metrics (name, value, updated_at) VALUES ($1, $2, $3);`
	insertCounterMetricQuery = `INSERT INTO counter_metrics (name, value, updated_at) VALUES ($1, $2, $3);`
	getGaugeMetricsQuery     = `SELECT name, value FROM gauge_metrics gm1 WHERE updated_at  = (
		SELECT MAX(updated_at)
		FROM gauge_metrics gm2
		WHERE gm2.name = gm1.name
	);`
	getCounterMetricsQuery = `SELECT name, value FROM counter_metrics cm1 WHERE updated_at  = (
		SELECT MAX(updated_at)
		FROM counter_metrics cm2
		WHERE cm2.name = cm1.name
	);`

	insertGaugeMetricQueryName   = "insertGaugeMetricQuery"
	insertCounterMetricQueryName = "insertCounterMetricQuery"
)
