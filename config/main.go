/**
 * File: /main.go
 * Project: integration-operator
 * File Created: 27-06-2021 02:53:17
 * Author: Clay Risser <email@clayrisser.com>
 * -----
 * Last Modified: 14-08-2022 14:34:43
 * Modified By: Clay Risser <email@clayrisser.com>
 * -----
 * Risser Labs LLC (c) Copyright 2021
 */

package config

import (
	"os"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var MaxRequeueDuration time.Duration = time.Duration(float64(time.Hour.Nanoseconds() * 6))

var StartTime metav1.Time = metav1.Now()

var DebugPlugEndpoint = os.Getenv("DEBUG_PLUG_ENDPOINT")

var DebugSocketEndpoint = os.Getenv("DEBUG_SOCKET_ENDPOINT")
