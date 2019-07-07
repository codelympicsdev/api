package challenge

import (
	"time"

	"github.com/codelympicsdev/api/database"
	"github.com/dop251/goja"
)

var timeout = 200 * time.Millisecond

// AttemptData is the generated data for an attempt
type AttemptData struct {
	Input  *database.AttemptInput
	Output *database.AttemptOutput
}

// Generate a pair response data pair
func Generate(challenge *database.Challenge) (*database.AttemptInput, *database.AttemptOutput, error) {
	vm := goja.New()
	time.AfterFunc(timeout, func() {
		vm.Interrupt(AttemptData{})
	})

	v, err := vm.RunString(challenge.Generator)
	if err != nil {
		return nil, nil, err
	}

	var data = new(AttemptData)
	err = vm.ExportTo(v, &data)
	if err != nil {
		return nil, nil, err
	}

	return data.Input, data.Output, nil
}
