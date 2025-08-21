package options

import "time"

type AdapterOptions struct {
	Corp                string
	Model               string
	Key                 string
	BaseUrl             string
	Path                string
	IsSupportSystemRole *bool
	IsSupportStream     *bool
	Timeout             time.Duration
	ProxyUrl            string
}
