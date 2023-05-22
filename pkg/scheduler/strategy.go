package scheduler

import (
	"math"
	"math/rand"
)

type Strategy interface {
	//SetMaxQuota(maxProcessedReqs int, maxRate float64, maxActives int)
	SetMaxActives(maxActives int64)
	SetMaxProcessedReqs(maxProcessedReqs int64)
	SetMaxRate(maxRate float64)

	ExpandOrShrink(incomingRegsItv int64, processedReqsItv int64, numActives int64) int64
}

type Gradienter struct {
	quotaList *CircularList

	maxActives       int64
	maxProcessedReqs int64
	maxRate          float64
}

const (
	NoNeedUpdating = iota
	NumActivesNeedsExpansion
	NumActivesNeedsShrinking
	MaxProcessedReqsNeedsExpansion
	MaxProcessedReqsNeedsShrinking
	MaxRateNeedsExpansion
	MaxRateNeedsShrinking
)

const (
	MaxDefaultGoRoutines    int64   = 500 * 10000
	MaxDefaultProcessedReqs int64   = int64(^uint(0) >> 1)
	MaxDefaultRate          float64 = 1000
)

func NewGradienter() *Gradienter {
	return &Gradienter{
		//quotaList:        NewCircularList(),
		maxActives:       MaxDefaultGoRoutines,
		maxProcessedReqs: MaxDefaultProcessedReqs,
		maxRate:          MaxDefaultRate,
	}
}

//already locked at sheduler
func (g *Gradienter) SetMaxActives(maxActives int64) {
	if maxActives == -1 {
		g.maxActives = MaxDefaultGoRoutines
		return
	}
	g.maxActives = maxActives
}

func (g *Gradienter) SetMaxProcessedReqs(maxProcessedReqs int64) {
	if maxProcessedReqs == -1 {
		g.maxProcessedReqs = MaxDefaultProcessedReqs
		return
	}
	g.maxProcessedReqs = maxProcessedReqs
}

func (g *Gradienter) SetMaxRate(maxRate float64) {
	if maxRate >= 0 && maxRate <= MaxDefaultRate {
		g.maxRate = maxRate
		return
	}
	g.maxRate = MaxDefaultRate
}

func (g *Gradienter) expand(ir int64, pr int64, numActives int64) int64 {
	weight := float64(ir) / (float64(ir) + float64(g.maxActives))
	reminder := float64(g.maxActives - numActives)
	randfloat := rand.Float64()
	uncasted := randfloat * reminder * weight
	casted := int64(uncasted)
	return casted
}

func (g *Gradienter) ExpandOrShrink(ir int64, pr int64, numActives int64) (diff int64) {
	flag := g.needToUpdate(ir, pr, numActives)

	switch flag {
	case NoNeedUpdating:
		return 0

	case NumActivesNeedsExpansion:
		return g.expand(ir, pr, numActives)

	case NumActivesNeedsShrinking:
		//directly shrink to maxActives
		return g.maxActives - numActives

	case MaxProcessedReqsNeedsExpansion:
		//map the processedReqs expanding section to goroutines expanding section
		//the goroutines expanding section should be: [0, maxActives - numActives]
		//the processedReqs expanding section shoulde be: [0, maxProcessedReqs - pr]
		//so the shrinking number should be mapped as: number / (maxProcessedReqs - pr) / (maxActives - numActives)
		//any error exists please point out
		return g.expand(ir, pr, numActives)

	case MaxProcessedReqsNeedsShrinking:
		//shrink 20%, should've add some weight like others
		shrinks := float64(numActives) * 0.2 * -1
		return int64(math.Floor(shrinks))

	case MaxRateNeedsExpansion:
		return g.expand(ir, pr, numActives)

	case MaxRateNeedsShrinking:
		//shrink 20%
		shrinks := float64(numActives) * 0.2 * -1
		return int64(math.Floor(shrinks))
	}
	return 0
}

//
//no updating first, then shrinking, expansion last
func (g *Gradienter) needToUpdate(ir int64, pr int64, numActives int64) int {
	flag := NoNeedUpdating
	if float64(g.maxActives) > float64(numActives)*(1+0.1) {
		flag = NumActivesNeedsExpansion
	} else if float64(g.maxActives) < float64(numActives)*(1-0.1) {
		flag = NumActivesNeedsShrinking
	}

	if float64(g.maxProcessedReqs) > float64(pr)*(1+0.1) && flag == NoNeedUpdating {
		flag = MaxProcessedReqsNeedsExpansion
	} else if float64(g.maxProcessedReqs) < float64(pr)*(1-0.1) && (flag == NoNeedUpdating || flag == NumActivesNeedsExpansion) {
		flag = MaxProcessedReqsNeedsShrinking
	}

	if ir == 0 && numActives > 0 {
		//means no incoming requests, needs to shrink
		flag = MaxRateNeedsShrinking
		return flag
	}

	if g.maxRate > float64(pr)/float64(ir)*(1+0.05) && flag == NoNeedUpdating {
		//over speed, need to shrink
		flag = MaxRateNeedsExpansion
	} else if g.maxRate < float64(pr)/float64(ir)*(1-0.05) && (flag == NoNeedUpdating || flag == NumActivesNeedsExpansion || flag == MaxProcessedReqsNeedsExpansion) {
		flag = MaxRateNeedsShrinking
	}
	return flag
}
