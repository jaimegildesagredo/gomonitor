package loads

import (
	"time"
)

func NewLoadServiceFactory() LoadService {
	return NewLoadService(NewLoadRepository())
}

type LoadService interface {
	Monitor(delay time.Duration) chan Load
}

func NewLoadService(loadRepo LoadRepository) LoadService {
	return &loadService{
		loadRepo: loadRepo,
	}
}

type loadService struct {
	loadRepo LoadRepository
}

func (service *loadService) Monitor(delay time.Duration) chan Load {
	output := make(chan Load)
	go func() {
		for {
			output <- service.loadRepo.Get()
			time.Sleep(delay)
		}
	}()
	return output
}
