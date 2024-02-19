package lox

import (
	"errors"
	"fmt"
)

type ScopeEnv struct {
	parent *ScopeEnv
	vars   map[string]any
}

func NewScopeEnv(parent *ScopeEnv) *ScopeEnv {
	return &ScopeEnv{
		parent: parent,
		vars:   make(map[string]any, 0),
	}
}

func (s *ScopeEnv) Declare(name string, value any) {
	s.vars[name] = value
}

// Assigns the value to the first available name variable
// in the scope stack
func (s *ScopeEnv) Assign(name string, value any) (any, error) {
	_, ok := s.vars[name]
	if ok {
		s.vars[name] = value
		return value, nil
	} else if s.parent != nil {
		return s.parent.Assign(name, value)
	}

	return nil, errors.New(fmt.Sprintf("Var %s has never been declared", name))
}

func (s *ScopeEnv) Lookup(name string) (any, error) {
	var err error
	v, ok := s.vars[name]
	if !ok {
		if s.parent != nil {
			v, err = s.parent.Lookup(name)
		} else {
			err = errors.New(fmt.Sprintf("Var %s has never been declared", name))
		}
	}

	if err != nil {
		return nil, err
	}

	return v, nil
}
