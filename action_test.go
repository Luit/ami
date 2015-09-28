package ami

import (
	"bytes"
	"testing"
)

type Login struct {
	Username, Secret string
}

// this Originate does not contain the full set of possible values, just the
// ones we want in the test.

type Originate struct {
	Channel string
	Exten   string
	Context string
	Async   string
}

type Redirect struct {
	Channel      string
	ExtraChannel string
	Exten        string
	Priority     uint
}

type MultilineAction struct {
	Test []int
}

type Tags struct {
	Zero     int `ami:"zero"`
	NonZero  int
	Empty    string `ami:"empty"`
	NonEmpty string
}

var marshalTests = []struct {
	in  interface{}
	id  ActionID
	out []byte
	err error
}{
	{
		in: Login{
			Username: "testuser",
			Secret:   "testsecret",
		},
		out: []byte("Action: Login" + "\r\n" +
			"Username: testuser" + "\r\n" +
			"Secret: testsecret" + "\r\n" +
			"\r\n"),
	},
	{
		in: Originate{
			Channel: "sip/12345",
			Exten:   "1234",
			Context: "default",
		},
		out: []byte("Action: Originate" + "\r\n" +
			"Channel: sip/12345" + "\r\n" +
			"Exten: 1234" + "\r\n" +
			"Context: default" + "\r\n" +
			"\r\n"),
	},
	{
		in: Originate{
			Channel: "sip/12345",
			Exten:   "1234",
			Context: "default",
			Async:   "yes",
		},
		id: 300,
		out: []byte("Action: Originate" + "\r\n" +
			"ActionID: 300" + "\r\n" +
			"Channel: sip/12345" + "\r\n" +
			"Exten: 1234" + "\r\n" +
			"Context: default" + "\r\n" +
			"Async: yes" + "\r\n" +
			"\r\n"),
	},
	{
		in: Redirect{
			Channel:      "DAHDI/1-1",
			ExtraChannel: "SIP/3064-7e00",
			Exten:        "680",
			Priority:     1,
		},
		out: []byte("Action: Redirect" + "\r\n" +
			"Channel: DAHDI/1-1" + "\r\n" +
			"ExtraChannel: SIP/3064-7e00" + "\r\n" +
			"Exten: 680" + "\r\n" +
			"Priority: 1" + "\r\n" +
			"\r\n"),
	},
	{
		in:  "test",
		err: errMarshalType,
	},
	{
		in: MultilineAction{
			Test: []int{1, 2, 3},
		},
		out: []byte("Action: MultilineAction" + "\r\n" +
			"Test: 1" + "\r\n" +
			"Test: 2" + "\r\n" +
			"Test: 3" + "\r\n" +
			"\r\n"),
	},
	{
		in: Tags{},
		out: []byte("Action: Tags" + "\r\n" +
			"Zero: 0" + "\r\n" +
			"Empty: " + "\r\n" +
			"\r\n"),
	},
}

func TestMarshalAction(t *testing.T) {
	for i, test := range marshalTests {
		b, err := marshalAction(test.in, test.id)
		if test.err != err {
			if test.err != nil {
				t.Fatalf("Expected test %d to fail with error %v, got %v", i, test.err, err)
			}
		}
		if !bytes.Equal(test.out, b) {
			t.Fatalf("Got wrong output on test %d\nexpected: \n%#v\ngot: \n%#v", i, string(test.out), string(b))
		}
	}
}
