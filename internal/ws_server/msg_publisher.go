package wsserver

import "github.com/tjarratt/babble"

type MessagePublisher interface {
	Publish(msg []byte)
}

type RandomWordPublisher struct {
	sendMsgChan chan<- []byte
}

func NewRandomWordPublisher(sendMsgChan chan<- []byte) MessagePublisher {
	return &RandomWordPublisher{
		sendMsgChan: sendMsgChan,
	}
}

func (rwp *RandomWordPublisher) Publish(msg []byte) {
	babbler := babble.NewBabbler()
	babbler.Count = 1
	rndWord := babbler.Babble()

	msgToSend := string(msg) + "-" + rndWord
	rwp.sendMsgChan <- []byte(msgToSend)
}
