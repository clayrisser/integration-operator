package coupler

import (
	"github.com/tidwall/gjson"

	ctrl "sigs.k8s.io/controller-runtime"
)

var (
	couplerLog = ctrl.Log.WithName("coupler")
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
			err := handlers.HandlePlugCreated(gjson.Parse(string(d.plug)))
			if err != nil {
				couplerLog.Error(err, "failed to handle plug created")
			}
		},
		OnJoined: func(data interface{}) {
			d := data.(struct {
				plug   []byte
				socket []byte
				config []byte
			})
			err := handlers.HandleJoined(
				gjson.Parse(string(d.plug)),
				gjson.Parse(string(d.socket)),
				gjson.Parse(string(d.config)),
			)
			if err != nil {
				couplerLog.Error(err, "failed to handle joined")
			}
		},
		OnPlugChanged: func(data interface{}) {
			d := data.(struct {
				plug   []byte
				socket []byte
				config []byte
			})
			err := handlers.HandlePlugChanged(
				gjson.Parse(string(d.plug)),
				gjson.Parse(string(d.socket)),
				gjson.Parse(string(d.config)),
			)
			if err != nil {
				couplerLog.Error(err, "failed to handle plug changed")
			}
		},
		OnSocketCreated: func(data interface{}) {
			d := data.(struct {
				socket []byte
			})
			err := handlers.HandleSocketCreated(gjson.Parse(string(d.socket)))
			if err != nil {
				couplerLog.Error(err, "failed to handle socket created")
			}
		},
		OnSocketChanged: func(data interface{}) {
			d := data.(struct {
				plug   []byte
				socket []byte
				config []byte
			})
			err := handlers.HandleSocketChanged(
				gjson.Parse(string(d.plug)),
				gjson.Parse(string(d.socket)),
				gjson.Parse(string(d.config)),
			)
			if err != nil {
				couplerLog.Error(err, "failed to handle socket changed")
			}
		},
		OnDeparted: func(data interface{}) {
			err := handlers.HandleDeparted()
			if err != nil {
				couplerLog.Error(err, "failed to handle departed")
			}
		},
		OnBroken: func(data interface{}) {
			err := handlers.HandleBroken()
			if err != nil {
				couplerLog.Error(err, "failed to handle broken")
			}
		},
	})
	return globalCoupler
}

var GlobalCoupler Coupler = CreateGlobalCoupler()
