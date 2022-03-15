package internal

import (
	"context"
	"github.com/nicseltzer/ddns/pkg/cloudflare"
)

type dnsRecord struct {
	id       string
	variant  string
	hostname string
	content  string
	zoneID   string
}

func (d *dnsRecord) Id() string {
	return d.id
}

func (d *dnsRecord) Variant() string {
	return d.variant
}

func (d *dnsRecord) Hostname() string {
	return d.hostname
}

func (d *dnsRecord) Content() string {
	return d.content
}

func (d *dnsRecord) ZoneID() string {
	return d.zoneID
}

// cfRecordToRecord converts a CloudFlare response record into a local domain struct
func cfRecordToRecord(record *cloudflare.Record) dnsRecord {
	if record == nil {
		return dnsRecord{}
	}

	return dnsRecord{
		id:       record.Id,
		variant:  record.Type,
		hostname: record.Name,
		content:  record.Content,
		zoneID:   record.ZoneId,
	}
}

// listPublicDNSRecords is a small helper to list all DNS records with a configured hostname
func (s *service) listPublicDNSRecord(ctx context.Context, variant string) dnsRecord {
	records := s.listDNSRecords(ctx, variant, s.config.PublicHostname)
	if len(records) == 0 {
		return dnsRecord{}
	}
	return records[0]
}

// listPrivateDNSRecords is a small helper to list all DNS records with a configured hostname
func (s *service) listPrivateDNSRecord(ctx context.Context, variant string) dnsRecord {
	records := s.listDNSRecords(ctx, variant, s.config.PrivateHostname)
	if len(records) == 0 {
		return dnsRecord{}
	}
	return records[0]
}

// listDNSRecords returns a list of DNS records associated with a given hostname-variant mapping
func (s *service) listDNSRecords(ctx context.Context, variant, hostname string) []dnsRecord {
	cf := cloudflare.NewClient(s.config.APIToken, s.config.Timeout)
	request := cloudflare.ListRecordRequest{
		Type: variant,
		Name: hostname,
	}
	response := cf.List(ctx, request, s.config.ZoneID)
	var records []dnsRecord
	for _, record := range response {
		cfRecordToRecord(record)
	}
	return records
}

// updateDNSRecord takes and maps a hostname to an IP for a given record ID
func (s *service) updateDNSRecord(ctx context.Context, record dnsRecord) {
	cf := cloudflare.NewClient(s.config.APIToken, s.config.Timeout)
	request := cloudflare.UpdateRecordRequest{
		Type:    record.Variant(),
		Name:    record.Hostname(),
		Content: record.Content(),
		TTL:     0,
		Proxied: false,
	}
	cf.Update(ctx, request, s.config.ZoneID, record.Id())
}

func (s *service) createPrivateDNSRecord(ctx context.Context) {
}

func (s *service) createDNSRecord(ctx context.Context, record dnsRecord) {
	cf := cloudflare.NewClient(s.config.APIToken, s.config.Timeout)
	request := cloudflare.CreateRecordRequest{
		Type:     record.Variant(),
		Name:     record.Hostname(),
		Content:  record.Content(),
		TTL:      0,
		Priority: 0,
		Proxied:  false,
	}

}
