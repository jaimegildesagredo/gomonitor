package network_test

import (
	"errors"
	"testing"
	"time"

	. "github.com/jaimegildesagredo/gomonitor/network"
)

const (
	FIRST_TX_BYTES          = 1350
	LAST_TX_BYTES           = 1500
	FIRST_RX_BYTES          = 2048
	LAST_RX_BYTES           = 4096
	A_DELAY                 = 10 * time.Millisecond
	EXPECTED_BANDWIDTH_UP   = 15000
	EXPECTED_BANDWIDTH_DOWN = 204800
	AN_INTERFACE_NAME       = "a-interface-name"
)

func TestMonitorBandwidth(t *testing.T) {
	bytesRepository := newInMemoryBytesRepository([]int{FIRST_TX_BYTES, LAST_TX_BYTES}, []int{FIRST_RX_BYTES, LAST_RX_BYTES})
	bandwidthService := NewBandwidthService(bytesRepository)
	bandwidths, _ := bandwidthService.MonitorBandwidth(AN_INTERFACE_NAME, A_DELAY)

	bandwidth := <-bandwidths

	if bandwidth.Up != EXPECTED_BANDWIDTH_UP {
		t.Fatal("Invalid bandwidth up value", bandwidth.Up, "expected", EXPECTED_BANDWIDTH_UP)
	}

	if bandwidth.Down != EXPECTED_BANDWIDTH_DOWN {
		t.Fatal("Invalid bandwidth down value", bandwidth.Down, "expected", EXPECTED_BANDWIDTH_DOWN)
	}
}

func TestMonitorBandwidthWhenInterfaceDoesNotExists(t *testing.T) {
	bytesRepository := newInMemoryBytesRepository([]int{}, []int{})
	bandwidthService := NewBandwidthService(bytesRepository)

	_, err := bandwidthService.MonitorBandwidth(AN_INTERFACE_NAME, A_DELAY)

	if err == nil {
		t.Fatal("Expected an error when interface does not exists")
	}
}

func newInMemoryBytesRepository(txBytes []int, rxBytes []int) BytesRepository {
	repository := inMemoryBytesRepository{
		txBytes: txBytes,
		rxBytes: rxBytes,
	}
	return &repository
}

type inMemoryBytesRepository struct {
	txBytes []int
	rxBytes []int
}

func (repo *inMemoryBytesRepository) GetTx(interfaceName string) (int, error) {
	var value int
	if len(repo.txBytes) > 0 {
		value = repo.txBytes[0]
	}

	if len(repo.txBytes) > 1 {
		repo.txBytes = repo.txBytes[1:]
	} else {
		repo.txBytes = []int{}
	}

	if value == 0 {
		return value, errors.New("")
	}

	return value, nil
}

func (repo *inMemoryBytesRepository) GetRx(interfaceName string) (int, error) {
	var value int
	if len(repo.rxBytes) > 0 {
		value = repo.rxBytes[0]
	}

	if len(repo.rxBytes) > 1 {
		repo.rxBytes = repo.rxBytes[1:]
	} else {
		repo.rxBytes = []int{}
	}

	if value == 0 {
		return value, errors.New("")
	}
	return value, nil
}
