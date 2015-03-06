package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/url"
	"time"

	"github.com/Clever/gearadmin"
	"gopkg.in/Clever/kayvee-go.v1"
)

func logMetrics(g gearadmin.GearmanAdmin) {
	statuses, err := g.Status()
	if err != nil {
		log.Fatalf("error retrieving gearman status: %s", err)
	}
	for _, status := range statuses {
		fmt.Println(kayvee.FormatLog("gearlogger", "info", "status", map[string]interface{}{
			"function": status.Function,
			"running":  status.Running,
			"total":    status.Total,
			"workers":  status.AvailableWorkers,
		}))
	}
}

func main() {
	host := flag.String("host", "tcp://localhost:4730", "gearman host")
	interval := flag.Duration("interval", time.Minute, "interval to log at")
	flag.Parse()

	url, err := url.Parse(*host)
	if err != nil {
		log.Fatalf("error parsing gearman url %s: %s", *host, err)
	}

	conn, err := net.Dial(url.Scheme, url.Host)
	if err != nil {
		log.Fatalf("error connecting to gearman: %s", err)
	}
	defer conn.Close()

	g := gearadmin.NewGearmanAdmin(conn)

	logMetrics(g)
	for _ = range time.Tick(*interval) {
		logMetrics(g)
	}
}
