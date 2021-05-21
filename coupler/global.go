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
		OnPlugCreated: func(data interface{}) error {
			d := data.(struct {
				plug []byte
			})
			if err := handlers.HandlePlugCreated(gjson.Parse(string(d.plug))); err != nil {
				return err
			}
			return nil
		},
		OnPlugCoupled: func(data interface{}) error {
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
				return err
			}
			return nil
		},
		OnPlugUpdated: func(data interface{}) error {
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
				return err
			}
			return nil
		},
		OnSocketCreated: func(data interface{}) error {
			d := data.(struct {
				socket []byte
			})
			if err := handlers.HandleSocketCreated(gjson.Parse(string(d.socket))); err != nil {
				return err
			}
			return nil
		},
		OnSocketCoupled: func(data interface{}) error {
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
				return err
			}
			return nil
		},
		OnSocketUpdated: func(data interface{}) error {
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
				return err
			}
			return nil
		},
		OnDecoupled: func(data interface{}) error {
			d := data.(struct {
				plug   []byte
				socket []byte
			})
			if err := handlers.HandleDecoupled(
				gjson.Parse(string(d.plug)), gjson.Parse(string(d.socket)),
			); err != nil {
				return err
			}
			return nil
		},
		OnBroken: func(data interface{}) error {
			err := handlers.HandleBroken()
			if err != nil {
				return err
			}
			return nil
		},
	})
	return globalCoupler
}

var GlobalCoupler Coupler = CreateGlobalCoupler()
