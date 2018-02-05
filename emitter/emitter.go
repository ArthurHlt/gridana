package emitter

import (
	"github.com/ArthurHlt/gridana/model"
	"github.com/olebedev/emitter"
)

var e *emitter.Emitter = emitter.New(uint(100))

func Emit(alert model.FormattedAlert) {
	e.Emit("alert", alert)
}

func On() <-chan emitter.Event {

	return e.On("alert", emitter.Sync)
}
func Off(events ...<-chan emitter.Event) {
	e.Off("alert", events...)
}
func Listeners() []<-chan emitter.Event {
	return e.Listeners("alert")
}
func ToAlert(evt emitter.Event) model.FormattedAlert {
	return evt.Args[0].(model.FormattedAlert)
}
