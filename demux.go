package go_demux

import "fmt"

//assuming output line ids are between 0 and 'n'
const noOutputLineSelected = -1

var noOutputChnlsErr = fmt.Errorf("output Channels not set for demux")
var noSelectorErr = fmt.Errorf("no selector function configured")

type OutputLineId int
type Signal interface{}

//on error returns negative outputline id. on success, returns chosen output line and nil error
type DemuxSelectorFunc func(inputSignal Signal) (outputLineId int, err error)

// i am a de-multiplexer http://electronics-course.com/demux
type Demultiplexer interface {
	//on error returns negative outputline id. on success, returns chosen output line and nil error
	Demux(input Signal) (outputLineId int, err error)
}

// a demux whose outputlines are channels.
type ChannelDemux struct {
	OutputChannels map[OutputLineId]chan Signal
	Selector       DemuxSelectorFunc
}

func (d *ChannelDemux) Demux(input Signal) (int, error) {
	if len(d.OutputChannels) == 0 {
		return noOutputLineSelected, noOutputChnlsErr
	}

	if d.Selector == nil {
		return noOutputLineSelected, noSelectorErr
	}

	outputlineId, err := d.Selector(input)
	if err != nil {
		return noOutputLineSelected, err
	}

	ch, ok := d.OutputChannels[OutputLineId(outputlineId)]
	if !ok {
		return noOutputLineSelected, fmt.Errorf("output line with id %v was selected but not found as an output line", outputlineId)
	}
	ch <- input
	return outputlineId, nil
}

type Adder interface {
	Add(elem interface{})
}

type AddFunc func(elem interface{})

func (callMeWith AddFunc) Add(elem interface{}) {
	callMeWith(elem)
}

// a generic demux. i.e output channels are defined as something which comply with the Adder interface.
type GenericDemux struct {
	OutputLines map[OutputLineId]Adder
	Selector    DemuxSelectorFunc
}

func (d *GenericDemux) Demux(input Signal) (int, error) {
	if len(d.OutputLines) == 0 {
		return noOutputLineSelected, noOutputChnlsErr
	}

	if d.Selector == nil {
		return noOutputLineSelected, noSelectorErr
	}

	outputlineId, err := d.Selector(input)
	if err != nil {
		return noOutputLineSelected, err
	}

	ch, ok := d.OutputLines[OutputLineId(outputlineId)]
	if !ok {
		return noOutputLineSelected, fmt.Errorf("output line with id %v was selected but not found as an output line", outputlineId)
	}
	ch.Add(input)
	return outputlineId, nil
}
