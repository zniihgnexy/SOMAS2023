/* honest matrix:
update points:
	Deliberative democracy:
	voting: the non-binding vote we can see the total voting pool and the ranking of this,
	if the real vote after that came out with a  surpring result (not the same as the original)
	decrease the honest of all the agent on the bike
	pedal: get total pedaling forces(?), compare with the voting results, if no match then reduce

	Leadership:
	leadership can decide to kickout and join. the rest still can vote, but the leader make the
	decision. if the ranking and the final result not match, then tell a lie (?)
	not democrary anyway

	Dictatorship:
	the leader decide everything. if the leader choose the lootbox has different colors
	with itself, then telling a lie. if kickout a same color agent, then telling a lie.
	if kickout a same oclor then tell a lie (not a good leader)

*/

package team4

import (
	"github.com/google/uuid"
)

// HonestyRecord represents the probability of an agent being honest in a specific context.
type HonestyRecord struct {
	HonestyProbability float64
	Context            string
}

type HonestyMatrix struct {
	Records map[uuid.UUID][]HonestyRecord
}

type ActionData struct {
	// define the actions
	// votingResult, leadershipDecision, etc
}

// NewHonestyMatrix creates a new HonestyMatrix.
func NewHonestyMatrix() *HonestyMatrix {
	return &HonestyMatrix{
		Records: make(map[uuid.UUID][]HonestyRecord),
	}
}

// UpdateHonesty updates the honesty record for an agent.
func (hm *HonestyMatrix) UpdateHonesty(agentID uuid.UUID, probability float64, context string) {
	hm.Records[agentID] = append(hm.Records[agentID], HonestyRecord{HonestyProbability: probability, Context: context})
}

// GetRecords returns all honesty records for a given agent.
func (hm *HonestyMatrix) GetRecords(agentID uuid.UUID) []HonestyRecord {
	return hm.Records[agentID]
}

/*
func (hm *HonestyMatrix) CalculateHonestyBasedOnActions(agentID uuid.UUID, actionData ActionData) {
	// implement algorithm
	bikeID := hm.GetBike()
	newHonestyValue := 0.0 // fixed or changing
	fellowBikers := hm.GameState.GetMegaBikes()[bikeID].GetAgents()
	for _, agentID := range fellowBikers {
		actionData := gatherActionDataForAgent(agentID) // Implement this function based on your game logic
		honestyMatrix.CalculateHonestyBasedOnActions(agentID, actionData)
	}

	// Update the honesty matrix with the new value
	hm.UpdateHonesty(agentID, newHonestyValue, "Description of the context")
}
*/
