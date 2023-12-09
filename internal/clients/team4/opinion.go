package team4

import (
	"math"

	"github.com/google/uuid"
)

// calc sigmoid
func sigmoid(x float64) float64 {
	return 1.0 / (1.0 + math.Exp(-x))
}

// get reputation value of all other agents
func (agent *BaselineAgent) GetReputation() map[uuid.UUID]float64 {
	return agent.reputation
}

// query for reputation value of specific agent with UUID
func (agent *BaselineAgent) QueryReputation(agentId uuid.UUID) float64 {
	return agent.reputation[agentId]
}

// changed version
// func (agent *BaselineAgent) CalculateReputation() {
// 	////////////////////////////
// 	//  As the program I used for debugging invoked "padal" and "break" with values of 0, I conducted tests using random numbers.
// 	// In case of an updated main program, I will need to adjust the parameters and expressions of the reputation matrix.
// 	// The current version lacks real data during the debugging process.
// 	////////////////////////////
// 	megaBikes := agent.GetGameState().GetMegaBikes()

// 	for _, bike := range megaBikes {
// 		// Get all agents on MegaBike
// 		fellowBikers := bike.GetAgents()
// 		epsilon:=1e10-12
// 		// Iterate over each agent on MegaBike, generate reputation assessment
// 		for _, otherAgent := range fellowBikers {
// 			// Exclude self
// 			selfTest := otherAgent.GetID() //nolint
// 			if selfTest == agent.GetID() {
// 				agent.reputation[otherAgent.GetID()] = 1.0
// 			}

// 			// Monitor otherAgent's location
// 			// location := otherAgent.GetLocation()
// 			// RAP := otherAgent.GetResourceAllocationParams()
// 			// fmt.Println("Agent ID:", otherAgent.GetID(), "Location:", location, "ResourceAllocationParams:", RAP)

// 			// Monitor otherAgent's forces
// 			historyenergy := agent.energyHistory[otherAgent.GetID()]
// 			lastEnergy := 1.0
// 			if len(historyenergy) > 2 {
// 				lastEnergy = historyenergy[len(historyenergy)-2]
// 				// rest of your code
// 			} else {
// 				lastEnergy = 0.0
// 			}
// 			energyLevel := otherAgent.GetEnergyLevel()
// 			consumption:= energyLevel-lastEnergy
// 			ReputationEnergy := consumption / (energyLevel+epsilon) //CAUTION: REMOVE THE RANDOM VALUE
// 			//print("我是大猴子")
// 			//fmt.Println("Agent ID:", otherAgent.GetID(), "Reputation_Forces:", ReputationEnergy)

// 			// Monitor otherAgent's bike status
// 			bikeStatus := otherAgent.GetBikeStatus()
// 			// Convert the boolean value to float64 and print the result
// 			ReputationBikeShift := 0.2
// 			if bikeStatus {
// 				ReputationBikeShift = 1.0
// 			}
// 			//fmt.Println("Agent ID:", otherAgent.GetID(), "Reputation_Bike_Shift", float64(ReputationBikeShift))

// 			// Calculate Overall_reputation
// 			OverallReputation := ReputationEnergy * ReputationBikeShift
// 			//fmt.Println("Agent ID:", otherAgent.GetID(), "Overall Reputation:", OverallReputation)

// 			// Store Overall_reputation in the reputation map
// 			agent.reputation[otherAgent.GetID()] = OverallReputation
// 		}
// 	}
// 	/* 	for agentID, agentReputation := range agent.reputation {
// 		print("Agent ID: ", agentID.String(), ", Reputation: ", agentReputation, "\n")
// 	} */

// }

// /
func (agent *BaselineAgent) CalculateReputation() {
	megaBikes := agent.GetGameState().GetMegaBikes()
	decay_factor := 0.1
	for _, bike := range megaBikes {
		fellowBikers := bike.GetAgents()
		//epsilon := 1e-10

		for _, otherAgent := range fellowBikers {
			selfTest := otherAgent.GetID()
			if selfTest == agent.GetID() {
				agent.reputation[otherAgent.GetID()] = 1.0
				continue // Skip the rest of the loop for the current agent
			}

			historyenergy := agent.energyHistory[otherAgent.GetID()]
			lastEnergy := 1.0
			if len(historyenergy) > 2 {
				lastEnergy = historyenergy[len(historyenergy)-2]
			} else {
				lastEnergy = 0.0
			}
			energyLevel := otherAgent.GetEnergyLevel()
			consumption := energyLevel - lastEnergy

			myhistoryenergy := agent.energyHistory[agent.GetID()]
			mylastEnergy := 1.0
			if len(myhistoryenergy) > 2 {
				mylastEnergy = myhistoryenergy[len(myhistoryenergy)-2]
			} else {
				mylastEnergy = 0.0
			}
			myenergyLevel := agent.GetEnergyLevel()
			myconsumption := myenergyLevel - mylastEnergy
			EnergyReputation := (consumption / (energyLevel + 0.001)) - (myconsumption / (myenergyLevel + 0.001))

			//consumption / (energyLevel + epsilon)

			// Check if ReputationEnergy is NaN or Inf before proceeding

			bikeStatus := otherAgent.GetBikeStatus()
			ReputationBikeShift := 0.2
			if bikeStatus {
				ReputationBikeShift = 1.0
			}

			OverallReputation := EnergyReputation * ReputationBikeShift

			// Check if OverallReputation is NaN or Inf before storing
			if math.IsNaN(OverallReputation) || math.IsInf(OverallReputation, 0) {
				agent.reputation[otherAgent.GetID()] = 0.0
				continue // Skip the rest of the loop for the current agent
			}
			OverallReputation = sigmoid(OverallReputation)
			// print((1-decay_factor)*(agent.reputation[otherAgent.GetID()])+decay_factor*(OverallReputation))
			// print("\n")
			agent.reputation[otherAgent.GetID()] = (1-decay_factor)*(agent.reputation[otherAgent.GetID()]) + decay_factor*(OverallReputation)
		}
	}
}

///
/* // Reputation and Honesty Matrix Teams Must Implement these or similar functions

func (agent *BaselineAgent) CalculateReputation() {
	////////////////////////////
	//  As the program I used for debugging invoked "padal" and "break" with values of 0, I conducted tests using random numbers.
	// In case of an updated main program, I will need to adjust the parameters and expressions of the reputation matrix.
	// The current version lacks real data during the debugging process.
	////////////////////////////
	megaBikes := agent.GetGameState().GetMegaBikes()

	for _, bike := range megaBikes {
		// Get all agents on MegaBike
		fellowBikers := bike.GetAgents()

		// Iterate over each agent on MegaBike, generate reputation assessment
		for _, otherAgent := range fellowBikers {
			// Exclude self
			selfTest := otherAgent.GetID() //nolint
			if selfTest == agent.GetID() {
				agent.reputation[otherAgent.GetID()] = 1.0
			}

			// Monitor otherAgent's location
			// location := otherAgent.GetLocation()
			// RAP := otherAgent.GetResourceAllocationParams()
			// fmt.Println("Agent ID:", otherAgent.GetID(), "Location:", location, "ResourceAllocationParams:", RAP)

			// Monitor otherAgent's forces
			forces := otherAgent.GetForces()
			energyLevel := otherAgent.GetEnergyLevel()
			ReputationForces := float64(forces.Pedal+forces.Brake+rand.Float64()) / energyLevel //CAUTION: REMOVE THE RANDOM VALUE
			// fmt.Println("Agent ID:", otherAgent.GetID(), "Reputation_Forces:", ReputationForces)

			// Monitor otherAgent's bike status
			bikeStatus := otherAgent.GetBikeStatus()
			// Convert the boolean value to float64 and print the result
			ReputationBikeShift := 0.2
			if bikeStatus {
				ReputationBikeShift = 1.0
			}
			// fmt.Println("Agent ID:", otherAgent.GetID(), "Reputation_Bike_Shift", float64(ReputationBikeShift))

			// Calculate Overall_reputation
			OverallReputation := ReputationForces * ReputationBikeShift
			// fmt.Println("Agent ID:", otherAgent.GetID(), "Overall Reputation:", OverallReputation)

			// Store Overall_reputation in the reputation map
			agent.reputation[otherAgent.GetID()] = OverallReputation
		}
	}
	// for agentID, agentReputation := range agent.reputation {
	// 	print("Agent ID: ", agentID.String(), ", Reputation: ", agentReputation, "\n")
	// }
} */

func (agent *BaselineAgent) CalculateHonestyMatrix() {
	for _, bike := range agent.GetGameState().GetMegaBikes() {
		for _, biker := range bike.GetAgents() {
			bikerID := biker.GetID()
			
			agent.honestyMatrix[bikerID] = 1.0
		}
	}
}

func (agent *BaselineAgent) IncreaseHonesty(agentID uuid.UUID, increaseAmount float64) {
	agent.honestyMatrix[agentID] += increaseAmount
}

func (agent *BaselineAgent) DecreaseHonestyHonesty(agentID uuid.UUID, decreaseAmount float64) {
	agent.honestyMatrix[agentID] -= decreaseAmount
}
