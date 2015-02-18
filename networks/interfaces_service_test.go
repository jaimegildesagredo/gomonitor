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
	ANOTHER_INTERFACE_NAME  = "another-interface-name"
)

func TestFindAllReturnsSliceOfNetworkInterfaceNames(t *testing.T) {
	interfacesRepository := newInMemoryInterfacesRepository(
		someInterfaces(),
		[]int{FIRST_TX_BYTES, LAST_TX_BYTES},
		[]int{FIRST_RX_BYTES, LAST_RX_BYTES})

	interfacesService := NewInterfacesService(interfacesRepository)

	interfaces := interfacesService.FindAll()

	expectedInterfaces := []string{AN_INTERFACE_NAME, ANOTHER_INTERFACE_NAME}
	if !equal(interfaces, expectedInterfaces) {
		t.Fatal("Expected", interfaces, "to equal", expectedInterfaces)
	}

}

func someInterfaces() []Interface {
	interfaces := []Interface{}
	for _, name := range []string{AN_INTERFACE_NAME, ANOTHER_INTERFACE_NAME} {
		interfaces = append(interfaces, Interface{Name: name})
	}
	return interfaces
}

func equal(actual []Interface, expected []string) bool {
	if len(actual) != len(expected) {
		return false
	}

	for i, item := range actual {
		if item.Name != expected[i] {
			return false
		}
	}

	return true
}

func TestMonitorBandwidth(t *testing.T) {
	interfacesRepository := newInMemoryInterfacesRepository(
		someInterfaces(),
		[]int{FIRST_TX_BYTES, LAST_TX_BYTES},
		[]int{FIRST_RX_BYTES, LAST_RX_BYTES})

	interfacesService := NewInterfacesService(interfacesRepository)
	bandwidths, _ := interfacesService.MonitorBandwidth(AN_INTERFACE_NAME, A_DELAY)

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
	interfacesRepository := newInMemoryInterfacesRepository([]Interface{}, []int{}, []int{})
	interfacesService := NewInterfacesService(interfacesRepository)

	_, err := interfacesService.MonitorBandwidth(AN_INTERFACE_NAME, A_DELAY)

	if err == nil {
		t.Fatal("Expected an error when interface does not exists")
	}
}

func newInMemoryInterfacesRepository(interfaces []Interface, txBytes []int, rxBytes []int) InterfacesRepository {
	repository := inMemoryInterfacesRepository{
		interfaces: interfaces,
		txBytes:    txBytes,
		rxBytes:    rxBytes,
	}
	return &repository
}

type inMemoryInterfacesRepository struct {
	interfaces []Interface
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

func (repo *inMemoryInterfacesRepository) FindAll() []Interface {
	return repo.interfaces
}

func (repo *inMemoryInterfacesRepository) Exists(name string) bool {
	for _, interface_ := range repo.interfaces {
		if name == interface_.Name {
			return true
		}
	}
	return false
}
