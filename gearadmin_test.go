package gearadmin

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

// MockGearmand implements io.ReadWriteCloser and can be used in place of an actual
// connection to gearmand.
type MockGearmand struct {
	buf       bytes.Buffer
	Responses map[string]string
}

func (mg *MockGearmand) Read(p []byte) (int, error) {
	return mg.buf.Read(p)
}

func (mg *MockGearmand) Write(p []byte) (int, error) {
	cmd := strings.Trim(string(p), "\n")
	val, ok := mg.Responses[cmd]
	if !ok {
		mg.buf.Write([]byte("ERR UNKNOWN_COMMAND"))
	} else {
		mg.buf.Write([]byte(val))
	}
	return len(p), nil
}

func (mg *MockGearmand) Close() error {
	return nil
}

func TestStatus(t *testing.T) {
	mockGearmand := MockGearmand{}
	defer mockGearmand.Close()
	mockGearmand.Responses = map[string]string{
		"status": `fn	3	2	1
.`,
	}
	ga := GearmanAdmin{&mockGearmand}
	statuses, err := ga.Status()
	if err != nil {
		t.Fatal(err)
	}
	if len(statuses) != 1 {
		t.Fatalf("Expected one status")
	}
	if statuses[0].Function != "fn" {
		t.Fatalf("Incorrect function: expected 'fn', got '%s'", statuses[0].Function)
	}
	if statuses[0].Total != 3 {
		t.Fatalf("Incorrect total: expected 3, got '%d'", statuses[0].Total)
	}
	if statuses[0].Running != 2 {
		t.Fatalf("Incorrect running: expected 2, got '%d'", statuses[0].Running)
	}
	if statuses[0].AvailableWorkers != 1 {
		t.Fatalf("Incorrect running: expected 1, got '%d'", statuses[0].AvailableWorkers)
	}
}

func TestWorkers(t *testing.T) {
	mockGearmand := MockGearmand{}
	defer mockGearmand.Close()
	mockGearmand.Responses = map[string]string{
		"workers": `74 10.0.1.167 - :
284 10.0.2.16 - : fn1
284 10.0.2.16 - : fn1 fn2
.`,
	}
	ga := GearmanAdmin{&mockGearmand}
	workers, err := ga.Workers()
	if err != nil {
		t.Fatal(err)
	}
	if len(workers) != 3 {
		t.Fatalf("Expected three workers: %v", workers)
	}
	checkWorker := func(worker Worker, fd, ip, cid string, functions []string) {
		if worker.Fd != fd {
			t.Fatalf("Expected worker Fd '%s', got '%s'", fd, worker.Fd)
		}
		if worker.IPAddress != ip {
			t.Fatalf("Expected worker IPAddress '%s', got '%s'", ip, worker.IPAddress)
		}
		if worker.ClientID != cid {
			t.Fatalf("Expected worker ClientID '%s', got '%s'", cid, worker.ClientID)
		}
		if fmt.Sprintf("%v", worker.Functions) != fmt.Sprintf("%v", functions) {
			t.Fatalf("Expected worker Functions '%#v', got '%#v'", functions, worker.Functions)
		}
	}
	checkWorker(workers[0], "74", "10.0.1.167", "-", []string{})
	checkWorker(workers[1], "284", "10.0.2.16", "-", []string{"fn1"})
	checkWorker(workers[2], "284", "10.0.2.16", "-", []string{"fn1", "fn2"})
}
