package go_demux

import "testing"

func TestMurmurKeyBasedDemuxFunc(t *testing.T) {

	type test struct {
		input        []byte
		lineExpected int
		errExpected  bool
	}

	tests := []test{
		{[]byte("1647fec8408-39c40000010a443b"), 62, false},
		{[]byte("1647fec8408-39c40000010a443b"), 62, false},
		{[]byte("1647ee3d0a6-58030000010f610d"), 97, false},
		{[]byte("1647ea7b83a-5ec60000010e2cbf"), 35, false},
	}

	md := NewMurmurKeyBasedDemuxFunc(100)

	for _, tc := range tests {
		got, err := md(tc.input)
		if err != nil && !tc.errExpected {
			t.Errorf("did not expect error for input: %v, but got %v", tc.input, err)
		}
		if got != tc.lineExpected {
			t.Errorf("expected %v , got %v", tc.lineExpected, got)
		}
	}
}
