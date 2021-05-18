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
			if err := handlers.HandlePlugCreated(gjson.Parse(string(d.plug))); err != nil {
				couplerLog.Error(err, "failed to handle plug created")
			}
		},
		OnPlugJoined: func(data interface{}) {
			d := data.(struct {
				plug   []byte
				socket []byte
				config []byte
			})
			if err := handlers.HandlePlugJoined(
				gjson.Parse(string(d.plug)),
				gjson.Parse(string(d.socket)),
				gjson.Parse(string(d.config)),
			); err != nil {
				couplerLog.Error(err, "failed to handle joined")
			}
		},
		OnPlugChanged: func(data interface{}) {
			d := data.(struct {
				plug   []byte
				socket []byte
				config []byte
			})
			if err := handlers.HandlePlugChanged(
				gjson.Parse(string(d.plug)),
				gjson.Parse(string(d.socket)),
				gjson.Parse(string(d.config)),
			); err != nil {
				couplerLog.Error(err, "failed to handle plug changed")
			}
		},
		OnSocketCreated: func(data interface{}) {
			d := data.(struct {
				socket []byte
			})
			if err := handlers.HandleSocketCreated(gjson.Parse(string(d.socket))); err != nil {
				couplerLog.Error(err, "failed to handle socket created")
			}
		},
		OnSocketJoined: func(data interface{}) {
			d := data.(struct {
				plug   []byte
				socket []byte
				config []byte
			})
			if err := handlers.HandleSocketJoined(
				gjson.Parse(string(d.plug)),
				gjson.Parse(string(d.socket)),
				gjson.Parse(string(d.config)),
			); err != nil {
				couplerLog.Error(err, "failed to handle joined")
			}
		},
		OnSocketChanged: func(data interface{}) {
			d := data.(struct {
				plug   []byte
				socket []byte
				config []byte
			})
			if err := handlers.HandleSocketChanged(
				gjson.Parse(string(d.plug)),
				gjson.Parse(string(d.socket)),
				gjson.Parse(string(d.config)),
			); err != nil {
				couplerLog.Error(err, "failed to handle socket changed")
			}
		},
		OnDeparted: func(data interface{}) {
			d := data.(struct {
				plug   []byte
				socket []byte
			})
			if err := handlers.HandleDeparted(
				gjson.Parse(string(d.plug)), gjson.Parse(string(d.socket)),
			); err != nil {
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
