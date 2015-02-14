package networks

import (
	"log"
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
	Up   int
	Down int
}

func (service *bandwidthService) MonitorBandwidth(interfaceName string, delay time.Duration) (chan Bandwidth, error) {
	output := make(chan Bandwidth)

	var currentTxBytes int
	var previousTxBytes int
	var currentRxBytes int
	var previousRxBytes int
	var err error

	currentTxBytes, err = service.interfacesRepo.GetTxBytes(interfaceName)

	if err != nil {
		return output, err
	}

	currentRxBytes, err = service.interfacesRepo.GetRxBytes(interfaceName)

	if err != nil {
		return output, err
	}

	go func() {
		var bandwidth Bandwidth

		for {
			time.Sleep(delay)

			previousTxBytes = currentTxBytes
			previousRxBytes = currentRxBytes
			currentTxBytes, _ = service.interfacesRepo.GetTxBytes(interfaceName)
			currentRxBytes, _ = service.interfacesRepo.GetRxBytes(interfaceName)

			bandwidth = Bandwidth{}

			if previousTxBytes != currentTxBytes && previousRxBytes != currentRxBytes {
				bandwidth.Up = int(float64(currentTxBytes-previousTxBytes) / delay.Seconds())
				bandwidth.Down = int(float64(currentRxBytes-previousRxBytes) / delay.Seconds())

				log.Println("Bandwidth for", interfaceName, bandwidth)
				output <- bandwidth
			}
		}
	}()

	return output, nil
}
