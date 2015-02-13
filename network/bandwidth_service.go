package network

import "time"

type BandwidthService interface {
	MonitorBandwidth(interfaceName string, delay time.Duration) chan Bandwidth
}

func NewBandwidthServiceFactory() BandwidthService {
	return NewBandwidthService(NewBytesRepository())
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

func (service *bandwidthService) MonitorBandwidth(interfaceName string, delay time.Duration) chan Bandwidth {
	output := make(chan Bandwidth)
	go func() {
		var bandwidth Bandwidth
		var currentTxBytes int
		var previousTxBytes int
		var currentRxBytes int
		var previousRxBytes int

		for {
			bandwidth = Bandwidth{}

			currentTxBytes = service.bytesRepo.GetTx(interfaceName)
			currentRxBytes = service.bytesRepo.GetRx(interfaceName)

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
