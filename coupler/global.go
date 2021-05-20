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
		OnPlugCoupled: func(data interface{}) {
			d := data.(struct {
				plug   []byte
				socket []byte
				config []byte
			})
			if err := handlers.HandlePlugCoupled(
				gjson.Parse(string(d.plug)),
				gjson.Parse(string(d.socket)),
				gjson.Parse(string(d.config)),
			); err != nil {
				couplerLog.Error(err, "failed to handle plug coupled")
			}
		},
		OnPlugUpdated: func(data interface{}) {
			d := data.(struct {
				plug   []byte
				socket []byte
				config []byte
			})
			if err := handlers.HandlePlugUpdated(
				gjson.Parse(string(d.plug)),
				gjson.Parse(string(d.socket)),
				gjson.Parse(string(d.config)),
			); err != nil {
				couplerLog.Error(err, "failed to handle plug updated")
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
		OnSocketCoupled: func(data interface{}) {
			d := data.(struct {
				plug   []byte
				socket []byte
				config []byte
			})
			if err := handlers.HandleSocketCoupled(
				gjson.Parse(string(d.plug)),
				gjson.Parse(string(d.socket)),
				gjson.Parse(string(d.config)),
			); err != nil {
				couplerLog.Error(err, "failed to handle socket coupled")
			}
		},
		OnSocketUpdated: func(data interface{}) {
			d := data.(struct {
				plug   []byte
				socket []byte
				config []byte
			})
			if err := handlers.HandleSocketUpdated(
				gjson.Parse(string(d.plug)),
				gjson.Parse(string(d.socket)),
				gjson.Parse(string(d.config)),
			); err != nil {
				couplerLog.Error(err, "failed to handle socket updated")
			}
		},
		OnDecoupled: func(data interface{}) {
			d := data.(struct {
				plug   []byte
				socket []byte
			})
			if err := handlers.HandleDecoupled(
				gjson.Parse(string(d.plug)), gjson.Parse(string(d.socket)),
			); err != nil {
				couplerLog.Error(err, "failed to handle decoupled")
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
