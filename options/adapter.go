package options

import "time"

type AdapterOptions struct {
	Provider            string
	Model               string
	Key                 string
	BaseUrl             string
	Path                string
	Header              map[string]string
	Timeout             time.Duration
	ProxyUrl            string
	Stream              bool
	Action              string
	IsSupportSystemRole *bool
	IsSupportStream     *bool
	ReqPassthroughParams  []string          // 请求透传参数
	ResPassthroughParams  []string          // 响应透传参数
	PassthroughHeader     map[string]string // 透传请求头
}
