package main

import (
	"errors"
	"math"
	"sort"
)

type statistics struct {
	numbers           []float64
	sum               float64
	mean              float64
	median            float64
	standardDeviation float64
	modes             []float64
	precision         int
	err               error
}

func NewStatistics(precision int) *statistics {
	return &statistics{
		precision: precision,
	}
}

func (s *statistics) Compute(inputs []float64) {
	s.numbers = inputs
	if !s.validInputs() {
		s.err = errors.New("Can't compute mean of empty inputs.")
		return
	}

	sort.Float64s(s.numbers)
	s.mean, _ = s.computeMean()
	s.median, _ = s.computeMedian()
	s.standardDeviation, _ = s.computeStandardDeviation()
	s.modes, _ = s.computeModes()
}

func (s *statistics) computeMean() (mean float64, err error) {
	if !s.validInputs() {
		return mean, errors.New("Can't compute mean of empty inputs.")
	}

	var sumErr error
	s.sum, sumErr = s.computeSum()
	if sumErr != nil {
		return mean, sumErr
	}

	mean = s.sum / float64(len(s.numbers))
	return s.roundToPrecision(mean), nil
}

func (s *statistics) computeSum() (sum float64, err error) {
	if !s.validInputs() {
		return sum, errors.New("Can't compute mean of empty inputs.")
	}

	for _, x := range s.numbers {
		sum += x
	}
	return s.roundToPrecision(sum), nil
}

func (s *statistics) computeMedian() (median float64, err error) {
	if !s.validInputs() {
		return median, errors.New("Can't compute mean of empty inputs.")
	}

	middle := len(s.numbers) / 2
	median = s.numbers[middle]
	if len(s.numbers)%2 == 0 {
		median = (median + s.numbers[middle-1]) / 2
	}
	return s.roundToPrecision(median), nil
}

func (s *statistics) computeStandardDeviation() (sd float64, err error) {
	if !s.validInputs() {
		return sd, errors.New("Can't compute mean of empty inputs.")
	}

	mean, meanErr := s.computeMean()
	if meanErr != nil {
		return sd, meanErr
	}

	var sum float64
	for _, number := range s.numbers {
		sum += math.Pow(number-mean, 2)
	}

	result := math.Sqrt(sum / float64((len(s.numbers) - 1)))
	return s.roundToPrecision(result), nil
}

func (s *statistics) computeModes() (modes []float64, err error) {
	if !s.validInputs() {
		return modes, errors.New("Can't compute mean of empty inputs.")
	}

	maxOccurrence := math.MinInt64
	occurrences := make(map[float64]int)
	for _, input := range s.numbers {
		if _, ok := occurrences[input]; ok {
			occurrences[input] += 1
		} else {
			occurrences[input] = 1
		}

		if occurrences[input] > maxOccurrence {
			maxOccurrence = occurrences[input]
		}
	}

	for input, occurrence := range occurrences {
		if occurrence == maxOccurrence {
			modes = append(modes, input)
		}
	}

	return
}

func (s *statistics) validInputs() bool {
	return len(s.numbers) > 0
}

func (s *statistics) roundToPrecision(input float64) float64 {
	multiplier := math.Pow(10, float64(s.precision))
	return (float64(int(input * multiplier))) / multiplier
}
