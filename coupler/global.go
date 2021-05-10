package coupler

import integrationv1alpha2 "github.com/silicon-hills/integration-operator/api/v1alpha2"

func CreateGlobalCoupler() Coupler {
	handlers := NewHandlers()
	globalCoupler := *NewCoupler(Options{
		MaxQueueSize: 99,
		MaxWorkers:   1,
	})
	globalCoupler.RegisterEvents(&Events{
		OnPlugCreated: func(data interface{}) {
			d := data.(struct {
				plug *integrationv1alpha2.Plug
			})
			handlers.HandlePlugCreated(d.plug)
		},
		OnJoined: func(data interface{}) {
			d := data.(struct {
				plug    *integrationv1alpha2.Plug
				socket  *integrationv1alpha2.Socket
				payload interface{}
			})
			handlers.HandleJoined(d.plug, d.socket, d.payload)
		},
		OnPlugChanged: func(data interface{}) {
			d := data.(struct {
				plug    *integrationv1alpha2.Plug
				socket  *integrationv1alpha2.Socket
				payload interface{}
			})
			handlers.HandlePlugChanged(d.plug, d.socket, d.payload)
		},
		OnSocketCreated: func(data interface{}) {
			d := data.(struct {
				socket *integrationv1alpha2.Socket
			})
			handlers.HandleSocketCreated(d.socket)
		},
		OnSocketChanged: func(data interface{}) {
			d := data.(struct {
				plug    *integrationv1alpha2.Plug
				socket  *integrationv1alpha2.Socket
				payload interface{}
			})
			handlers.HandleSocketChanged(d.plug, d.socket, d.payload)
		},
		OnDeparted: func(data interface{}) {
			handlers.HandleDeparted(nil, nil, nil)
		},
		OnBroken: func(data interface{}) {
			d := data.(struct {
				plug   *integrationv1alpha2.Plug
				socket *integrationv1alpha2.Socket
			})
			handlers.HandleBroken(d.plug, d.socket)
		},
	})
	return globalCoupler
}

var GlobalCoupler Coupler = CreateGlobalCoupler()
