package api

import (
	// "encoding/json"
	"errors"
	"fmt"
	sg "github.com/sebastianring/simulationgame"
)

type RoundDataType byte

const (
	AliveAtEnd RoundDataType = iota
	Killed
	Spawned
)

type simpleRoundData struct {
	ID              int
	CreatureSummary map[sg.BoardObjectType]*sg.CreatureSummary
}

func getRoundData(b *sg.Board, datatype RoundDataType) ([]*simpleRoundData, error) {
	fmt.Println("Starting to get round data for a specific board.")
	compiledRounds := []*simpleRoundData{}

	switch datatype {
	case AliveAtEnd:
		for _, val := range b.Rounds {
			srd := simpleRoundData{
				ID:              val.Id,
				CreatureSummary: val.CreaturesAliveAtEndSum,
			}
			compiledRounds = append(compiledRounds, &srd)
		}

	case Killed:
		for _, val := range b.Rounds {
			srd := simpleRoundData{
				ID:              val.Id,
				CreatureSummary: val.CreaturesKilledSum,
			}
			compiledRounds = append(compiledRounds, &srd)
		}

	case Spawned:
		for _, val := range b.Rounds {
			srd := simpleRoundData{
				ID:              val.Id,
				CreatureSummary: val.CreaturesSpawnedSum,
			}
			compiledRounds = append(compiledRounds, &srd)
		}

	default:
		fmt.Println("Error when trying to get round data - data type does not exist.")
		return nil, errors.New("Datatype can't be found")
	}

	return compiledRounds, nil
}
