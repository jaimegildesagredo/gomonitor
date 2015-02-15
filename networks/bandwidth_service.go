package networks

import (
	"fmt"
	"time"
)

type BandwidthService interface {
	MonitorBandwidth(interfaceName string, delay time.Duration) (chan Bandwidth, error)
}

func NewBandwidthServiceFactory() BandwidthService {
	return NewBandwidthService(NewInterfacesRepository())
}

func NewBandwidthService(interfacesRepo InterfacesRepository) BandwidthService {
	service := bandwidthService{
		interfacesRepo: interfacesRepo,
	}
	return &service
}

type bandwidthService struct {
	interfacesRepo InterfacesRepository
}

type Bandwidth struct {
	Up        int
	Down      int
	CreatedAt time.Time
}

func (service *bandwidthService) MonitorBandwidth(interfaceName string, delay time.Duration) (chan Bandwidth, error) {
	output := make(chan Bandwidth)
	txBytes := []int{}
	rxBytes := []int{}

	if !service.interfacesRepo.Exists(interfaceName) {
		return output, fmt.Errorf("Interface %s does not exist", interfaceName)
	}

	go func() {
		var bandwidth Bandwidth

		for {
			txBytes = append(txBytes, service.interfacesRepo.GetTxBytes(interfaceName))
			if len(txBytes) > 2 {
				txBytes = txBytes[1:]
			}

			rxBytes = append(rxBytes, service.interfacesRepo.GetRxBytes(interfaceName))
			if len(rxBytes) > 2 {
				rxBytes = rxBytes[1:]
			}

			bandwidth = Bandwidth{}
			bandwidth.CreatedAt = time.Now().UTC()

			if len(txBytes) == 2 {
				bandwidth.Up = int(float64(txBytes[1]-txBytes[0]) / delay.Seconds())
			}

			if len(rxBytes) == 2 {
				bandwidth.Down = int(float64(rxBytes[1]-rxBytes[0]) / delay.Seconds())
			}

			output <- bandwidth

			time.Sleep(delay)
		}
	}()

	return output, nil
}
