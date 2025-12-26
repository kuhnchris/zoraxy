package main

import (
	"encoding/json"
)

// ProxyConfigs is a convenience alias for an array of ProxyConfig.
type ProxyConfigs []ProxyConfig

// ProxyConfig represents a single proxy configuration entry.
// It stores all known fields explicitly and also keeps any extra (unknown) fields
// in ExtraFields so it can round-trip JSON with arbitrary additional properties.
type ProxyConfig struct {
	ProxyType                      *int                   `json:"ProxyType,omitempty"`
	RootOrMatchingDomain           *string                `json:"RootOrMatchingDomain,omitempty"`
	MatchingDomainAlias            []string               `json:"MatchingDomainAlias,omitempty"`
	ActiveOrigins                  []Origin               `json:"ActiveOrigins,omitempty"`
	InactiveOrigins                []Origin               `json:"InactiveOrigins,omitempty"`
	UseStickySession               *bool                  `json:"UseStickySession,omitempty"`
	UseActiveLoadBalance           *bool                  `json:"UseActiveLoadBalance,omitempty"`
	Disabled                       *bool                  `json:"Disabled,omitempty"`
	ListeningPorts                 json.RawMessage        `json:"ListeningPorts,omitempty"` // allow null or array
	BypassGlobalTLS                *bool                  `json:"BypassGlobalTLS,omitempty"`
	TlsOptions                     TLSOptions             `json:"TlsOptions,omitempty"`
	VirtualDirectories             []string               `json:"VirtualDirectories,omitempty"`
	HeaderRewriteRules             HeaderRewriteRules     `json:"HeaderRewriteRules,omitempty"`
	EnableWebsocketCustomHeaders   *bool                  `json:"EnableWebsocketCustomHeaders,omitempty"`
	AuthenticationProvider         AuthenticationProvider `json:"AuthenticationProvider,omitempty"`
	RequireRateLimit               *bool                  `json:"RequireRateLimit,omitempty"`
	RateLimit                      *int                   `json:"RateLimit,omitempty"`
	RequireCaptcha                 *bool                  `json:"RequireCaptcha,omitempty"`
	CaptchaConfig                  json.RawMessage        `json:"CaptchaConfig,omitempty"` // allow null
	DisableUptimeMonitor           *bool                  `json:"DisableUptimeMonitor,omitempty"`
	DisableAutoFallback            *bool                  `json:"DisableAutoFallback,omitempty"`
	DisableLogging                 *bool                  `json:"DisableLogging,omitempty"`
	DisableStatisticCollection     *bool                  `json:"DisableStatisticCollection,omitempty"`
	BlockCommonExploits            *bool                  `json:"BlockCommonExploits,omitempty"`
	BlockAICrawlers                *bool                  `json:"BlockAICrawlers,omitempty"`
	MitigationAction               *int                   `json:"MitigationAction,omitempty"`
	DisableChunkedTransferEncoding *bool                  `json:"DisableChunkedTransferEncoding,omitempty"`
	AccessFilterUUID               *string                `json:"AccessFilterUUID,omitempty"`
	DefaultSiteOption              *int                   `json:"DefaultSiteOption,omitempty"`
	DefaultSiteValue               *string                `json:"DefaultSiteValue,omitempty"`
	Tags                           []string               `json:"Tags,omitempty"`

	// ExtraFields holds any additional fields present in the JSON that are not
	// explicitly declared above. This makes the struct flexible to future changes.
	ExtraFields map[string]json.RawMessage `json:"-"`
}

type Origin struct {
	OriginIpOrDomain         *string `json:"OriginIpOrDomain,omitempty"`
	RequireTLS               *bool   `json:"RequireTLS,omitempty"`
	SkipCertValidations      *bool   `json:"SkipCertValidations,omitempty"`
	SkipWebSocketOriginCheck *bool   `json:"SkipWebSocketOriginCheck,omitempty"`
	Weight                   *int    `json:"Weight,omitempty"`
	MaxConn                  *int    `json:"MaxConn,omitempty"`
	RespTimeout              *int    `json:"RespTimeout,omitempty"`

	ExtraFields map[string]json.RawMessage `json:"-"`
}

type TLSOptions struct {
	DisableSNI                       *bool                      `json:"DisableSNI,omitempty"`
	DisableLegacyCertificateMatching *bool                      `json:"DisableLegacyCertificateMatching,omitempty"`
	EnableAutoHTTPS                  *bool                      `json:"EnableAutoHTTPS,omitempty"`
	PreferredCertificate             map[string]json.RawMessage `json:"-"`

	ExtraFields map[string]json.RawMessage `json:"-"`
}

type HeaderRewriteRules struct {
	UserDefinedHeaders            []map[string]string `json:"UserDefinedHeaders,omitempty"`
	RequestHostOverwrite          *string             `json:"RequestHostOverwrite,omitempty"`
	HSTSMaxAge                    *int                `json:"HSTSMaxAge,omitempty"`
	EnablePermissionPolicyHeader  *bool               `json:"EnablePermissionPolicyHeader,omitempty"`
	PermissionPolicy              json.RawMessage     `json:"PermissionPolicy,omitempty"` // allow null
	DisableHopByHopHeaderRemoval  *bool               `json:"DisableHopByHopHeaderRemoval,omitempty"`
	DisableUserAgentHeaderRemoval *bool               `json:"DisableUserAgentHeaderRemoval,omitempty"`

	ExtraFields map[string]json.RawMessage `json:"-"`
}
type AuthenticationProvider struct {
	AuthMethod                        *int                `json:"AuthMethod,omitempty"`
	BasicAuthCredentials              []map[string]string `json:"BasicAuthCredentials,omitempty"`
	BasicAuthExceptionRules           []map[string]string `json:"BasicAuthExceptionRules,omitempty"`
	BasicAuthGroupIDs                 json.RawMessage     `json:"BasicAuthGroupIDs,omitempty"` // allow null
	ForwardAuthURL                    *string             `json:"ForwardAuthURL,omitempty"`
	ForwardAuthResponseHeaders        json.RawMessage     `json:"ForwardAuthResponseHeaders,omitempty"`        // allow null
	ForwardAuthResponseClientHeaders  json.RawMessage     `json:"ForwardAuthResponseClientHeaders,omitempty"`  // allow null
	ForwardAuthRequestHeaders         json.RawMessage     `json:"ForwardAuthRequestHeaders,omitempty"`         // allow null
	ForwardAuthRequestExcludedCookies json.RawMessage     `json:"ForwardAuthRequestExcludedCookies,omitempty"` // allow null

	ExtraFields map[string]json.RawMessage `json:"-"`
}

// ProxyUpstream represents the whole JSON payload you posted.
// It can hold any additional fields that are not explicitly declared.
type ProxyUpstream struct {
	ActiveOrigins   []Origin `json:"ActiveOrigins,omitempty"`
	InactiveOrigins []Origin `json:"InactiveOrigins,omitempty"`

	// Any extra keys present in the source JSON will be stored here
	ExtraFields map[string]json.RawMessage `json:"-"`
}
