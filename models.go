package westcn

import (
	"fmt"
	"strings"
	"time"

	"github.com/libdns/libdns"
)

type APIResponse[T any] struct {
	Result    int    `json:"result,omitempty"`
	ClientID  string `json:"clientid,omitempty"`
	Message   string `json:"msg,omitempty"`
	ErrorCode int    `json:"errcode,omitempty"`
	Data      T      `json:"data,omitempty"`
}

func (a APIResponse[T]) Error() string {
	return fmt.Sprintf("%d: %s (%d)", a.ErrorCode, a.Message, a.Result)
}

type RecordID struct {
	ID int `json:"id,omitempty"`
}

type Records struct {
	Records []Record `json:"items,omitempty"`
}

type Record struct {
	ID     int    `json:"id,omitempty"`
	Domain string `url:"domain,omitempty"` // 域名，删除记录时需要
	Item   string `url:"item,omitempty"`   // 主机名称，获取记录时返回
	Host   string `url:"host,omitempty"`   // 主机名称，添加记录时需要
	Type   string `url:"type,omitempty"`
	Value  string `url:"value,omitempty"`
	TTL    int    `url:"ttl,omitempty"` // 60~86400 seconds
	Level  int    `url:"level,omitempty"`
}

func (r Record) libdnsRecord(zone string) (libdns.Record, error) {
	// format host name
	if r.Host == "" {
		r.Host = r.Item
	}
	// format MX record
	if r.Type == "MX" {
		r.Value = fmt.Sprintf("%d %s", r.Level, r.Value)
	}
	return libdns.RR{
		Name: libdns.RelativeName(r.Host, zone),
		TTL:  time.Duration(r.TTL) * time.Second,
		Type: r.Type,
		Data: r.Value,
	}.Parse()
}

func westcnRecord(zone string, r libdns.Record) (Record, error) {
	rr := r.RR()
	westcnRec := Record{
		Domain: strings.TrimSuffix(zone, "."),
		Host:   rr.Name,
		Type:   rr.Type,
		TTL:    int(rr.TTL.Seconds()),
		Value:  rr.Data,
	}
	// Set default values for TTL and Level if not provided
	if westcnRec.TTL <= 0 {
		westcnRec.TTL = 600
	}
	if westcnRec.Level <= 0 {
		westcnRec.Level = 10
	}
	return westcnRec, nil
}
