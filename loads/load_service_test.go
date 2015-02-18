package loads_test

import (
	"testing"
	"time"

	. "github.com/jaimegildesagredo/gomonitor/loads"
)

const (
	A_DELAY        = 10 * time.Millisecond
	A_LOAD_ONE     = 1.69
	A_LOAD_FIVE    = 2.32
	A_LOAD_FIFTEEN = 4.48
)

func TestMonitorSystemLoad(t *testing.T) {
	loadService := NewLoadService(newInMemoryLoadRepository(Load{
		A_LOAD_ONE, A_LOAD_FIVE, A_LOAD_FIFTEEN}))

	loads := loadService.Monitor(A_DELAY)

	load := <-loads

	if load[0] != A_LOAD_ONE || load[1] != A_LOAD_FIVE || load[2] != A_LOAD_FIFTEEN {
		t.Fatal("Invalid load values", load, "expected", A_LOAD_ONE, A_LOAD_FIVE, A_LOAD_FIFTEEN)
	}
}

func newInMemoryLoadRepository(load Load) LoadRepository {
	return &inMemoryLoadRepository{
		load: load,
	}
}

type inMemoryLoadRepository struct {
	load Load
}

func (repo *inMemoryLoadRepository) Get() Load {
	return repo.load
}
