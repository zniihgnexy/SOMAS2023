package team_4

import (
	"SOMAS2023/internal/common/objects"
	"SOMAS2023/internal/common/physics"
	"SOMAS2023/internal/common/utils"
	"SOMAS2023/internal/common/voting"
	"fmt"
	"math/rand"
	"sort"

	"github.com/MattSScott/basePlatformSOMAS/messaging"
	"github.com/google/uuid"
)

type IBaselineAgent interface {
	objects.IBaseBiker
}

type BaselineAgent struct {
	*objects.BaseBiker
	currentBike          *objects.MegaBike
	OtherBiker           objects.IBaseBiker
	OtherBikerReputation float64
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
	return rank
}

// breaks code
// func (agent *BaselineAgent) DecideAllocation() voting.IdVoteMap {
// 	currentBike := agent.GetGameState().GetMegaBikes()[agent.GetBike()]
// 	rank, e := agent.rankAgentsReputation(currentBike.GetAgents())
// 	if e != nil {
// 		panic("unexpected error!")
// 	}
// 	return rank
// }

// rankTargetProposals ranks by distance and processes kickout messages to adjust rankings.
func (agent *BaselineAgent) rankTargetProposals(proposedLootBox []objects.ILootBox) (map[uuid.UUID]float64, error) {

	currentMegaBike := agent.GetGameState().GetMegaBikes()[agent.GetBike()]
	fellowBikers := currentMegaBike.GetAgents()

	// sort lootBox by distance
	sort.Slice(proposedLootBox, func(i, j int) bool {
		return physics.ComputeDistance(currentMegaBike.GetPosition(), proposedLootBox[i].GetPosition()) <
			physics.ComputeDistance(currentMegaBike.GetPosition(), proposedLootBox[j].GetPosition())
	})
	rank := make(map[uuid.UUID]float64)
	ranksum := make(map[uuid.UUID]float64)
	totalsum := float64(0)
	totaloptions := len(proposedLootBox)

	// Process messages and update honesty matrix
	msgs := agent.GetAllMessages(fellowBikers)
	for _, msg := range msgs {
		switch m := msg.(type) {
		case objects.KickOffAgentMessage:
			// Print out the ID of the agent who might be kicked off
			fmt.Printf("Received kickout message for agent ID: %s\n", m.AgentId)

			// Calculate the honesty value for the agent in the message
			//honestyValue := GetHonestyValue(m.AgentId)
			if m.AgentId == agent.GetID() {
				senderId := m.BaseMessage.GetSender()
				// Decrease the sender's honesty by 0.05
				GlobalHonestyMatrix.DecreaseHonesty(senderId.GetID(), 0.05)
			}
		}
	}

	// Original ranking logic
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

func (agent *BaselineAgent) HandleKickOffMessage(msg objects.KickOffAgentMessage) {
	if msg.AgentId == agent.GetID() {
		senderId := msg.BaseMessage.GetSender()
		GlobalHonestyMatrix.DecreaseHonesty(senderId.GetID(), 0.05)
	} else {
		senderId := msg.BaseMessage.GetSender()
		GlobalHonestyMatrix.IncreaseHonesty(senderId.GetID(), 0.05)
	}
}

func (agent *BaselineAgent) CreateKickOffMessage() []objects.KickOffAgentMessage {
	currentMegaBike := agent.GetGameState().GetMegaBikes()[agent.GetBike()]
	fellowBikers := currentMegaBike.GetAgents()

	var messages []objects.KickOffAgentMessage

	for _, fellowBiker := range fellowBikers {
		messages = append(messages, objects.KickOffAgentMessage{
			BaseMessage: messaging.CreateMessage[objects.IBaseBiker](agent, []objects.IBaseBiker{fellowBiker}),
			AgentId:     fellowBiker.GetID(),
			KickOff:     false,
		})
	}

	return messages
}

func (agent *BaselineAgent) GetAllMessages(fellowBikers []objects.IBaseBiker) []messaging.IMessage[objects.IBaseBiker] {
	var messages []messaging.IMessage[objects.IBaseBiker]

	if wantToSendMsg := true; wantToSendMsg {
		//reputationMsgs := agent.CreateReputationMessage()
		kickOffMsgs := agent.CreateKickOffMessage()
		/* 		lootboxMsg := agent.CreateLootboxMessage()
		   		joiningMsg := agent.CreateJoiningMessage()
		   		governceMsg := agent.CreateGoverenceMessage() */

		for _, msg := range kickOffMsgs {
			messages = append(messages, msg)
		}
	}

	return messages
}

// GetHonestyValue returns the honesty value for the given agent ID from the global honesty matrix.
func GetHonestyValue(agentID uuid.UUID) float64 {
	return GlobalHonestyMatrix.Records[agentID]
}
