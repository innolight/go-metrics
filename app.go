package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/innolight/go-metrics/promdb"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	// Open connection to a DB (could also use the https://github.com/jmoiron/sqlx library)
	db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost:5432/app?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	// configure database connection
	db.SetConnMaxIdleTime(time.Second * 10)
	db.SetConnMaxLifetime(time.Second * 30)
	db.SetMaxIdleConns(100)
	db.SetMaxOpenConns(200)

	// Create a new collector, the name will be used as a label on the metrics
	collector := promdb.NewCollector("app_db", db)

	// Register it with Prometheus
	prometheus.MustRegister(collector)

	// Register the metrics handler
	http.Handle("/metrics", promhttp.Handler())

	// Run the web server
	http.ListenAndServe(":8080", nil)
}
