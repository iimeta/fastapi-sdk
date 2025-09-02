package options

import "time"

type AdapterOptions struct {
	Corp                    string
	Model                   string
	Key                     string
	BaseUrl                 string
	Path                    string
	Timeout                 time.Duration
	ProxyUrl                string
	IsSupportSystemRole     *bool
	IsSupportStream         *bool
	IsOfficialFormatRequest bool
}
