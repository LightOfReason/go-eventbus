package eventbus

type Subscriber interface{}
type BusEvent interface{}

type EventBus interface {
	Subscribe(subscriber Subscriber)
	Unsubscribe(subscriber Subscriber)
	Publish(event BusEvent)
	Shutdown()
}

var eventBus EventBus = newStdEventBus()

func Subscribe(subscriber Subscriber) {
	eventBus.Subscribe(subscriber)
}

func Unsubscribe(subscriber Subscriber) {
	eventBus.Unsubscribe(subscriber)
}

func Publish(event BusEvent) {
	eventBus.Publish(event)
}

func Shutdown() {
	eventBus.Shutdown()
}
