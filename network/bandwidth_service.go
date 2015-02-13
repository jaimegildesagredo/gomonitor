package network

import "time"

type BandwidthService interface {
	MonitorBandwidth(interfaceName string, delay time.Duration) (chan Bandwidth, error)
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

func (service *bandwidthService) MonitorBandwidth(interfaceName string, delay time.Duration) (chan Bandwidth, error) {
	output := make(chan Bandwidth)

	var currentTxBytes int
	var previousTxBytes int
	var currentRxBytes int
	var previousRxBytes int
	var err error

	currentTxBytes, err = service.bytesRepo.GetTx(interfaceName)

	if err != nil {
		return output, err
	}

	currentRxBytes, err = service.bytesRepo.GetRx(interfaceName)

	if err != nil {
		return output, err
	}

	go func() {
		var bandwidth Bandwidth

		for {
			time.Sleep(delay)

			previousTxBytes = currentTxBytes
			previousRxBytes = currentRxBytes
			currentTxBytes, _ = service.bytesRepo.GetTx(interfaceName)
			currentRxBytes, _ = service.bytesRepo.GetRx(interfaceName)

			bandwidth = Bandwidth{}

			if previousTxBytes != currentTxBytes && previousRxBytes != currentRxBytes {
				bandwidth.Up = int(float64(currentTxBytes-previousTxBytes) / delay.Seconds())
				bandwidth.Down = int(float64(currentRxBytes-previousRxBytes) / delay.Seconds())
				output <- bandwidth
			}
		}
	}()

	return output, nil
}
