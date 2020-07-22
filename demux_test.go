package go_demux

import (
	"fmt"
	"testing"
)

func TestChannelDemux_UsingMurmur(t *testing.T) {

	const numChnls = 3

	////prepare test
	outputChannels := make(map[OutputLineId]chan Signal, numChnls)
	for i := 0; i < numChnls; i++ {
		outputChannels[OutputLineId(i)] = make(chan Signal, 5)
	}

	//set up
	dmux := ChannelDemux{
		OutputChannels: outputChannels,
		Selector:       NewMurmurKeyBasedDemuxFunc(numChnls),
	}

	////test cases
	type test struct {
		input        []byte
		lineExpected int
		errExpected  bool
	}

	tests := []test{
		{[]byte("1647fec8408-39c40000010a443b"), 0, false},
		{[]byte("1647fec8408-39c40000010a443b"), 0, false},
		{[]byte("1647ee3d0a6-58030000010f610d"), 0, false},
		{[]byte("1647ea7b83a-5ec60000010e2cbf"), 1, false},
	}

	//prepare expected results from test input above.
	xpectedCountPerOutputLine := make(map[OutputLineId]int) // key is output line id , value is count
	for _, tc := range tests {
		if tc.errExpected {
			continue
		}
		if cnt, ok := xpectedCountPerOutputLine[OutputLineId(tc.lineExpected)]; !ok {
			xpectedCountPerOutputLine[OutputLineId(tc.lineExpected)] = 1
		} else {
			xpectedCountPerOutputLine[OutputLineId(tc.lineExpected)] = cnt + 1
		}
	}

	////do test
	for _, tc := range tests {
		got, err := dmux.Demux(tc.input)
		if !tc.errExpected && err != nil {
			t.Errorf("did not expect error for tc %v bot got %v", tc.input, err)
		}
		if err != nil {
			fmt.Println(err)
		}
		if got != tc.lineExpected {
			t.Errorf("expected output region chanl %v, got %v", tc.lineExpected, got)
		}
	}

	//check num of signals received per channel
	for outputLineId, ch := range outputChannels {
		if xpected, ok := xpectedCountPerOutputLine[outputLineId]; !ok {
			//but if
			if len(outputChannels[outputLineId]) > 0 {
				t.Errorf("expected 0 elements for outputline %v but got %v", outputLineId, len(outputChannels[outputLineId]))
			}
		} else {
			if xpected != len(ch) {
				t.Errorf("expected %v elements in output line for outputline  %v, got %v", xpected, outputLineId, len(ch))
			}
		}
	}
}

func TestChannelDemux_WithNoChannls(t *testing.T) {
	dmux := ChannelDemux{
		OutputChannels: nil,
		Selector:       NewRandomDemuxFunc(2),
	}
	id, err := dmux.Demux("abc")
	if err == nil {
		t.Errorf("expected error (but got none) while demuxing with 0 output lines but got none")
	}
	if id >= 0 {
		t.Errorf("no output lines were set, but demux chose a valid output line with id %v", id)
	}
}

func TestChannelDemux_WithInvalidChannel(t *testing.T) {

	outputChannels := make(map[OutputLineId]chan Signal, 1)
	outputChannels[0] = make(chan Signal, 1)

	dmux := ChannelDemux{
		OutputChannels: outputChannels,
		Selector: func(inputSignal Signal) (outputLineId int, err error) {
			return -2, nil //return some id not in output channels
		},
	}
	id, err := dmux.Demux("abc")
	if err == nil {
		t.Errorf("expected error (but got none) while demuxing , as configured selector returned a channel id does not exist")
	}
	if id >= 0 {
		t.Errorf("configured selector always returns an invalid output line id, but we got a valid one : %v", id)
	}
}

func TestChannelDemux_WithNoSelector(t *testing.T) {

	outputChannels := make(map[OutputLineId]chan Signal, 1)
	outputChannels[0] = make(chan Signal, 1)

	dmux := ChannelDemux{
		OutputChannels: outputChannels,
		Selector:       nil,
	}
	id, err := dmux.Demux("abc")
	if err == nil {
		t.Errorf("expected error as selector function was not set, but gone.")
	}
	if id >= 0 {
		t.Errorf("selector not set but we got a valid output line id : %v", id)
	}
}

type blackHole struct {
}

func (b blackHole) Add(signal interface{}) {
	fmt.Println("blackhole got", signal)
}

func ExampleGenericDemux() {
	numLines := 2
	outputLines := make(map[OutputLineId]Adder, numLines)
	for i := 0; i < numLines; i++ {
		outputLines[OutputLineId(i)] = blackHole{}
	}

	dmux := GenericDemux{
		OutputLines: outputLines,
		Selector:    NewRandomDemuxFunc(numLines),
	}

	_, _ = dmux.Demux("abc")
	// Output: blackhole got abc
}
