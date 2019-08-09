package variance

import (
	"math"
	"math/rand"
	"strings"
	"testing"
	"time"
)

func TestVarianceMinMaxFloat64(t *testing.T) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	const eps = 1e-13
	matchFloat64 := func(a, b float64) bool {
		diffRatio := (a - b) / b
		if math.IsNaN(diffRatio) {
			return math.Abs(a) < eps
		}
		return math.Abs(diffRatio) < eps
	}

	t.Run("noValue", func(t *testing.T) {
		v := NewVarianceMinMaxFloat64()

		gotCount := v.Count()
		if gotCount != 0 {
			t.Errorf("count unmatch, got=%d, want=%d", gotCount, 0)
		}

		if !math.IsNaN(v.Mean()) {
			t.Errorf("mean must be NaN")
		}

		if !math.IsNaN(v.Variance()) {
			t.Errorf("variance must be NaN")
		}

		if !math.IsNaN(v.SampleVariance()) {
			t.Errorf("sample variance must be NaN")
		}

		if !math.IsInf(v.Min(), 1) {
			t.Errorf("min must be +Inf")
		}

		if !math.IsInf(v.Max(), -1) {
			t.Errorf("max must be -Inf")
		}
	})
	t.Run("oneValue", func(t *testing.T) {
		v := NewVarianceMinMaxFloat64()
		value := math.MaxInt64 * r.NormFloat64()
		v.Update(value)

		gotCount := v.Count()
		if gotCount != 1 {
			t.Errorf("count unmatch, got=%d, want=%d", gotCount, 1)
		}

		gotMean := v.Mean()
		if gotMean != value {
			t.Errorf("mean unmatch, got=%g, want=%g", gotMean, value)
		}

		gotVariance := v.Variance()
		if gotVariance != float64(0) {
			t.Errorf("variance unmatch, got=%g, want=%g", gotVariance, float64(0))
		}

		if !math.IsNaN(v.SampleVariance()) {
			t.Errorf("sample variance must be NaN")
		}

		gotMin := v.Min()
		if gotMin != value {
			t.Errorf("min unmatch, got=%g, want=%g", gotMin, value)
		}

		gotMax := v.Max()
		if gotMax != value {
			t.Errorf("max unmatch, got=%g, want=%g", gotMax, value)
		}
	})
	t.Run("multipleValues", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			v := NewVarianceMinMaxFloat64()

			n := r.Intn(10) + 2
			values := make([]float64, n)
			randomMax := float64(math.MaxInt64 / 2)
			var sum, min, max float64
			for i := 0; i < n; i++ {
				value := randomMax * r.NormFloat64()
				v.Update(value)
				sum += value
				values[i] = value
				if i == 0 {
					min = value
					max = value
				} else {
					if value < min {
						min = value
					}
					if max < value {
						max = value
					}
				}
			}

			mean := sum / float64(n)
			var m2 float64
			for i := 0; i < n; i++ {
				m2 += (values[i] - mean) * (values[i] - mean)
			}
			variance := m2 / float64(n)
			sampleVariance := m2 / float64(n-1)

			gotCount := v.Count()
			gotMean := v.Mean()
			gotVariance := v.Variance()
			gotSampleVariance := v.SampleVariance()
			gotMin := v.Min()
			gotMax := v.Max()

			countUnmatch := gotCount != n
			meanUnmatch := !matchFloat64(gotMean, mean)
			varianceUnmatch := !matchFloat64(gotVariance, variance)
			sampleVarianceUnmatch := !matchFloat64(gotSampleVariance, sampleVariance)
			minUnmatch := !matchFloat64(gotMin, min)
			maxUnmatch := !matchFloat64(gotMax, max)
			if countUnmatch || meanUnmatch || varianceUnmatch || sampleVarianceUnmatch ||
				minUnmatch || maxUnmatch {
				var unmatched []string
				if countUnmatch {
					unmatched = append(unmatched, "count")
				}
				if meanUnmatch {
					unmatched = append(unmatched, "mean")
				}
				if varianceUnmatch {
					unmatched = append(unmatched, "variance")
				}
				if sampleVarianceUnmatch {
					unmatched = append(unmatched, "sampleVariance")
				}
				if minUnmatch {
					unmatched = append(unmatched, "min")
				}
				if maxUnmatch {
					unmatched = append(unmatched, "max")
				}
				t.Errorf(`case %d: %s unmatched:
				values=%v,
				count          got=%d, want=%d,
				mean           got=%g, want=%g,
				variance       got=%g, want=%g,
				sampleVariance got=%g, want=%g,
				min            got=%g, want=%g,
				max            got=%g, want=%g
				`,
					i, strings.Join(unmatched, ", "),
					values,
					gotCount, n,
					gotMean, mean,
					gotVariance, variance,
					gotSampleVariance, sampleVariance,
					gotMin, min,
					gotMax, max)
			}
		}
	})
}
