package statemachine

import (
	"context"
)

type stateMachine struct {
	state      string
	actionChan chan func()
}

func (sm *stateMachine) Run(ctx context.Context) error {
	for {
		select {
		case f := <-sm.actionChan:
			f()
		case <-ctx.Done():
			return ctx.Err()
		}
	}

}
