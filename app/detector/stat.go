package detector

import "math"

type Accumulator struct {
	sum     float64
	numbers int64

	meanS float64
	qS    float64

	calculated bool
	mean       float64
	stdDev     float64
}

func (a *Accumulator) Add(v float64) {
	a.sum += v
	a.numbers++

	prevMeanS := a.meanS

	a.meanS += +((v - a.meanS) / float64(a.numbers))
	a.qS += (v - prevMeanS) * (v - a.meanS)

	a.calculated = false
}

func (a *Accumulator) MeanStdDev() (float64, float64) {
	if a.calculated {
		return a.mean, a.stdDev
	}

	a.mean = a.sum / float64(a.numbers)
	a.stdDev = math.Sqrt(a.qS / float64(a.numbers-1))

	return a.mean, a.stdDev
}

type PartitionStat struct {
	stat map[string][3]Accumulator
}

func NewPartitionStat() PartitionStat {
	return PartitionStat{
		stat: make(map[string][3]Accumulator, 128),
	}
}

func (p *PartitionStat) Add(key string, features [3]float64) {
	acc := p.stat[key]

	acc[0].Add(features[0])
	acc[1].Add(features[1])
	acc[2].Add(features[2])

	p.stat[key] = acc
}

func (p *PartitionStat) ZTest(key string, features [3]float64) [3]float64 {
	acc := p.stat[key]

	mean0, stdDev0 := acc[0].MeanStdDev()
	mean1, stdDev1 := acc[1].MeanStdDev()
	mean2, stdDev2 := acc[2].MeanStdDev()

	return [3]float64{
		(features[0] - mean0) / stdDev0,
		(features[1] - mean1) / stdDev1,
		(features[2] - mean2) / stdDev2,
	}
}
