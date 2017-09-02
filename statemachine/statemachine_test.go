package statemachine

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Run(t *testing.T) {
	sm := stateMachine{
		state:      "initial",
		actionChan: make(chan func()),
	}
	ctx := context.Background()
	receive := make(chan string)
	go sm.Run(ctx)
	go func() {
		sm.actionChan <- func() { receive <- "action" }
		sm.actionChan <- func() { receive <- "action2" }
	}()
	assert.Equal(t, "action", <-receive)
	assert.Equal(t, "action2", <-receive)
	ctx.Done()
}

func Test_ErrorContext(t *testing.T) {
	sm := stateMachine{
		state:      "initial",
		actionChan: make(chan func()),
	}
	ctx := context.Background()
	receive := make(chan string)
	go sm.Run(ctx)
	go func() {
		sm.actionChan <- func() { receive <- "action3" }
		ctx.Done()

	}()
	assert.Equal(t, "action3", <-receive)

}
