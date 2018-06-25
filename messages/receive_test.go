package messages

import "testing"

func Test_Receive_From_Queue(t *testing.T) {

	SendToQueue("hi")
	receiveFromQueue()
}
