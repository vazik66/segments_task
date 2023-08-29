package events

import (
	"fmt"
)

type EventManager interface {
    Subscribe(string, func(interface{}) error) 
    Publish(string, interface{}) error 
}

type LocalEventManager struct {
	channels map[string][]chan interface{}
}

func NewLocalEventManager() EventManager {
	return &LocalEventManager{
		channels: make(map[string][]chan interface{}),
	}
}

func (m *LocalEventManager) Subscribe(event string, handler func(interface{}) error) {
    ch := make(chan interface{}, 1024)

	go func(c chan interface{}, h func(interface{}) error) {
		for {
			data, ok := <-c
			if !ok {
				break
			}
            _ = h(data)
		}
	}(ch, handler)

	m.channels[event] = append(m.channels[event], ch)
}

func (m *LocalEventManager) Publish(event string, data interface{}) error {
	_, ok := m.channels[event]
	if !ok {
		return fmt.Errorf("event not found")
	}

	for _, ch := range m.channels[event] {
		ch <- data
	}

	return nil
}
