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

func mustAtoi(s string) int {
	num, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return num
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
		statuses = append(statuses, Status{
			Function:         toks[0],
			Total:            mustAtoi(toks[1]),
			Running:          mustAtoi(toks[2]),
			AvailableWorkers: mustAtoi(toks[3]),
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
