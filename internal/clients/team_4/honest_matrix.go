package team_4

import (
	"github.com/google/uuid"
)

type HonestyMatrix struct {
	Records map[uuid.UUID]float64 // Make sure to export Records
}

// GlobalHonestyMatrix holds the honesty values for all agents.
// It needs to be initialized before usage.
var GlobalHonestyMatrix *HonestyMatrix

// NewHonestyMatrix creates a new HonestyMatrix with default honesty values for each agent.
func NewHonestyMatrix(agentIDs []uuid.UUID) *HonestyMatrix {
	hm := &HonestyMatrix{
		Records: make(map[uuid.UUID]float64),
	}
	for _, agentID := range agentIDs {
		hm.Records[agentID] = 1.0
	}
	return hm
}

// UpdateHonesty
func (hm *HonestyMatrix) UpdateHonesty(agentID uuid.UUID, newHonestyValue float64) {
	hm.Records[agentID] = newHonestyValue
}

// GetHonesty
func (hm *HonestyMatrix) GetHonesty(agentID uuid.UUID) float64 {
	return hm.Records[agentID]
}

func (hm *HonestyMatrix) DecreaseHonesty(agentID uuid.UUID, decreaseAmount float64) {
	if currentHonesty, ok := hm.Records[agentID]; ok {
		newHonesty := currentHonesty - decreaseAmount
		if newHonesty < 0 {
			newHonesty = 0.0000001
		}
		hm.Records[agentID] = newHonesty
	}
}

func (hm *HonestyMatrix) IncreaseHonesty(agentID uuid.UUID, decreaseAmount float64) {
	if currentHonesty, ok := hm.Records[agentID]; ok {
		newHonesty := currentHonesty + decreaseAmount
		if newHonesty > 1 {
			newHonesty = 1
		}
		hm.Records[agentID] = newHonesty
	}
}
