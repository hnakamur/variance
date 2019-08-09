package variance

import "math"

// VarianceMinMaxFloat64 is used to calculate the mean,
// the variance, the min, and the max value in one pass,
// without keeping all values.
//
// It uses Welford's online algorithm for calculating
// the mean and the variance.
// https://en.wikipedia.org/wiki/Algorithms_for_calculating_variance
//
// To initialize, NewVarianceMinMaxFloat64 method must be used,
// instead of the zero value of VarianceMinMaxFloat64.
//
// VarianceMinMaxFloat64 is not goroutine safe. It is caller's
// responsibility to lock VarianceFloat64 appropriately
// if it is accessed from multiple goroutines.
type VarianceMinMaxFloat64 struct {
	v   VarianceFloat64
	min float64
	max float64
}

// NewVarianceMinMaxFloat64 returns a new initialized
// VarianceMinMaxFloat64.
func NewVarianceMinMaxFloat64() VarianceMinMaxFloat64 {
	return VarianceMinMaxFloat64{
		min: math.Inf(1),
		max: math.Inf(-1),
	}
}

// Update updates the count, the mean, the m2,
// the min, and the max in the v, m2 is the sum
// of the squared distance from the mean.
func (v *VarianceMinMaxFloat64) Update(value float64) {
	v.v.Update(value)
	if value < v.min {
		v.min = value
	}
	if value > v.max {
		v.max = value
	}
}

// Count returns the count of sample values,
// that is count of Update method called.
func (v *VarianceMinMaxFloat64) Count() int {
	return v.v.Count()
}

// Mean returns the mean value.
// It returns NaN if Update was not called.
func (v *VarianceMinMaxFloat64) Mean() float64 {
	return v.v.Mean()
}

// Variance returns the variance.
// It returns NaN if Update was not called.
func (v *VarianceMinMaxFloat64) Variance() float64 {
	return v.v.Variance()
}

// SampleVariance returns the sample variance.
// It returns NaN if sample count (= count of Update method called)
// is less than two.
func (v *VarianceMinMaxFloat64) SampleVariance() float64 {
	return v.v.SampleVariance()
}

// Min returns the minimum value.
func (v *VarianceMinMaxFloat64) Min() float64 {
	return v.min
}

// Max returns the maximum value.
func (v *VarianceMinMaxFloat64) Max() float64 {
	return v.max
}
