/**
 * File: /global.go
 * Project: integration-operator
 * File Created: 23-06-2021 09:14:26
 * Author: Clay Risser <email@clayrisser.com>
 * -----
 * Last Modified: 25-06-2023 14:21:54
 * Modified By: Clay Risser <email@clayrisser.com>
 * -----
 * Risser Labs LLC (c) Copyright 2021
 */

package coupler

import (
	"context"
	"encoding/json"

	integrationv1alpha2 "gitlab.com/bitspur/rock8s/integration-operator/api/v1alpha2"
	ctrl "sigs.k8s.io/controller-runtime"
)

var couplerLog = ctrl.Log.WithName("coupler")

func CreateGlobalCoupler() Coupler {
	handlers := NewHandlers()
	globalCoupler := *NewCoupler(Options{
		MaxQueueSize: 99,
		MaxWorkers:   1,
	})
	globalCoupler.RegisterEvents(&Events{
		OnPlugCreated: func(data interface{}) error {
			d := data.(struct {
				ctx  *context.Context
				plug []byte
			})
			var plug integrationv1alpha2.Plug
			if err := json.Unmarshal(d.plug, &plug); err != nil {
				return err
			}
			if err := handlers.HandlePlugCreated(d.ctx, &plug); err != nil {
				return err
			}
			return nil
		},
		OnPlugCoupled: func(data interface{}) error {
			d := data.(struct {
				ctx          *context.Context
				plug         []byte
				socket       []byte
				plugConfig   map[string]string
				socketConfig map[string]string
			})
			var plug integrationv1alpha2.Plug
			if err := json.Unmarshal(d.plug, &plug); err != nil {
				return err
			}
			var socket integrationv1alpha2.Socket
			if err := json.Unmarshal(d.socket, &socket); err != nil {
				return err
			}
			if err := handlers.HandlePlugCoupled(
				d.ctx,
				&plug,
				&socket,
				&d.plugConfig,
				&d.socketConfig,
			); err != nil {
				return err
			}
			return nil
		},
		OnPlugUpdated: func(data interface{}) error {
			d := data.(struct {
				ctx          *context.Context
				plug         []byte
				socket       []byte
				plugConfig   map[string]string
				socketConfig map[string]string
			})
			var plug integrationv1alpha2.Plug
			if err := json.Unmarshal(d.plug, &plug); err != nil {
				return err
			}
			var socket integrationv1alpha2.Socket
			if err := json.Unmarshal(d.socket, &socket); err != nil {
				return err
			}
			if err := handlers.HandlePlugUpdated(
				d.ctx,
				&plug,
				&socket,
				&d.plugConfig,
				&d.socketConfig,
			); err != nil {
				return err
			}
			return nil
		},
		OnPlugDecoupled: func(data interface{}) error {
			d := data.(struct {
				ctx          *context.Context
				plug         []byte
				socket       []byte
				plugConfig   map[string]string
				socketConfig map[string]string
			})
			var plug integrationv1alpha2.Plug
			if err := json.Unmarshal(d.plug, &plug); err != nil {
				return err
			}
			var socket integrationv1alpha2.Socket
			if err := json.Unmarshal(d.socket, &socket); err != nil {
				return err
			}
			if err := handlers.HandlePlugDecoupled(
				d.ctx,
				&plug,
				&socket,
				&d.plugConfig,
				&d.socketConfig,
			); err != nil {
				return err
			}
			return nil
		},
		OnPlugDeleted: func(data interface{}) error {
			d := data.(struct {
				ctx  *context.Context
				plug []byte
			})
			var plug integrationv1alpha2.Plug
			if err := json.Unmarshal(d.plug, &plug); err != nil {
				return err
			}
			if err := handlers.HandlePlugDeleted(
				d.ctx,
				&plug,
			); err != nil {
				return err
			}
			return nil
		},
		OnPlugBroken: func(data interface{}) error {
			d := data.(struct {
				ctx  *context.Context
				plug []byte
			})
			var plug integrationv1alpha2.Plug
			if err := json.Unmarshal(d.plug, &plug); err != nil {
				return err
			}
			err := handlers.HandlePlugBroken(
				d.ctx,
				&plug,
			)
			if err != nil {
				return err
			}
			return nil
		},
		OnSocketCreated: func(data interface{}) error {
			d := data.(struct {
				ctx    *context.Context
				socket []byte
			})
			var socket integrationv1alpha2.Socket
			if err := json.Unmarshal(d.socket, &socket); err != nil {
				return err
			}
			if err := handlers.HandleSocketCreated(d.ctx, &socket); err != nil {
				return err
			}
			return nil
		},
		OnSocketCoupled: func(data interface{}) error {
			d := data.(struct {
				ctx          *context.Context
				plug         []byte
				socket       []byte
				plugConfig   map[string]string
				socketConfig map[string]string
			})
			var plug integrationv1alpha2.Plug
			if err := json.Unmarshal(d.plug, &plug); err != nil {
				return err
			}
			var socket integrationv1alpha2.Socket
			if err := json.Unmarshal(d.socket, &socket); err != nil {
				return err
			}
			if err := handlers.HandleSocketCoupled(
				d.ctx,
				&plug,
				&socket,
				&d.plugConfig,
				&d.socketConfig,
			); err != nil {
				return err
			}
			return nil
		},
		OnSocketUpdated: func(data interface{}) error {
			d := data.(struct {
				ctx          *context.Context
				plug         []byte
				socket       []byte
				plugConfig   map[string]string
				socketConfig map[string]string
			})
			var plug integrationv1alpha2.Plug
			if err := json.Unmarshal(d.plug, &plug); err != nil {
				return err
			}
			var socket integrationv1alpha2.Socket
			if err := json.Unmarshal(d.socket, &socket); err != nil {
				return err
			}
			if err := handlers.HandleSocketUpdated(
				d.ctx,
				&plug,
				&socket,
				&d.plugConfig,
				&d.socketConfig,
			); err != nil {
				return err
			}
			return nil
		},
		OnSocketDecoupled: func(data interface{}) error {
			d := data.(struct {
				ctx          *context.Context
				plug         []byte
				socket       []byte
				plugConfig   map[string]string
				socketConfig map[string]string
			})
			var plug integrationv1alpha2.Plug
			if err := json.Unmarshal(d.plug, &plug); err != nil {
				return err
			}
			var socket integrationv1alpha2.Socket
			if err := json.Unmarshal(d.socket, &socket); err != nil {
				return err
			}
			if err := handlers.HandleSocketDecoupled(
				d.ctx,
				&plug,
				&socket,
				&d.plugConfig,
				&d.socketConfig,
			); err != nil {
				return err
			}
			return nil
		},
		OnSocketDeleted: func(data interface{}) error {
			d := data.(struct {
				ctx    *context.Context
				socket []byte
			})
			var socket integrationv1alpha2.Socket
			if err := json.Unmarshal(d.socket, &socket); err != nil {
				return err
			}
			if err := handlers.HandleSocketDeleted(
				d.ctx,
				&socket,
			); err != nil {
				return err
			}
			return nil
		},
		OnSocketBroken: func(data interface{}) error {
			d := data.(struct {
				ctx    *context.Context
				socket []byte
			})
			var socket integrationv1alpha2.Socket
			if err := json.Unmarshal(d.socket, &socket); err != nil {
				return err
			}
			err := handlers.HandleSocketBroken(
				d.ctx,
				&socket,
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
