// Package westcn implements a DNS record management client compatible
// with the libdns interfaces for west.cn.
package westcn

import (
	"context"
	"fmt"
	"strings"

	"github.com/libdns/libdns"
)

// Provider facilitates DNS record manipulation with west.cn.
type Provider struct {
	// Username is your username for west.cn, see https://www.west.cn/CustomerCenter/doc/apiv2.html#12u3001u8eabu4efdu9a8cu8bc10a3ca20id3d12u3001u8eabu4efdu9a8cu8bc13e203ca3e
	Username string `json:"username,omitempty"`
	// APIPassword is your API password for west.cn, see https://www.west.cn/CustomerCenter/doc/apiv2.html#12u3001u8eabu4efdu9a8cu8bc10a3ca20id3d12u3001u8eabu4efdu9a8cu8bc13e203ca3e
	APIPassword string `json:"api_password,omitempty"`
}

// GetRecords lists all the records in the zone.
func (p *Provider) GetRecords(ctx context.Context, zone string) ([]libdns.Record, error) {
	client, err := p.getClient()
	if err != nil {
		return nil, err
	}

	records, err := client.GetRecords(ctx, strings.TrimSuffix(zone, "."))
	if err != nil {
		return nil, err
	}

	results := make([]libdns.Record, 0, len(records))
	for _, rec := range records {
		libdnsRec, err := rec.libdnsRecord(zone)
		if err != nil {
			return nil, fmt.Errorf("parsing West.cn DNS record %+v: %v", rec, err)
		}
		results = append(results, libdnsRec)
	}

	return results, nil
}

// AppendRecords adds records to the zone. It returns the records that were added.
func (p *Provider) AppendRecords(ctx context.Context, zone string, records []libdns.Record) ([]libdns.Record, error) {
	client, err := p.getClient()
	if err != nil {
		return nil, err
	}

	var results []libdns.Record
	for _, rec := range records {
		westcnRec, err := westcnRecord(zone, rec)
		if err != nil {
			return nil, fmt.Errorf("parsing libdns record %+v: %v", rec, err)
		}
		if _, err = client.AddRecord(ctx, westcnRec); err != nil {
			return nil, err
		}
		libdnsRec, err := westcnRec.libdnsRecord(zone)
		if err != nil {
			return nil, fmt.Errorf("parsing West.cn DNS record %+v: %v", rec, err)
		}
		results = append(results, libdnsRec)
	}

	return results, nil
}

// SetRecords sets the records in the zone, either by updating existing records or creating new ones.
// It returns the updated records.
func (p *Provider) SetRecords(ctx context.Context, zone string, records []libdns.Record) ([]libdns.Record, error) {
	client, err := p.getClient()
	if err != nil {
		return nil, err
	}

	var results []libdns.Record
	for _, rec := range records {
		// West.cn does not support updating records, so we need to delete the existing record first
		rr := rec.RR()
		id, err := p.getRecordId(ctx, zone, rr.Name, rr.Type, rr.Data)
		if err == nil {
			if err = client.DeleteRecord(ctx, strings.TrimSuffix(zone, "."), id); err != nil {
				return nil, err
			}
		}
		// Now we can add the record
		westcnRec, err := westcnRecord(zone, rec)
		if err != nil {
			return nil, fmt.Errorf("parsing libdns record %+v: %v", rec, err)
		}
		_, err = client.AddRecord(ctx, westcnRec)
		if err != nil {
			return nil, err
		}
		// Add the record to the results
		libdnsRec, err := westcnRec.libdnsRecord(zone)
		if err != nil {
			return nil, fmt.Errorf("parsing West.cn DNS record %+v: %v", rec, err)
		}
		results = append(results, libdnsRec)
	}

	return results, nil
}

// DeleteRecords deletes the records from the zone. It returns the records that were deleted.
func (p *Provider) DeleteRecords(ctx context.Context, zone string, records []libdns.Record) ([]libdns.Record, error) {
	client, err := p.getClient()
	if err != nil {
		return nil, err
	}

	for _, record := range records {
		rr := record.RR()
		id, err := p.getRecordId(ctx, zone, rr.Name, rr.Type, rr.Data)
		if err != nil {
			return nil, err
		}
		if err = client.DeleteRecord(ctx, strings.TrimSuffix(zone, "."), id); err != nil {
			return nil, err
		}
	}

	return records, nil
}

func (p *Provider) getRecordId(ctx context.Context, zone, recName, recType string, recVal ...string) (int, error) {
	client, err := p.getClient()
	if err != nil {
		return 0, err
	}

	records, err := client.GetRecords(ctx, strings.TrimSuffix(zone, "."))
	if err != nil {
		return 0, err
	}

	for _, rec := range records {
		if recName == rec.Item && recType == rec.Type {
			if len(recVal) > 0 && recVal[0] != "" && rec.Value != recVal[0] {
				continue
			}
			return rec.ID, nil
		}
	}

	return 0, fmt.Errorf("record %q not found", recName)
}

func (p *Provider) getClient() (*Client, error) {
	return NewClient(p.Username, p.APIPassword)
}

// Interface guards
var (
	_ libdns.RecordGetter   = (*Provider)(nil)
	_ libdns.RecordAppender = (*Provider)(nil)
	_ libdns.RecordSetter   = (*Provider)(nil)
	_ libdns.RecordDeleter  = (*Provider)(nil)
)
