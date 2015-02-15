package networks_test

import (
	"testing"
	"time"

	. "github.com/jaimegildesagredo/gomonitor/networks"
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
	interfacesRepository := newInMemoryInterfacesRepository(
		[]string{AN_INTERFACE_NAME},
		[]int{FIRST_TX_BYTES, LAST_TX_BYTES},
		[]int{FIRST_RX_BYTES, LAST_RX_BYTES})

	bandwidthService := NewBandwidthService(interfacesRepository)
	bandwidths, _ := bandwidthService.MonitorBandwidth(AN_INTERFACE_NAME, A_DELAY)

	bandwidth := <-bandwidths
	bandwidth = <-bandwidths

	if bandwidth.Up != EXPECTED_BANDWIDTH_UP {
		t.Fatal("Invalid bandwidth up value", bandwidth.Up, "expected", EXPECTED_BANDWIDTH_UP)
	}

	if bandwidth.Down != EXPECTED_BANDWIDTH_DOWN {
		t.Fatal("Invalid bandwidth down value", bandwidth.Down, "expected", EXPECTED_BANDWIDTH_DOWN)
	}
}

func TestMonitorBandwidthWhenInterfaceDoesNotExists(t *testing.T) {
	interfacesRepository := newInMemoryInterfacesRepository([]string{}, []int{}, []int{})
	bandwidthService := NewBandwidthService(interfacesRepository)

	_, err := bandwidthService.MonitorBandwidth(AN_INTERFACE_NAME, A_DELAY)

	if err == nil {
		t.Fatal("Expected an error when interface does not exists")
	}
}

func newInMemoryInterfacesRepository(interfaces []string, txBytes []int, rxBytes []int) InterfacesRepository {
	repository := inMemoryInterfacesRepository{
		interfaces: interfaces,
		txBytes:    txBytes,
		rxBytes:    rxBytes,
	}
	return &repository
}

type inMemoryInterfacesRepository struct {
	interfaces []string
	txBytes    []int
	rxBytes    []int
}

func (repo *inMemoryInterfacesRepository) GetTxBytes(interfaceName string) int {
	var value int
	if len(repo.txBytes) > 0 {
		value = repo.txBytes[0]
	}

	if len(repo.txBytes) > 1 {
		repo.txBytes = repo.txBytes[1:]
	} else {
		repo.txBytes = []int{}
	}

	return value
}

func (repo *inMemoryInterfacesRepository) GetRxBytes(interfaceName string) int {
	var value int
	if len(repo.rxBytes) > 0 {
		value = repo.rxBytes[0]
	}

	if len(repo.rxBytes) > 1 {
		repo.rxBytes = repo.rxBytes[1:]
	} else {
		repo.rxBytes = []int{}
	}

	return value
}

func (repo *inMemoryInterfacesRepository) GetAllInterfaces() []string {
	return repo.interfaces
}

func (repo *inMemoryInterfacesRepository) Exists(name string) bool {
	for _, interfaceName := range repo.interfaces {
		if name == interfaceName {
			return true
		}
	}
	return false
}
