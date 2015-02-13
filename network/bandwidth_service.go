package network

import "time"

type BandwidthService interface {
	MonitorBandwidth(delay time.Duration) chan Bandwidth
}

func NewBandwidthServiceFactory(interfaceName string) BandwidthService {
	return NewBandwidthService(NewBytesRepository(interfaceName))
}

func NewBandwidthService(bytesRepo BytesRepository) BandwidthService {
	service := bandwidthService{
		bytesRepo: bytesRepo,
	}
	return &service
}

type bandwidthService struct {
	bytesRepo BytesRepository
}

type Bandwidth struct {
	Up   int
	Down int
}

func (service *bandwidthService) MonitorBandwidth(delay time.Duration) chan Bandwidth {
	output := make(chan Bandwidth)
	go func() {
		var bandwidth Bandwidth
		var currentTxBytes int
		var previousTxBytes int
		var currentRxBytes int
		var previousRxBytes int

		for {
			bandwidth = Bandwidth{}

			currentTxBytes = service.bytesRepo.GetTx()
			currentRxBytes = service.bytesRepo.GetRx()

			if previousTxBytes != 0 && previousRxBytes != 0 {
				bandwidth.Up = int(float64(currentTxBytes-previousTxBytes) / delay.Seconds())
				bandwidth.Down = int(float64(currentRxBytes-previousRxBytes) / delay.Seconds())
				output <- bandwidth
			}

			previousTxBytes = currentTxBytes
			previousRxBytes = currentRxBytes

			time.Sleep(delay)
		}
	}()
	return output
}
