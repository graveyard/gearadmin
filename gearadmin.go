package gearadmin

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// GearmanAdmin communicates with a gearman server.
type GearmanAdmin struct {
	Conn io.ReadWriter
}

// Status represents the status of a queue for a function as returned by the "status" command.
type Status struct {
	Function         string
	Total            int
	Running          int
	AvailableWorkers int
}

// Worker represents a worker connected to gearman as returned by the "workers" command.
type Worker struct {
	Fd        string
	IPAddress string
	ClientID  string
	Functions []string
}

// Status returns the status of all function queues.
func (ga GearmanAdmin) Status() ([]Status, error) {
	var statuses []Status
	fmt.Fprintf(ga.Conn, "status\n")
	scanner := bufio.NewScanner(ga.Conn)
	for scanner.Scan() && scanner.Text() != "." {
		toks := strings.Split(scanner.Text(), "\t")
		if len(toks) != 4 {
			return statuses, fmt.Errorf("unexpected status: '%s'", scanner.Text())
		}
		total, err := strconv.Atoi(toks[1])
		if err != nil {
			return statuses, fmt.Errorf("could not parse total: '%s'", scanner.Text())
		}
		running, err := strconv.Atoi(toks[2])
		if err != nil {
			return statuses, fmt.Errorf("could not parse running: '%s'", scanner.Text())
		}
		available, err := strconv.Atoi(toks[3])
		if err != nil {
			return statuses, fmt.Errorf("could not parse available: '%s'", scanner.Text())
		}
		statuses = append(statuses, Status{
			Function:         toks[0],
			Total:            total,
			Running:          running,
			AvailableWorkers: available,
		})
	}
	return statuses, scanner.Err()
}

// Workers returns a summary of workers connected to gearman.
func (ga GearmanAdmin) Workers() ([]Worker, error) {
	var workers []Worker
	fmt.Fprintf(ga.Conn, "workers\n")
	scanner := bufio.NewScanner(ga.Conn)
	for scanner.Scan() && scanner.Text() != "." {
		toks := strings.Split(scanner.Text(), " ")
		if len(toks) < 4 {
			return workers, fmt.Errorf("unexpected worker: '%s'", scanner.Text())
		}
		workers = append(workers, Worker{
			Fd:        toks[0],
			IPAddress: toks[1],
			ClientID:  toks[2],
			Functions: toks[4:],
		})
	}
	return workers, scanner.Err()
}
