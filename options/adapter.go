package options

import "time"

type AdapterOptions struct {
	Provider                string
	Model                   string
	Key                     string
	BaseUrl                 string
	Path                    string
	Header                  map[string]string
	Timeout                 time.Duration
	ProxyUrl                string
	Action                  string
	IsSupportSystemRole     *bool
	IsSupportStream         *bool
	IsOfficialFormatRequest bool
}
