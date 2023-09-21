package helper

import (
	"github.com/sebastianring/simulationgame"
)

func testSimulation() {
	testSC := simulationgame.SimulationConfig{
		Rows:      40,
		Cols:      100,
		Foods:     75,
		Draw:      true,
		Creature1: 20,
		Creature2: 10,
	}

	simulationgame.RunSimulation(&testSC)
}
