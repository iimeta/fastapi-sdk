package options

import "time"

type AdapterOptions struct {
	Provider                string
	Model                   string
	Key                     string
	BaseUrl                 string
	Path                    string
	Action                  string
	Timeout                 time.Duration
	ProxyUrl                string
	IsSupportSystemRole     *bool
	IsSupportStream         *bool
	IsOfficialFormatRequest bool
}
