package engine

import (
	"errors"
	"log"

	"github.com/statechannels/go-nitro/client/engine/store"
	"github.com/statechannels/go-nitro/types"
)

type VirtualApp interface {
	handleVirtualAppRequest(s *store.Store, request *VirtualAppRequest) error
}

type VirtualAppRequest struct {
	ChannelId types.Destination
	AppId     string
	AppData   interface{}
}

type VirtualAppManager struct {
	logger        *log.Logger
	store         *store.Store
	vAppsRegistry map[string]VirtualApp
}

func NewVirtualAppManager(logger *log.Logger, store *store.Store) *VirtualAppManager {
	return &VirtualAppManager{
		logger:        logger,
		store:         store,
		vAppsRegistry: map[string]VirtualApp{},
	}
}

func (vam *VirtualAppManager) RegisterVirtualApp(appId string, app VirtualApp) {
	_, exists := vam.vAppsRegistry[appId]
	if exists {
		vam.logger.Printf("WARN: Virtual App %s already registered, ignoring", appId)
	} else {
		vam.vAppsRegistry[appId] = app
		vam.logger.Printf("INFO: Virtual App %s registered", appId)
	}
}

func (vam *VirtualAppManager) UnregisterVirtualApp(appId string) {
	delete(vam.vAppsRegistry, appId)
	vam.logger.Printf("INFO: Virtual App %s unregistered", appId)
}

func (vam *VirtualAppManager) HandleRequest(req *VirtualAppRequest) error {
	vApp, ok := vam.vAppsRegistry[req.AppId]
	if !ok {
		return errors.New("virtual app " + req.AppId + " not registered")
	}
	return vApp.handleVirtualAppRequest(vam.store, req)
}
