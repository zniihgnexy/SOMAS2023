package server

import (
	"SOMAS2023/internal/common/objects"
	"SOMAS2023/internal/common/utils"

	team4 "SOMAS2023/internal/clients/team_4"

	baseserver "github.com/MattSScott/basePlatformSOMAS/BaseServer"
	"github.com/google/uuid"
)

const BikerAgentCount = 10

func GetAgentGenerators() []baseserver.AgentGeneratorCountPair[objects.IBaseBiker] {
	return []baseserver.AgentGeneratorCountPair[objects.IBaseBiker]{
		baseserver.MakeAgentGeneratorCountPair[objects.IBaseBiker](BikerAgentGenerator, BikerAgentCount),
	}
}

func BikerAgentGenerator() objects.IBaseBiker {
	// Create a new instance of BaselineAgent
	baselineAgent := &team4.BaselineAgent{
		BaseBiker: *objects.GetBaseBiker(utils.GenerateRandomColour(), uuid.New()),
	}
	// Return the baselineAgent as an IBaseBiker
	return baselineAgent
}

func (s *Server) spawnLootBox() {
	lootBox := objects.GetLootBox()
	s.lootBoxes[lootBox.GetID()] = lootBox
}

func (s *Server) replenishLootBoxes() {
	count := LootBoxCount - len(s.lootBoxes)
	for i := 0; i < count; i++ {
		s.spawnLootBox()
	}
}

func (s *Server) spawnMegaBike() {
	megaBike := objects.GetMegaBike()
	s.megaBikes[megaBike.GetID()] = megaBike
}

func (s *Server) replenishMegaBikes() {
	for i := 0; i < MegaBikeCount-len(s.megaBikes); i++ {
		s.spawnMegaBike()
	}
}
