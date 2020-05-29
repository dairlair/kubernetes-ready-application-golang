package core

import (
	"github.com/stepsisters/kgb/mocks/pkg/core"
	"testing"
)

func TestRun_WillExecuteComponentRun(t *testing.T) {
	component := &mocks.ComponentInterface{}
	component.On("Run").Return(func() {}, func() error {return nil}, nil)

	Run(map[string]ComponentInterface{"test": component})

	component.AssertCalled(t, "Run")
	component.AssertExpectations(t)
}