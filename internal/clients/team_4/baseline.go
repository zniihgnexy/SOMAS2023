package team_4

import (
	"SOMAS2023/internal/common/objects"
	"SOMAS2023/internal/common/physics"
	"SOMAS2023/internal/common/utils"
	"SOMAS2023/internal/common/voting"
	"fmt"

	//"go/printer"
	"math"
	"math/rand"
	"sort"
	"time"

	"github.com/google/uuid"
)

type IBaselineAgent interface {
	objects.IBaseBiker
}

type BaselineAgent struct {
	objects.BaseBiker
	currentBike *objects.MegaBike
}

// DecideAction only pedal
func (agent *BaselineAgent) DecideAction() objects.BikerAction {
	return objects.Pedal

}

// DecideForces randomly based on current energyLevel
func (agent *BaselineAgent) DecideForces(direction uuid.UUID) {
	energyLevel := agent.GetEnergyLevel() // 当前能量

	randomBreakForce := float64(0)
	randomPedalForce := rand.Float64() * energyLevel // 使用 rand 包生成随机的 pedal 力量，可以根据需要调整范围

	if randomPedalForce == 0 {
		// just random break force based on energy level, but not too much
		randomBreakForce += rand.Float64() * energyLevel * 0.5
	} else {
		randomBreakForce = 0
	}
	fmt.Print("energyLevel", energyLevel)
	fmt.Println("randomPedalForce", randomPedalForce)

	// 因为force是一个struct,包括pedal, brake,和turning，因此需要一起定义，不能够只有pedal
	forces := utils.Forces{
		Pedal: randomPedalForce,
		Brake: randomBreakForce, // random for now
		Turning: utils.TurningDecision{
			SteerBike:     true,
			SteeringForce: physics.ComputeOrientation(agent.GetLocation(), agent.GetGameState().GetMegaBikes()[direction].GetPosition()) - agent.GetGameState().GetMegaBikes()[agent.currentBike.GetID()].GetOrientation(),
		},
	}

	agent.SetForces(forces)
}

// DecideJoining accept all
func (agent *BaselineAgent) DecideJoining(pendingAgents []uuid.UUID) map[uuid.UUID]bool {
	decision := make(map[uuid.UUID]bool)
	for _, agent := range pendingAgents {
		decision[agent] = true
	}
	return decision
}

func (agent *BaselineAgent) FinalDirectionVote(proposals []uuid.UUID) voting.LootboxVoteMap {
	boxesInMap := agent.GetGameState().GetLootBoxes()
	boxProposed := make([]objects.ILootBox, len(proposals))
	for i, pp := range proposals {
		boxProposed[i] = boxesInMap[pp]
	}
	rank, e := agent.rankTargetProposals(boxProposed)
	if e != nil {
		panic("unexpected error!")
	}
	agent.calc_Reputation_value(0.1, 0.5, 0.8)
	//agent.Reputation_loop()
	return rank
}

func (agent *BaselineAgent) rankTargetProposals(proposedLootBox []objects.ILootBox) (map[uuid.UUID]float64, error) {

	currentMegaBike := agent.GetGameState().GetMegaBikes()[agent.GetBike()]
	fellowBikers := currentMegaBike.GetAgents()

	// sort lootBox by distance
	sort.Slice(proposedLootBox, func(i, j int) bool {
		return physics.ComputeDistance(currentMegaBike.GetPosition(), proposedLootBox[i].GetPosition()) < physics.ComputeDistance(currentMegaBike.GetPosition(), proposedLootBox[j].GetPosition())
	})
	rank := make(map[uuid.UUID]float64)
	ranksum := make(map[uuid.UUID]float64)
	totalsum := float64(0)
	totaloptions := len(proposedLootBox)

	for i, lootBox := range proposedLootBox {
		for _, fellows := range fellowBikers {
			if fellows.GetColour() == lootBox.GetColour() {
				ranksum[lootBox.GetID()] += float64(totaloptions-i) / 1.2
				totalsum += float64(totaloptions-i) / 1.2
			}
		}
		if lootBox.GetColour() == agent.GetColour() {
			ranksum[lootBox.GetID()] += 2 * float64(totaloptions-i)
			totalsum += 2 * float64(totaloptions-i)
		} else {
			ranksum[lootBox.GetID()] += float64(totaloptions-i) / float64(totaloptions)
			totalsum += float64(totaloptions-i) / float64(totaloptions)
		}
	}
	for _, lootBox := range proposedLootBox {
		rank[lootBox.GetID()] = ranksum[lootBox.GetID()] / totalsum
	}

	return rank, nil
}

// rankAgentReputation randomly rank agents
func (agent *BaselineAgent) rankAgentsReputation(agentsOnBike []objects.IBaseBiker) (map[uuid.UUID]float64, error) {
	rank := make(map[uuid.UUID]float64)
	for i, agent := range agentsOnBike {
		//getReputationMatrix()
		//choose the highest one
		rank[agent.GetID()] = float64(i)
	}
	return rank, nil
}

func (agent *BaselineAgent) DecideGovernance() voting.GovernanceVote {
	rank := make(map[utils.Governance]float64)
	rank[utils.Democracy] = 1
	rank[utils.Leadership] = 0
	rank[utils.Dictatorship] = 0
	rank[utils.Invalid] = 0
	//for i := utils.Democracy; i <= utils.Invalid; i++ {
	//  rank[i] = 0.25
	//}
	fmt.Println(rank)
	return rank
}

// ///////////////////////////////////////////////
// Yanzhou's  Changes (the function is called in line 82)
// ----> 82 agent.calc_Reputation_value(0.1, 0.5, 0.8)
//
// Description:
// 1. it returns value of the reputation of a agent.
// 2. it have to be tuned and normalized in the reputation matrix
// 3. now, we can calculate the value of the reputation value, but cannot calculate the matrix because datatype issue.
//
// ISSUE:
// 1. The 'calc_Reputation_value' function can only calculate values of "*BaselineAgent", can not work for "*IBaseBiker"
// 2. To generate the reputation, we can enumerate all agents in the system, and calculate reputations for each of them
//    But as I enumerate all agents, (in line 196) the varible otherAgent is "*IBaseBiker",
//    it can only be called in the Gameloop.go, I have no access to it
//    Thus, I don't think there is a way to calculate the reputation matrix in this file
// ///////////////////////////////////////////////

// /////////////////////////
// func (agent *BaselineAgent) GetReputation() map[uuid.UUID]float64 {
// 	reputation := make(map[uuid.UUID]float64)
// 	print("1232132213123212312\n")
// 	return reputation
// }

func (agent *BaselineAgent) leaving() float64 { //get input
	rand.Seed(time.Now().UnixNano())
	randomBit := float64(rand.Intn(2))
	return randomBit
}

func (agent *BaselineAgent) energy_consumption() float64 {
	// Generate a random energy consumption for the current loop
	consumption := rand.Float64() * agent.GetEnergyLevel()
	return consumption
}

func (agent *BaselineAgent) selected_Lootbox_distance() float64 {
	// Randomly select a lootbox as the selected lootbox for the current loop
	lootBoxes := agent.GetGameState().GetLootBoxes()

	var distance float64
	for _, lootBox := range lootBoxes {
		// Implement specific logic here to calculate the distance between agent and LootBox
		distance = physics.ComputeDistance(agent.GetLocation(), lootBox.GetPosition())
		if rand.Float64() < 0.5 {
			break
		}
	}
	return distance
}

// func (agent *BaselineAgent) Reputation_loop() map[uuid.UUID]float64 {
// 	reputation := make(map[uuid.UUID]float64)
// 	megaBikes := agent.GetGameState().GetMegaBikes()

// 	for _, bike := range megaBikes {
// 		// Get all agents on MegaBike
// 		fellowBikers := bike.GetAgents()

// 		// Iterate over each agent on MegaBike, generate reputation assessment
// 		for _, otherAgent := range fellowBikers {
// 			// Exclude self
// 			test := otherAgent.GetID() //nolint
// 			if test == agent.GetID() {
// 				continue
// 			}
// 			otherAgent.GetReputation()

// 		}
// 	}
// 	return reputation
// }

func (agent *BaselineAgent) calc_Reputation_value(b, c, d float64) float64 {
	reputation_value := 0.0
	// Get all LootBoxes on the field
	lootBoxes := agent.GetGameState().GetLootBoxes()

	// 1. Minimum distance and selected target value between agent and all LootBoxes
	var minDistance float64 = math.Inf(1)
	i := 0
	for _, lootBox := range lootBoxes {
		// Implement specific logic here to calculate the distance between agent and LootBox
		distance := physics.ComputeDistance(agent.GetLocation(), lootBox.GetPosition())

		// Check if it is a new minimum distance
		if distance < minDistance {
			minDistance = distance
		}

		// If the current loop index equals the specified index a, record the distance to the selected LootBox
		i = i + 1
	}

	// Calculate distance_Reputation
	distance_Reputation := 1 - (b * agent.selected_Lootbox_distance() / minDistance)

	// Here you can use minDistance, selected_Lootdistance, and distance_Reputation for further logic
	print("Agent " + agent.GetID().String() + " Reputation = \n")
	print("distance_Reputation: \t", distance_Reputation)
	print("\n")

	// 2. EnergyLevel of the agent in this and the previous loop
	// Get the current EnergyLevel
	currentEnergyLevel := agent.GetEnergyLevel()
	currentConsumption := agent.energy_consumption()

	energy_Reputation := currentEnergyLevel - c*currentConsumption

	print("energy_Reputation: \t", energy_Reputation)
	print("\n")

	// Get whether the agent leaves the bike this round 0 means no, 1 means yes
	leave_value := agent.leaving()
	Leave_Reputation := 1 - d*leave_value
	print("Leave_Reputation: \t", Leave_Reputation)
	print("\n")

	///
	///////////////
	reputation_value = Leave_Reputation*energy_Reputation - distance_Reputation
	print("reputation_value: \t", reputation_value)
	print("\n")
	print("\n")
	return reputation_value
}
