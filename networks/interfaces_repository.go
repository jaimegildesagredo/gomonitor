package networks

import (
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

type InterfacesRepository interface {
	GetTxBytes(interfaceName string) int
	GetRxBytes(interfaceName string) int
	FindAll() []Interface
	Exists(interfaceName string) bool
}

func NewInterfacesRepository() InterfacesRepository {
	repo := interfacesRepository{
		baseDir: "/sys/class/net",
	}
	return &repo
}

type interfacesRepository struct {
	baseDir string
}

func (repo *interfacesRepository) GetTxBytes(interfaceName string) int {
	return readIntFromFile(repo.pathFor(interfaceName, "statistics", "tx_bytes"))
}

func (repo *interfacesRepository) pathFor(parts ...string) string {
	result := repo.baseDir

	for _, part := range parts {
		result += "/" + part
	}

	return result
}

func readIntFromFile(path string) int {
	var rawValue []byte
	var value int
	var err error

	rawValue, err = ioutil.ReadFile(path)

	if err != nil {
		return 0
	}

	value, err = strconv.Atoi(strings.Trim(string(rawValue), "\n"))

	if err != nil {
		return 0
	}

	return value
}

func (repo *interfacesRepository) GetRxBytes(interfaceName string) int {
	return readIntFromFile(repo.pathFor(interfaceName, "statistics", "rx_bytes"))
}

func (repo *interfacesRepository) FindAll() []Interface {
	interfaces := []Interface{}

	contents, err := ioutil.ReadDir(repo.baseDir)

	if err != nil {
		log.Println("Error getting all network interfaces", err)
		return interfaces
	}

	for _, item := range contents {
		interfaceName := item.Name()
		interfaces = append(interfaces, Interface{
			Name:  interfaceName,
			State: repo.readStateFor(interfaceName),
		})
	}

	return interfaces
}

func (repo *interfacesRepository) readStateFor(interfaceName string) string {
	content, err := ioutil.ReadFile(repo.pathFor(interfaceName, "operstate"))

	if err != nil {
		log.Println("Error getting network interface state", err)
	}

	return strings.TrimSpace(string(content))
}

func (repo *interfacesRepository) Exists(interfaceName string) bool {
	for _, interface_ := range repo.FindAll() {
		if interfaceName == interface_.Name {
			return true
		}
	}
	return false
}
