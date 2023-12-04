package team_4

import (
	"SOMAS2023/internal/common/objects"
	"SOMAS2023/internal/common/physics"
	"SOMAS2023/internal/common/utils"
	"SOMAS2023/internal/common/voting"
	"fmt"
	"math"
	"math/rand"
	"sort"

	"github.com/google/uuid"
)

type IBaselineAgent interface {
	objects.IBaseBiker

	//MISSING FUNCTIONS
	CalculateReputation( /*choose*/ ) map[uuid.UUID]float64    //calculate reputation matrix
	CalculateHonestyMatrix( /*choose*/ ) map[uuid.UUID]float64 //calculate honesty matrix

	//INCOMPLETE/NO STRATEGY FUNCTIONS
	DecideAction() objects.BikerAction                         //determines what action the agent is going to take this round. (changeBike or Pedal)
	DecideForce(direction uuid.UUID)                           //defines the vector you pass to the bike: [pedal, brake, turning]
	DecideJoining(pendinAgents []uuid.UUID) map[uuid.UUID]bool //decide whether to accept or not accept bikers, ranks the ones
	ChangeBike() uuid.UUID                                     //called when biker wants to change bike, it will choose which bike to try and join
	VoteForKickout() map[uuid.UUID]int

	//CURRENTLY UNUSED/NOT CONSIDERED FUNCTIONS

	VoteDictator() voting.IdVoteMap
	VoteLeader() voting.IdVoteMap
	DictateDirection() uuid.UUID //called only when the agent is the dictator
	LeadDirection() uuid.UUID    //called only when the agent is the leader

	//IMPLEMENTED FUNCTIONS
	ProposeDirection() uuid.UUID                                    //returns the id of the desired lootbox
	FinalDirectionVote(proposals []uuid.UUID) voting.LootboxVoteMap //returns rank of proposed lootboxes
	DecideGovernance() voting.GovernanceVote                        //decide the governance system
	//some parts commented - awaiting further implementation of reputation and honesty matrix
	DecideAllocation() voting.IdVoteMap //decide the allocation parameters

	//HELPER FUNCTIONS
	UpdateDecisionData() //updates all the data needed for the decision making process(call at the start of any decision making function)

	rankFellowsReputation(agentsOnBike []objects.IBaseBiker) (map[uuid.UUID]float64, error) //returns normal rank of fellow bikers reputation
	rankFellowsHonesty(agentsOnBike []objects.IBaseBiker) (map[uuid.UUID]float64, error)    //returns normal rank of fellow bikers honesty

	rankTargetProposals(proposedLootBox []objects.ILootBox) (map[uuid.UUID]float64, error) //returns ranking of the proposed lootboxes

}
type BaselineAgent struct {
	objects.BaseBiker
	currentBike       *objects.MegaBike
	lootBoxColour     utils.Colour
	proposedLootBox   objects.ILootBox
	mylocationHistory []utils.Coordinates     //log location history for this agent
	energyHistory     map[uuid.UUID][]float64 //log energy level for all agents
	reputation        map[uuid.UUID]float64   //record reputation for other agents, 0-1
	honestyMatrix     map[uuid.UUID]float64   //record honesty for other agents, 0-1
}

func (agent *BaselineAgent) UpdateDecisionData() {
	//Initialize mapping if not initialized yet (= nil)
	if agent.energyHistory == nil {
		agent.energyHistory = make(map[uuid.UUID][]float64)
	}
	if len(agent.mylocationHistory) == 0 {
		agent.mylocationHistory = make([]utils.Coordinates, 0)
	}
	fmt.Println("")
	fmt.Println("Updating decision data ...")
	//update location history for the agent
	agent.mylocationHistory = append(agent.mylocationHistory, agent.GetLocation())
	//get current bike
	currentBike := agent.GetGameState().GetMegaBikes()[agent.GetBike()]
	//get fellow bikers
	fellowBikers := currentBike.GetAgents()
	//update energy history for each fellow biker
	for _, fellow := range fellowBikers {
		fellowID := fellow.GetID()
		currentEnergyLevel := fellow.GetEnergyLevel()
		//Append bikers current energy level to the biker's history
		agent.energyHistory[fellowID] = append(agent.energyHistory[fellowID], currentEnergyLevel)
	}
	//call reputation and honesty matrix to calcuiate/update them
	//save updated reputation and honesty matrix
	agent.reputation = agent.CalculateReputation()
	agent.honestyMatrix = agent.CalculateHonestyMatrix()
}

func (agent *BaselineAgent) rankFellowsReputation(agentsOnBike []objects.IBaseBiker) (map[uuid.UUID]float64, error) {
	totalsum := float64(0)
	rank := make(map[uuid.UUID]float64)

	for _, fellow := range agentsOnBike {
		fellowID := fellow.GetID()
		totalsum += agent.reputation[fellowID]
	}
	//normalize the reputation
	for _, fellow := range agentsOnBike {
		fellowID := fellow.GetID()
		rank[fellowID] = float64(agent.reputation[fellowID] / totalsum)
	}
	return rank, nil
}

func (agent *BaselineAgent) rankFellowsHonesty(agentsOnBike []objects.IBaseBiker) (map[uuid.UUID]float64, error) {
	totalsum := float64(0)
	rank := make(map[uuid.UUID]float64)

	for _, fellow := range agentsOnBike {
		fellowID := fellow.GetID()
		totalsum += agent.honestyMatrix[fellowID]
	}
	//normalize the honesty
	for _, fellow := range agentsOnBike {
		fellowID := fellow.GetID()
		rank[fellowID] = float64(agent.honestyMatrix[fellowID] / totalsum)
	}
	return rank, nil
}

func (agent *BaselineAgent) rankTargetProposals(proposedLootBox []objects.ILootBox) (map[uuid.UUID]float64, error) {
	rank := make(map[uuid.UUID]float64)
	ranksum := make(map[uuid.UUID]float64)
	totalsum := float64(0)
	totaloptions := len(proposedLootBox)
	minEnergyThreshold := 0.2 //if energy level is below this threshold, the agent will increase voting towards its colour lootbox

	currentBike := agent.GetGameState().GetMegaBikes()[agent.GetBike()]
	fellowBikers := currentBike.GetAgents()
	reputationRank, e1 := agent.rankFellowsReputation(fellowBikers)
	honestyRank, e2 := agent.rankFellowsHonesty(fellowBikers)
	if e1 != nil || e2 != nil {
		panic("unexpected error!")
	}

	//sort Proposed Loot Boxes by distance
	sort.Slice(proposedLootBox, func(i, j int) bool {
		return physics.ComputeDistance(currentBike.GetPosition(), proposedLootBox[i].GetPosition()) < physics.ComputeDistance(currentBike.GetPosition(), proposedLootBox[j].GetPosition())
	})

	for i, lootBox := range proposedLootBox {
		//loop through all fellow bikers and check if they have the same colour as the lootbox
		for _, fellow := range fellowBikers {
			fellowID := fellow.GetID()
			if fellow.GetColour() == lootBox.GetColour() {
				weight := (float64(totaloptions-i) * reputationRank[fellowID] * honestyRank[fellowID]) / 1.2
				ranksum[lootBox.GetID()] += weight
				totalsum += weight
			}
		}

		if lootBox.GetColour() == agent.GetColour() {
			weight := float64(totaloptions-i) * 2.0
			//if energy level is below threshold, increase weighting towards own colour lootbox
			if agent.GetEnergyLevel() < minEnergyThreshold {
				weight *= 2
			}
			ranksum[lootBox.GetID()] += weight
			totalsum += weight

		} else {
			weight := float64(totaloptions-i) / float64(totaloptions)
			ranksum[lootBox.GetID()] += weight
			totalsum += weight
		}
	}
	for _, lootBox := range proposedLootBox {
		rank[lootBox.GetID()] = ranksum[lootBox.GetID()] / totalsum
	}

	return rank, nil
}

func (agent *BaselineAgent) FinalDirectionVote(proposals []uuid.UUID) voting.LootboxVoteMap {
	agent.UpdateDecisionData()
	boxesInMap := agent.GetGameState().GetLootBoxes()
	boxProposed := make([]objects.ILootBox, len(proposals))
	for i, pp := range proposals {
		boxProposed[i] = boxesInMap[pp]
	}
	rank, e := agent.rankTargetProposals(boxProposed)
	if e != nil {
		panic("unexpected error!")
	}
	return rank
}

func (agent *BaselineAgent) DecideAllocation() voting.IdVoteMap {
	agent.UpdateDecisionData()
	distribution := make(map[uuid.UUID]float64) //make(voting.IdVoteMap)
	currentBike := agent.GetGameState().GetMegaBikes()[agent.GetBike()]
	fellowBikers := currentBike.GetAgents()
	totalEnergySpent := float64(0)
	totalAllocation := float64(0)

	// reputationRank, e1 := agent.rankFellowsReputation(fellowBikers)
	// honestyRank, e2 := agent.rankFellowsHonesty(fellowBikers)
	// if e1 != nil || e2 != nil {
	// 	panic("unexpected error!")
	// }

	for _, fellow := range fellowBikers {
		// w1 := 0.0 //weight for reputation
		// w2 := 0.0 //weight for honesty
		// w3 := 0.5 //weight for energy spent
		// w4 := 0.5 //weight for energy level
		fellowID := fellow.GetID()
		energyLog := agent.energyHistory[fellowID]
		energySpent := energyLog[len(energyLog)-2] - energyLog[len(energyLog)-1]
		totalEnergySpent += energySpent
		// In the case where my fellow biker is the same colour as the lootbox
		if fellow.GetColour() == agent.lootBoxColour {
			// w1 = 0.0
			// w2 = 0.0
			// w3 = 1.0
			// w4 = 1.0
			// In the case where the I am the same colour as the lootbox
			if fellow.GetColour() == agent.GetColour() {
				// w1 = 0.0
				// w2 = 0.0
				// w3 = 1.0
				// w4 = 1.0
			}
		}
		// distribution[fellow.GetID()] = (w1 * reputationRank[fellowID]) + (w2 * honestyRank[fellowID]) + (w3 * energySpent) + (w4 * fellow.GetEnergyLevel()))
		distribution[fellow.GetID()] = energySpent * rand.Float64() // random for now
		totalAllocation += distribution[fellow.GetID()]
	}

	//normalize the distribution
	for _, fellow := range fellowBikers {
		fellowID := fellow.GetID()
		distribution[fellowID] = distribution[fellowID] / totalAllocation
	}

	return distribution
}

// Reputation and Honesty Matrix Teams Must Implement these or similar functions

func (agent *BaselineAgent) CalculateReputation( /*choose*/ ) map[uuid.UUID]float64 {
	random := make(map[uuid.UUID]float64)
	return random
}
func (agent *BaselineAgent) CalculateHonestyMatrix( /*choose*/ ) map[uuid.UUID]float64 {
	random := make(map[uuid.UUID]float64)
	return random
}

func (agent *BaselineAgent) DisplayFellowsEnergyHistory() {
	currentBike := agent.GetGameState().GetMegaBikes()[agent.GetBike()]
	fellowBikers := currentBike.GetAgents()
	for _, fellow := range fellowBikers {
		fellowID := fellow.GetID()
		fmt.Println("")
		fmt.Println("Energy history for: ", fellowID)
		fmt.Print(agent.energyHistory[fellowID])
		fmt.Println("")
	}
}

func (agent *BaselineAgent) ProposeDirection() uuid.UUID {
	agent.proposedLootBox = nil
	lootBoxes := agent.GetGameState().GetLootBoxes()
	agentLocation := agent.GetLocation() //agent's location
	shortestDistance := math.MaxFloat64

	for _, lootbox := range lootBoxes {
		lootboxLocation := lootbox.GetPosition()
		distance := physics.ComputeDistance(agentLocation, lootboxLocation)
		if agent.proposedLootBox == nil && distance < shortestDistance {
			shortestDistance = distance
			agent.proposedLootBox = lootbox
		}
		if distance < shortestDistance || agent.GetColour() == lootbox.GetColour() {
			shortestDistance = distance
			agent.proposedLootBox = lootbox
		}
	}
	return agent.proposedLootBox.GetID()
}

// DecideJoining accept all
func (agent *BaselineAgent) DecideJoining(pendingAgents []uuid.UUID) map[uuid.UUID]bool {
	decision := make(map[uuid.UUID]bool)
	for _, agent := range pendingAgents {
		decision[agent] = true
	}
	return decision
}

func (agent *BaselineAgent) DecideGovernance() voting.GovernanceVote {
	my_energy := agent.GetEnergyLevel()
	governanceRanking := make(voting.GovernanceVote)
	// decide energy thresholds are arbitary
	// can add another threshold for deciding governance if required
	if my_energy <= 2 {
		governanceRanking[utils.Democracy] = 1.0
		governanceRanking[utils.Dictatorship] = 0.0
		governanceRanking[utils.Leadership] = 0.0
	} else if my_energy > 2 {
		governanceRanking[utils.Democracy] = 0.0
		governanceRanking[utils.Dictatorship] = 0.0
		governanceRanking[utils.Leadership] = 1.0
	} else {
		governanceRanking[utils.Democracy] = 1.0
		governanceRanking[utils.Dictatorship] = 0.0
		governanceRanking[utils.Leadership] = 0.0
	}
	//fmt.Println(governanceRanking)
	return governanceRanking
}

func (agent *BaselineAgent) DecideAction() objects.BikerAction {
	currentBike := agent.GetGameState().GetMegaBikes()[agent.GetBike()]
	fellowBikers := currentBike.GetAgents()
	myColor := agent.GetColour()
	// Initialize count of fellow bikers with the same color
	// Reputation and honesty
	// Leadership + Dictatorhsip
	// Resource allocation
	count := 0

	for _, fellow := range fellowBikers {
		colorAgent := fellow.GetColour()

		// Check if the color of the fellow biker matches the agent's color
		if myColor == colorAgent {
			count++
		}
	}
	// Initialize action variable
	var action objects.BikerAction
	// Check the count and set the action accordingly
	if count < 2 {
		action = objects.ChangeBike
	} else {
		action = objects.Pedal
	}
	// Return the action
	return action
}

func (agent *BaselineAgent) ChangeBike() uuid.UUID {
	megaBikes := agent.GetGameState().GetMegaBikes()
	i, targetI := 0, rand.Intn(len(megaBikes))
	// Go doesn't have a sensible way to do this...
	for id := range megaBikes {
		if i == targetI {
			return id
		}
		i++
	}
	panic("no bikes")
}
