package usecase

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestParseAIResponse(t *testing.T) {
	s := &TourMaker{}

	tests := []struct {
		input  string
		output []string
	}{
		{
			input:  "sdfsdf ['история', 'культура', 'достопримечательности']sdfsdd",
			output: []string{"история", "культура", "достопримечательности"},
		},
		{
			input:  "sdfsdf ['история', 'культура', 'достопри",
			output: []string{},
		},
		{
			input:  "sdfsdf [\"история', 'культура', 'достопримечательности']",
			output: []string{"история", "культура", "достопримечательности"},
		},
	}

	for _, test := range tests {
		out := s.parseAIPreferences(test.input)

		require.Equal(t, test.output, out)
	}

}

func TestFindAINum(t *testing.T) {
	s := &TourMaker{}

	require.Equal(t, 10, s.findNum("fhfgfg  fgh10dfghdfhg"))
	require.Equal(t, 5, s.findNum("fhfgfg  fgh5dfghdfhg"))
	require.Equal(t, 123, s.findNum("123"))
	require.Equal(t, 0, s.findNum("fhfgfg  fghdfghdfhg"))
}
