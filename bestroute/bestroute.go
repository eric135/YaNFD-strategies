/* YaNFD-strategies - Yet another NDN Forwarding Daemon forwarding strategies
 *
 * Copyright (C) 2021 Eric Newberry.
 *
 * This file is licensed under the terms of the MIT License, as found in LICENSE.md.
 */

package main

import (
	"strconv"

	"github.com/eric135/YaNFD/core"
	"github.com/eric135/YaNFD/fw"
	"github.com/eric135/YaNFD/ndn"
	"github.com/eric135/YaNFD/table"
)

// BestRoute is a forwarding strategy that forwards Interests to the nexthop with the lowest cost.
type BestRoute struct {
	fib  *table.FibStrategyEntry
	name *ndn.Name
	fw.StrategyBase
}

// New creates a new instance of the BestRoute strategy.
func New(fib *table.FibStrategyEntry) *BestRoute {
	bestRoute := new(BestRoute)
	bestRoute.fib = fib
	bestRoute.name, _ = ndn.NameFromString(fw.StrategyPrefix + "/best-route/%FD%01")
	return bestRoute
}

func (s *BestRoute) String() string {
	return "Strategy-BestRoute"
}

// GetName ...
func (s *BestRoute) GetName() *ndn.Name {
	return s.name
}

// AfterContentStoreHit ...
func (s *BestRoute) AfterContentStoreHit(pitEntry *table.PitEntry, inFace int, data *ndn.Data) {
	// Send downstream
	core.LogTrace(s, "Forwarding content store hit Data "+data.Name().String()+" to "+strconv.Itoa(inFace))
	s.SendData(data, pitEntry, inFace)
}

// AfterReceiveData ...
func (s *BestRoute) AfterReceiveData(pitEntry *table.PitEntry, inFace int, data *ndn.Data) {
	for faceID := range pitEntry.InRecords {
		core.LogTrace(s, "Forwarding Data "+data.Name().String()+" to "+strconv.Itoa(faceID))
		s.SendData(data, pitEntry, faceID)
	}
}

// AfterReceiveInterest ...
func (s *BestRoute) AfterReceiveInterest(pitEntry *table.PitEntry, inFace int, interest *ndn.Interest) {
	nexthops := s.fib.LongestPrefixNexthops(interest.Name())
	if len(nexthops) == 0 {
		core.LogDebug(s, "No nexthop for Interest "+interest.Name().String()+" - DROP")
		return
	}

	lowestCost := nexthops[0]
	for _, nexthop := range nexthops {
		if nexthop.Cost < lowestCost.Cost {
			lowestCost = nexthop
		}
	}

	core.LogTrace(s, "Forwarding Interest "+interest.Name().String()+" to "+strconv.Itoa(lowestCost.Nexthop))
	s.SendInterest(interest, lowestCost.Nexthop)
}

// BeforeSatisfyInterest ...
func (s *BestRoute) BeforeSatisfyInterest(pitEntry *table.PitEntry, inFace int, interest *ndn.Interest) {
	// Does nothing in BestRoute
}
