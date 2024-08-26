package command

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDiceCommand(t *testing.T) {
	t.Parallel()
	assert.NotNil(t, NewDiceCommand())
}

func Test_Command(t *testing.T) {
	t.Parallel()
	assert.NotNil(t, (&DiceCommand{}).Command())
}
