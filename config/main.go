package config

import "time"

var MaxRequeueDuration time.Duration = time.Duration(float64(time.Hour.Nanoseconds() * 6))
