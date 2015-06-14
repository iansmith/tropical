package std

import (
	"fmt"

	"github.com/iansmith/tropical"
)

type TreeManipulator struct {
	parent   tropical.Interactor
	children []tropical.Interactor
}

//
// Tree Manipulator defaults to having any number of children
//
func NewTreeManipulator(parent tropical.Interactor) tropical.TreeManipulator {
	return &TreeManipulator{
		parent:   parent,
		children: []tropical.Interactor{},
	}
}

func (t *TreeManipulator) Children() []tropical.Interactor {
	return t.children
}

func (t *TreeManipulator) AppendChild(i tropical.Interactor) []tropical.Interactor {
	found := false
	for _, child := range t.children {
		if child == i {
			found = true
			break
		}
	}
	if !found {
		t.children = append(t.children, i)
	}
	return t.children
}

func (t *TreeManipulator) Parent() tropical.Interactor {
	return t.parent
}

//
// SingleChild limits children to 1
//

type SingleChild struct {
	parent   tropical.Interactor
	children []tropical.Interactor
}

func NewSingleChild(parent tropical.Interactor) tropical.TreeManipulator {
	return &SingleChild{
		parent:   parent,
		children: []tropical.Interactor{},
	}
}

func (s *SingleChild) Children() []tropical.Interactor {
	return s.children
}

func (s *SingleChild) AppendChild(i tropical.Interactor) []tropical.Interactor {
	found := false
	if len(s.children) > 0 && s.children[0] == i {
		found = true
	}
	if !found && len(s.children) > 0 {
		panic(fmt.Sprintf("%T only accepts one child! have %T and trying to add %T", s, s.children[0], i))
	}
	if !found {
		s.children = append(s.children, i)
	}
	return s.children
}

func (s *SingleChild) Parent() tropical.Interactor {
	return s.parent
}
