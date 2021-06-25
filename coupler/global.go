package coupler

import (
	"github.com/tidwall/gjson"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"

	ctrl "sigs.k8s.io/controller-runtime"
)

var (
	couplerLog      = ctrl.Log.WithName("coupler")
	decUnstructured = yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
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
				plug         []byte
				socket       []byte
				plugConfig   map[string]string
				socketConfig map[string]string
			})
			if err := handlers.HandlePlugCoupled(
				gjson.Parse(string(d.plug)),
				gjson.Parse(string(d.socket)),
				d.plugConfig,
				d.socketConfig,
			); err != nil {
				return err
			}
			return nil
		},
		OnPlugUpdated: func(data interface{}) error {
			d := data.(struct {
				plug         []byte
				socket       []byte
				plugConfig   map[string]string
				socketConfig map[string]string
			})
			if err := handlers.HandlePlugUpdated(
				gjson.Parse(string(d.plug)),
				gjson.Parse(string(d.socket)),
				d.plugConfig,
				d.socketConfig,
			); err != nil {
				return err
			}
			return nil
		},
		OnPlugDecoupled: func(data interface{}) error {
			d := data.(struct {
				plug         []byte
				socket       []byte
				plugConfig   map[string]string
				socketConfig map[string]string
			})
			if err := handlers.HandlePlugDecoupled(
				gjson.Parse(string(d.plug)),
				gjson.Parse(string(d.socket)),
				d.plugConfig,
				d.socketConfig,
			); err != nil {
				return err
			}
			return nil
		},
		OnPlugDeleted: func(data interface{}) error {
			d := data.(struct {
				plug []byte
			})
			if err := handlers.HandlePlugDeleted(
				gjson.Parse(string(d.plug)),
			); err != nil {
				return err
			}
			return nil
		},
		OnPlugBroken: func(data interface{}) error {
			d := data.(struct {
				plug []byte
			})
			err := handlers.HandlePlugBroken(
				gjson.Parse(string(d.plug)),
			)
			if err != nil {
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
				plug         []byte
				socket       []byte
				plugConfig   map[string]string
				socketConfig map[string]string
			})
			if err := handlers.HandleSocketCoupled(
				gjson.Parse(string(d.plug)),
				gjson.Parse(string(d.socket)),
				d.plugConfig,
				d.socketConfig,
			); err != nil {
				return err
			}
			return nil
		},
		OnSocketUpdated: func(data interface{}) error {
			d := data.(struct {
				plug         []byte
				socket       []byte
				plugConfig   map[string]string
				socketConfig map[string]string
			})
			if err := handlers.HandleSocketUpdated(
				gjson.Parse(string(d.plug)),
				gjson.Parse(string(d.socket)),
				d.plugConfig,
				d.socketConfig,
			); err != nil {
				return err
			}
			return nil
		},
		OnSocketDecoupled: func(data interface{}) error {
			d := data.(struct {
				plug         []byte
				socket       []byte
				plugConfig   map[string]string
				socketConfig map[string]string
			})
			if err := handlers.HandleSocketDecoupled(
				gjson.Parse(string(d.plug)),
				gjson.Parse(string(d.socket)),
				d.plugConfig,
				d.socketConfig,
			); err != nil {
				return err
			}
			return nil
		},
		OnSocketDeleted: func(data interface{}) error {
			d := data.(struct {
				socket []byte
			})
			if err := handlers.HandleSocketDeleted(
				gjson.Parse(string(d.socket)),
			); err != nil {
				return err
			}
			return nil
		},
		OnSocketBroken: func(data interface{}) error {
			d := data.(struct {
				socket []byte
			})
			err := handlers.HandlePlugBroken(
				gjson.Parse(string(d.socket)),
			)
			if err != nil {
				return err
			}
			return nil
		},
	})
	return globalCoupler
}

var GlobalCoupler Coupler = CreateGlobalCoupler()
