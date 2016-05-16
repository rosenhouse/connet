package store

import (
	"errors"
	"policy-server/models"
	"sync"

	"github.com/pivotal-golang/lager"
)

type MemoryStore struct {
	rules []models.Rule
	lock  sync.Mutex
}

func (s *MemoryStore) Add(logger lager.Logger, rule models.Rule) error {
	logger = logger.Session("memory-store-add")
	logger.Info("start")
	defer logger.Info("done")

	s.lock.Lock()
	defer s.lock.Unlock()

	s.rules = append(s.rules, rule)
	logger.Info("added", lager.Data{"rule": rule})

	return nil
}

func (s *MemoryStore) Delete(logger lager.Logger, rule models.Rule) error {
	logger = logger.Session("memory-store-delete")
	logger.Info("start")
	defer logger.Info("done")

	s.lock.Lock()
	defer s.lock.Unlock()

	newRules := []models.Rule{}

	for _, r := range s.rules {
		if !rule.Equals(r) {
			newRules = append(newRules, r)
		}
	}

	if len(newRules) == len(s.rules) {
		return errors.New("not found")
	}

	s.rules = newRules

	logger.Info("deleted", lager.Data{"rule": rule})
	return nil
}

func (s *MemoryStore) List(logger lager.Logger) ([]models.Rule, error) {
	logger = logger.Session("memory-store-list")
	logger.Info("start")
	defer logger.Info("done")

	s.lock.Lock()
	defer s.lock.Unlock()

	toReturn := make([]models.Rule, len(s.rules))
	copy(toReturn, s.rules)

	return toReturn, nil
}
