// parser_test
package main

import (
	"testing"

	"github.com/bmizerany/assert"
)

func Test_CheckCalledNumber(t *testing.T) {
	assert.Equal(t, checkCalledNumber("01088822200"), true, "01088822200")
	assert.Equal(t, checkCalledNumber("010-88822200"), true, "010-88822200")
	assert.Equal(t, checkCalledNumber("057427614632"), true, "057427614632")
	assert.Equal(t, checkCalledNumber("0574-27614632"), true, "0574-27614632")
	assert.Equal(t, checkCalledNumber("13587010754"), true, "13587010754")

	assert.Equal(t, checkCalledNumber("750140112"), false, "750140112")
	assert.Equal(t, checkCalledNumber("804003801613"), false, "804003801613")
}
