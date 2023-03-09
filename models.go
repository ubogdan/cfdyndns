package cfdyndns

import (
	"time"
)

const (
	userAgent = "cfDynDNS/0.1.0 (github.com/ubogdan/cfdyndns)"
	jsonMime  = "application/json"
)

type ZoneQuery struct {
	Result     []Zone        `json:"result"`
	ResultInfo ResultInfo    `json:"result_info"`
	Success    bool          `json:"success"`
	Errors     []interface{} `json:"errors"`
	Messages   []interface{} `json:"messages"`
}

type RecordQuery struct {
	Result     []Record      `json:"result"`
	Success    bool          `json:"success"`
	Errors     []interface{} `json:"errors"`
	Messages   []interface{} `json:"messages"`
	ResultInfo ResultInfo    `json:"result_info"`
}

type Zone struct {
	Id              string    `json:"id"`
	Name            string    `json:"name"`
	Status          string    `json:"status"`
	Paused          bool      `json:"paused"`
	Type            string    `json:"type"`
	DevelopmentMode int       `json:"development_mode"`
	NameServers     []string  `json:"name_servers"`
	ModifiedOn      time.Time `json:"modified_on"`
	CreatedOn       time.Time `json:"created_on"`
	ActivatedOn     time.Time `json:"activated_on"`
	Meta            ZoneMeta  `json:"meta"`
	Owner           struct {
		Id    interface{} `json:"id"`
		Type  string      `json:"type"`
		Email interface{} `json:"email"`
	} `json:"owner"`
	Account struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	} `json:"account"`
	Permissions []string `json:"permissions"`
	Plan        Plan     `json:"plan"`
}

type ZoneMeta struct {
	Step                    int  `json:"step"`
	CustomCertificateQuota  int  `json:"custom_certificate_quota"`
	PageRuleQuota           int  `json:"page_rule_quota"`
	PhishingDetected        bool `json:"phishing_detected"`
	MultipleRailgunsAllowed bool `json:"multiple_railguns_allowed"`
}

type Record struct {
	Id         string        `json:"id"`
	ZoneId     string        `json:"zone_id"`
	ZoneName   string        `json:"zone_name"`
	Name       string        `json:"name"`
	Type       string        `json:"type"`
	Content    string        `json:"content"`
	Proxiable  bool          `json:"proxiable"`
	Proxied    bool          `json:"proxied"`
	Ttl        int           `json:"ttl"`
	Locked     bool          `json:"locked"`
	Meta       RecordMeta    `json:"meta"`
	Comment    interface{}   `json:"comment"`
	Tags       []interface{} `json:"tags"`
	CreatedOn  time.Time     `json:"created_on"`
	ModifiedOn time.Time     `json:"modified_on"`
}

type RecordMeta struct {
	AutoAdded           bool   `json:"auto_added"`
	ManagedByApps       bool   `json:"managed_by_apps"`
	ManagedByArgoTunnel bool   `json:"managed_by_argo_tunnel"`
	Source              string `json:"source"`
}

type Plan struct {
	Id                string `json:"id"`
	Name              string `json:"name"`
	Price             int    `json:"price"`
	Currency          string `json:"currency"`
	Frequency         string `json:"frequency"`
	IsSubscribed      bool   `json:"is_subscribed"`
	CanSubscribe      bool   `json:"can_subscribe"`
	LegacyId          string `json:"legacy_id"`
	LegacyDiscount    bool   `json:"legacy_discount"`
	ExternallyManaged bool   `json:"externally_managed"`
}

type ResultInfo struct {
	Page       int `json:"page"`
	PerPage    int `json:"per_page"`
	TotalPages int `json:"total_pages"`
	Count      int `json:"count"`
	TotalCount int `json:"total_count"`
}
