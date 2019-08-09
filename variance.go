package variance

import "math"

// VarianceFloat64 is used to calculate the mean value
// and the variance value in one pass, without keeping all values.
// It uses Welford's online algorithm.
// https://en.wikipedia.org/wiki/Algorithms_for_calculating_variance
//
// VarianceFloat64 is not goroutine safe. It is caller's
// responsibility to lock VarianceFloat64 appropriately
// if it is accessed from multiple goroutines.
type VarianceFloat64 struct {
	count int
	mean  float64
	m2    float64
}

// Update updates the count, the mean, and the m2
// in the v, m2 is the sum of the squared distance
// from the mean.
func (v *VarianceFloat64) Update(value float64) {
	v.count++
	delta := value - v.mean
	v.mean += delta / float64(v.count)
	delta2 := value - v.mean
	v.m2 += delta * delta2
}

// Count returns the count of sample values,
// that is count of Update method called.
func (v *VarianceFloat64) Count() int {
	return v.count
}

// Mean returns the mean value.
// It returns NaN if Update was not called.
func (v *VarianceFloat64) Mean() float64 {
	if v.count == 0 {
		return math.NaN()
	}
	return v.mean
}

// Variance returns the variance.
// It returns NaN if Update was not called.
func (v *VarianceFloat64) Variance() float64 {
	return v.m2 / float64(v.count)
}

// SampleVariance returns the sample variance.
// It returns NaN if sample count (= count of Update method called)
// is less than two.
func (v *VarianceFloat64) SampleVariance() float64 {
	if v.count == 0 {
		return math.NaN()
	}
	return v.m2 / float64(v.count-1)
}
