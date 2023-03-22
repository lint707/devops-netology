package resources

import (
	"backend/pkg/logging"
	"fmt"
	"math/rand"
)

var _ Service = &service{}

type service struct {
	logger  logging.Logger
	storage map[string]*resource
}

type Service interface {
	getAllocatedName(name string) *resource
	putAllocatedName(name, resourceType, region string) *resource
	generateName(resourceType, region string) string
	info(msg string)
}

func NewService(l *logging.Logger) (Service, error) {
	return &service{
		logger:  *l,
		storage: make(map[string]*resource),
	}, nil
}

func (s service) getAllocatedName(name string) *resource {
	if v, ok := s.storage[name]; ok {
		return v
	}
	return nil
}

func (s service) putAllocatedName(name, resourceType, region string) *resource {
	resource := &resource{
		Name:   name,
		Type:   resourceType,
		Region: region,
	}
	s.storage[name] = resource
	return resource
}

func (s service) generateName(resourceType, region string) string {
	newName := fmt.Sprint("N", rand.Intn(10000), "T", resourceType, "R", region[0:2])
	if _, exists := s.storage[newName]; exists {
		return s.generateName(resourceType, region)
	} else {
		return newName
	}
}

func (s service) info(msg string) {
	s.logger.Info(msg)
}
