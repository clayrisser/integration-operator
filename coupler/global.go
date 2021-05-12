package coupler

import (
	"github.com/tidwall/gjson"
)

func CreateGlobalCoupler() Coupler {
	handlers := NewHandlers()
	globalCoupler := *NewCoupler(Options{
		MaxQueueSize: 99,
		MaxWorkers:   1,
	})
	globalCoupler.RegisterEvents(&Events{
		OnPlugCreated: func(data interface{}) {
			d := data.(struct {
				plug []byte
			})
			handlers.HandlePlugCreated(gjson.Parse(string(d.plug)))
		},
		OnJoined: func(data interface{}) {
			d := data.(struct {
				plug   []byte
				socket []byte
				config []byte
			})
			handlers.HandleJoined(
				gjson.Parse(string(d.plug)),
				gjson.Parse(string(d.socket)),
				gjson.Parse(string(d.config)),
			)
		},
		OnPlugChanged: func(data interface{}) {
			d := data.(struct {
				plug   []byte
				socket []byte
				config []byte
			})
			handlers.HandlePlugChanged(
				gjson.Parse(string(d.plug)),
				gjson.Parse(string(d.socket)),
				gjson.Parse(string(d.config)),
			)
		},
		OnSocketCreated: func(data interface{}) {
			d := data.(struct {
				socket []byte
			})
			handlers.HandleSocketCreated(gjson.Parse(string(d.socket)))
		},
		OnSocketChanged: func(data interface{}) {
			d := data.(struct {
				plug   []byte
				socket []byte
				config []byte
			})
			handlers.HandleSocketChanged(
				gjson.Parse(string(d.plug)),
				gjson.Parse(string(d.socket)),
				gjson.Parse(string(d.config)),
			)
		},
		OnDeparted: func(data interface{}) {
			handlers.HandleDeparted()
		},
		OnBroken: func(data interface{}) {
			handlers.HandleBroken()
		},
	})
	return globalCoupler
}

var GlobalCoupler Coupler = CreateGlobalCoupler()
