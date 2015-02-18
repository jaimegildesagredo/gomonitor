package loads

import (
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

type LoadRepository interface {
	Get() Load
}

func NewLoadRepository() LoadRepository {
	return &loadRepository{}
}

type loadRepository struct{}

func (repo *loadRepository) Get() Load {
	raw, err := ioutil.ReadFile("/proc/loadavg")

	if err != nil {
		log.Println("Error getting system load", err)
		return Load{}
	}

	values := strings.Split(string(raw), " ")

	return Load{parseValue(values[0]), parseValue(values[1]), parseValue(values[2])}
}

func parseValue(raw string) float32 {
	result, err := strconv.ParseFloat(raw, 32)

	if err != nil {
		log.Println("Error parsing float", err)
	}

	return float32(result)
}
