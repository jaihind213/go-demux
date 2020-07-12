package go_demux

import (
	"fmt"
	"math/rand"

	"github.com/spaolacci/murmur3"
)

const seed = 213123

var notBytesErr = fmt.Errorf("input signal not of type %T", []byte{})

func NewMurmurKeyBasedDemuxFunc(numOutputlines int) DemuxSelectorFunc {
	return func(inputSignal Signal) (int, error) {
		key, ok := inputSignal.([]byte)
		if !ok {
			return noOutputLineSelected, notBytesErr
		}
		var murmur = murmur3.New32WithSeed(seed)
		if _, err := murmur.Write(key); err != nil {
			return noOutputLineSelected, fmt.Errorf("failed to generate hash while doing demux using key: %v, err: %v", key, err)
		}
		hash := murmur.Sum32()

		outputLine := hash % uint32(numOutputlines)
		return int(outputLine), nil
	}
}

//assuming that output line ids are between 0 and 'n'
func NewRandomDemuxFunc(numOutputlines int) DemuxSelectorFunc {
	return func(ignoreMe Signal) (outputLineId int, err error) {
		return rand.Intn(numOutputlines), nil
	}
}
