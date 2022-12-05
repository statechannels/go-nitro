package apps

import (
	"log"

	"github.com/statechannels/go-nitro/client/engine/store"
)

type AppManager struct {
	logger *log.Logger
	store  store.Store

	apps map[string]App
}

func NewAppManager(logger *log.Logger, sto store.Store) *AppManager {
	m := &AppManager{
		logger: logger,
		store:  sto,

		apps: make(map[string]App),
	}

	return m
}

func (m *AppManager) RegisterApp(app App) {
	if _, ok := m.apps[app.Type()]; ok {
		m.logger.Printf("WARN: App %s already registered, ignoring", app.Type())

		return
	}

	m.apps[app.Type()] = app

	m.logger.Printf("INFO: App %s registered", app.Type())
}

func (m *AppManager) UnregisterApp(app App) {
	delete(m.apps, app.Type())

	m.logger.Printf("INFO: App %s unregistered", app.Type())
}

func (m *AppManager) HandleRequest(req AppRequest) error {
	app, ok := m.apps[req.AppType]
	if !ok {
		return ErrAppNotRegistered
	}

	ch, ok := m.store.GetChannelById(req.ChannelId)
	if !ok {
		return ErrChannelNotFound
	}

	err := app.HandleRequest(ch, req.RequestType, req.Data)
	if err != nil {
		return err
	}

	// NOTE: this is a temporary solution, more optimized way can be achieved later on
	return m.store.SetChannel(ch)
}
