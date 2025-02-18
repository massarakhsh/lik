package metric

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetrics(t *testing.T) {
	met := MetricValue{}
	for n := 0; n < 1000; n++ {
		met.SetValueFloat("", rand.Float64())
	}
	avg := met.GetValue("")
	ok := avg > 0.4 && avg < 0.6
	if !ok {
		t.Errorf("ERROR of metric.Get/SetValue")
	}
	assert.Equal(t, ok, true)
	if avg < 0.4 || avg > 0.6 {
		t.Errorf("ERROR of metric.Get/SetValue")
	}
}
