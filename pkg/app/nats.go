package app

import (
	"encoding/json"
	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
	"time"
)

type PubSub struct {
	Conn   *nats.Conn
	Logger *zap.SugaredLogger
}

// InitializeNats : initialize nats server and return the PubSub with Conn object
func InitializeNats(logger *zap.SugaredLogger) *PubSub {
	opts := &server.Options{}
	// Initialize new server with options
	ns, err := server.NewServer(opts)
	if err != nil {
		logger.Panicw("Nats server unable to start", "error", err)
		return nil
	}
	go ns.Start()

	if !ns.ReadyForConnections(4 * time.Second) {
		logger.Panic("Nats not ready for connection")
		return nil
	}

	nc, err := nats.Connect(ns.ClientURL())
	if err != nil {
		logger.Panicw("Failed in Connecting", "error", err)
		return nil
	}
	return &PubSub{
		Logger: logger,
		Conn:   nc,
	}
}

// RegisterSubscribers : register the map of subscriber
func (p *PubSub) RegisterSubscribers(handlerMap map[string]nats.MsgHandler) {
	for topic, handler := range handlerMap {
		_, err := p.Conn.Subscribe(topic, handler)
		if err != nil {
			p.Logger.Errorw("Error while subscribing the topic", "error", err)
		}
	}
}

// Publish : Publish message to nats
func (p *PubSub) Publish(stage Stage, data interface{}) {
	jsonData, _ := json.Marshal(data)
	err := p.Conn.Publish(string(stage), jsonData)
	if err != nil {
		p.Logger.Errorw("Error while publishing to nats", "error", err)
	}
}
