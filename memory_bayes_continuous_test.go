package bayes

import (
	// "fmt"
	"github.com/stretchr/testify/assert"
	// "regexp"
	"testing"
)

func TestNumberPositionalBinning(t *testing.T) {
	assert.Equal(t, NumberPositionalBinning("1945")[0], "1***")
	assert.Equal(t, NumberPositionalBinning("1945")[3], "1945")
	assert.Equal(t, len(NumberPositionalBinning("")), 0)
	assert.Equal(t, NumberPositionalBinning("12")[0], "1*")
	assert.Equal(t, NumberPositionalBinning("12")[1], "12")
}

func TestNumberDiffing(t *testing.T) {
	diffing := Diff([]string{"1943", "asda", "12", "blab", "34:12:12"}, ":")
	data := []string{"1***", "19**", "194*", "1943", "1943",
		"asda", "1*", "12", "12", "blab",
		"3*******", "34******", "34:*****", "34:1****",
		"34:12***", "34:12:**", "34:12:1*", "34:12:12", "34:12:12"}
	assert.Equal(t, diffing, data)
}

func TestContinuousBayesSimple(t *testing.T) {
	b := BayesMemory{}
	splitLocal := func(data string) []string {
		return b.TokenizeSimple(data, "")
	}
	atestLocal := func(data, result string) {
		assert.Equal(t, b.Classify(splitLocal(data))[0].Category, result)
	}
  
	b.Train(splitLocal("1945"), []string{"old"})
	b.Train(splitLocal("1980"), []string{"middle"})
	b.Train(splitLocal("1990"), []string{"young"})
	b.Train(splitLocal("23"), []string{"young"})
	b.Train(splitLocal("60"), []string{"quite old"})
	b.Train(splitLocal("80"), []string{"old"})
	b.Train(splitLocal("100"), []string{"very old"})
	b.Train(splitLocal("1940"), []string{"old"})
	b.Train(splitLocal("1960"), []string{"old"})

	atestLocal("1948", "old")
	atestLocal("1985", "middle")
	atestLocal("1996", "young")
	atestLocal("67", "quite old")
	atestLocal("89", "old")
	atestLocal("1967", "old")
}

func TestContinuousBayesDuration(t *testing.T) {
	b := BayesMemory{}
	splitLocal := func(data string) []string {
		return b.TokenizeSimple(data, ":")
	}

	atestLocal := func(data, result string) {
		assert.Equal(t, b.Classify(splitLocal(data))[0].Category, result)
	}
  
	b.Train(splitLocal("An audio  12:23:34"), []string{"very long"})
	b.Train(splitLocal("An audio 03:43:45"), []string{"very long"})
	b.Train(splitLocal("00:23:56 some data"), []string{"short"})
	b.Train(splitLocal("it is duration 00:45:53"), []string{"quite long"})
	b.Train(splitLocal("00:03:12 but it works"), []string{"short"})
	b.Train(splitLocal("00:08:12"), []string{"short"})
	b.Train(splitLocal("00:12:12"), []string{"short"})
	b.Train(splitLocal("02:14:14"), []string{"very long"})
	b.Train(splitLocal("05:15:67"), []string{"very long"})
	// ========================
	atestLocal("00:05:04", "short")
	atestLocal("00:25:04", "short")
	atestLocal("00:44:23", "quite long")
	atestLocal("03:45:20", "very long")
	atestLocal("15:05:50", "very long")
	atestLocal("12:31:56", "very long")
	atestLocal("10:31:56", "very long")
	atestLocal("05:31:56", "very long")
	// TODO:: make a geospatial data
}
