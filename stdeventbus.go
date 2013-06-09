package eventbus

import (
	"log"
	"reflect"
	"strings"
)

const methodPrefix = "EventHandler"

//******************************************************************************
//events
//******************************************************************************
type subscribeEvent struct {
	subscriber Subscriber
}

type unsubscribeEvent struct {
	subscriber Subscriber
}

type shutdownEvent struct {
}

//******************************************************************************
//standard event bus
//******************************************************************************
type stdEventBus struct {
	type2Handlers map[reflect.Type][]*handler
	eventQueue    chan BusEvent
}

type handler struct {
	method     *reflect.Method
	subscriber Subscriber
}

func newStdEventBus() *stdEventBus {
	bus := new(stdEventBus)
	bus.type2Handlers = make(map[reflect.Type][]*handler)
	bus.eventQueue = make(chan BusEvent)
	go bus.listen()
	return bus
}

func (bus *stdEventBus) Subscribe(subscriber Subscriber) {
	bus.eventQueue <- &subscribeEvent{subscriber}
}

func (bus *stdEventBus) Unsubscribe(subscriber Subscriber) {
	bus.eventQueue <- &unsubscribeEvent{subscriber}
}

func (bus *stdEventBus) Shutdown() {
	bus.eventQueue <- &shutdownEvent{}
}

func (bus *stdEventBus) Publish(event BusEvent) {
	bus.eventQueue <- event
}

func (bus *stdEventBus) listen() {
	for {
		evt := <-bus.eventQueue

		switch evt.(type) {
		case *subscribeEvent:
			bus.subscribe(evt.(*subscribeEvent).subscriber)

		case *unsubscribeEvent:
			bus.unsubscribe(evt.(*unsubscribeEvent).subscriber)

		case *shutdownEvent:
			bus.shutdown()
			return

		default:
			bus.notifySubscribers(evt)
		}
	}
}

//unsubscribe
func (bus *stdEventBus) unsubscribe(subscriber Subscriber) {
	for key, handlers := range bus.type2Handlers {
		for i, h := range handlers {
			if h.subscriber == subscriber {
				//remove handler
				handlers = append(handlers[:i], handlers[i+1:]...)
			}
		}

		if len(handlers) <= 0 {
			delete(bus.type2Handlers, key)
		}
	}
}

//subscribe
func (bus *stdEventBus) subscribe(subscriber Subscriber) {
	if subscriber == nil {
		return
	}

	if bus.isAlreadySubscribed(subscriber) {
		log.Printf("EventBus: (%v, %p) was already subscribed.", subscriber, subscriber)
		return
	}

	subscriberType := reflect.TypeOf(subscriber)
	numMethods := subscriberType.NumMethod()
	handlerAdded := false
	for i := 0; i < numMethods; i++ {
		method := subscriberType.Method(i)
		methodname := method.Name

		//check method methodname
		if !strings.HasPrefix(methodname, methodPrefix) {
			continue
		}

		if !method.Func.CanInterface() {
			log.Printf("EventBus: method %s ignored (capital letter required)", methodname)
			continue
		}

		if method.Type.NumIn() != 2 {
			log.Printf("EventBus: method %s ignored (single parameter required).", methodname)
			continue
		}

		in1 := method.Type.In(1)
		h := &handler{&method, subscriber}
		bus.putHandler(in1, h)
		handlerAdded = true
	}

	if !handlerAdded {
		log.Printf("EventBus: (%v, %p) has no event handler.", subscriber, subscriber)
	}
}

func (bus *stdEventBus) isAlreadySubscribed(subscriber Subscriber) bool {
	for _, handlers := range bus.type2Handlers {
		for _, h := range handlers {
			if h.subscriber == subscriber {
				return true
			}
		}
	}
	return false
}

func (bus *stdEventBus) putHandler(key reflect.Type, h *handler) {
	bus.type2Handlers[key] = append(bus.type2Handlers[key], h)
}

//shutdown eventbus
func (bus *stdEventBus) shutdown() {
	close(bus.eventQueue)

	for key, _ := range bus.type2Handlers {
		delete(bus.type2Handlers, key)
	}
}

//notify subscriber
func (bus *stdEventBus) notifySubscribers(event BusEvent) {
	key := reflect.TypeOf(event)
	handlers, present := bus.type2Handlers[key]

	if !present {
		log.Printf("EventBus: no handlers for type %v (%v) found", key, event)
		return
	}

	for _, h := range handlers {
		m := h.method
		m.Func.Call([]reflect.Value{reflect.ValueOf(h.subscriber),
			reflect.ValueOf(event)})
	}
}
