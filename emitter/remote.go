package emitter

import (
	"github.com/ArthurHlt/gridana/model"
)

var remoteEmitter RemoteEmitter

type RemoteEmitter interface {
	Emit(alert model.FormattedAlert)
	Receive(func(model.FormattedAlert))
}

func AttachRemote(rmtEmitter RemoteEmitter) {
	rmtEmitter.Receive(Emit)
	go func() {
		for event := range On() {
			alert := ToAlert(event)
			rmtEmitter.Emit(alert)
		}
	}()
}
