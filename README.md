# go-metrics

A set of Go libraries to collect metrics for various infrastructure components and exporting them in Prometheus format.

## promdb module

promdb module collects database metrics ([sql.DBStats](https://golang.org/pkg/database/sql/#DBStats)) of database connection pools and expose the metrics in the Prometheus format.

**Install**
```
go get github.com/innolight/go-metrics/promdb
```

**Example usage**

```go
package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/innolight/go-metrics/promdb"
)

func main() {
	// Create a database connection pool
	db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost:5432/app?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	// Create a new collector for database named "app_db"
	collector := promdb.NewCollector("app_db", db)

	// Register Prometheus collector
	prometheus.MustRegister(collector)

	// Register the metrics handler
	http.Handle("/metrics", promhttp.Handler())

	// Run the web server
	http.ListenAndServe(":8080", nil)
}
```

To quickly get metrics from the application above:
```bash
curl localhost:8080/metrics | grep db
```

**Exposed Metrics**

| Name                                 | Type    | Description                                                       | Labels  |
|--------------------------------------|---------|-------------------------------------------------------------------|---------|
| db_connections_max_open              | Gauge   | Maximum number of open connections to the database.               | db_name |
| db_connections_open                  | Gauge   | The number of established connections both in use and idle.       | db_name |
| db_connections_in_use                | Gauge   | The number of connections currently in use.                       | db_name |
| db_connections_idle                  | Gauge   | The number of idle connections.                                   | db_name |
| db_connections_wait_count            | Counter | The total number of connections waited for.                       | db_name |
| db_connections_wait_duration_seconds | Counter | The total time blocked waiting for a new connection.              | db_name |
| db_connections_max_idle_closed       | Counter | The total number of connections closed due to SetMaxIdleConns.    | db_name |
| db_connections_max_idle_time_closed  | Counter | The total number of connections closed due to SetConnMaxIdleTime. | db_name |
| db_connections_max_lifetime_closed   | Counter | The total number of connections closed due to SetConnMaxLifetime. | db_name |

 
