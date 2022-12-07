package app

import (
	"log"

	"github.com/statechannels/go-nitro/client/engine/store"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

type AppManager struct {
	logger *log.Logger
	store  store.Store

	apps map[string]*App

	stopChs   map[string]chan struct{}
	MessageCh chan protocols.Message
}

func NewAppManager(logger *log.Logger, sto store.Store) *AppManager {
	return &AppManager{
		logger: logger,
		store:  sto,

		apps: map[string]*App{},

		stopChs:   map[string]chan struct{}{},
		MessageCh: make(chan protocols.Message),
	}
}

func (m *AppManager) RegisterApp(app *App) {
	if _, ok := m.apps[app.Id()]; ok {
		m.logger.Printf("WARN: App %s already registered, ignoring", app.Id())

		return
	}

	m.apps[app.Id()] = app

	m.stopChs[app.Id()] = make(chan struct{})

	go func() {
		for {
			select {
			case <-m.stopChs[app.Id()]:
				return

			case msg := <-app.MessageCh:
				m.MessageCh <- msg
			}
		}
	}()

	m.logger.Printf("INFO: App %s registered", app.Id())
}

func (m *AppManager) UnregisterApp(app *App) {
	if _, ok := m.apps[app.Id()]; !ok {
		m.logger.Printf("WARN: App %s not registered, ignoring", app.Id())

		return
	}

	close(m.stopChs[app.Id()])

	delete(m.apps, app.Id())

	m.logger.Printf("INFO: App %s unregistered", app.Id())
}

func (m *AppManager) HandleRequest(req *types.AppRequest) error {
	app, ok := m.apps[req.AppId]
	if !ok {
		return ErrAppNotRegistered
	}

	ch, err := m.store.GetConsensusChannelById(req.ChannelId)
	if err != nil {
		m.logger.Printf("ERR: %s", err)

		return ErrChannelNotFound
	}

	err = app.HandleRequest(ch, req.From, req.RequestType, req.Data)
	if err != nil {
		return err
	}

	// NOTE: this is a temporary solution, more optimized way can be achieved later on
	return m.store.SetConsensusChannel(ch)
}
