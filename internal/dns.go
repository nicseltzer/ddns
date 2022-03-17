package internal

import (
	"context"
	"fmt"
	"github.com/nicseltzer/ddns/pkg/cloudflare"
	"github.com/nicseltzer/ddns/pkg/ifconfigme"
)

type dnsRecord struct {
	id       string
	variant  string
	hostname string
	content  string
	zoneID   string
}

func newDnsRecord(variant, hostname, content string) *dnsRecord {
	return &dnsRecord{
		variant:  variant,
		hostname: hostname,
		content:  content,
	}
}

func (d *dnsRecord) ID() string {
	return d.id
}

func (d *dnsRecord) SetID(id string) {
	d.id = id
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

func (d *dnsRecord) SetContent(content string) {
	d.content = content
}

func (d *dnsRecord) ZoneID() string {
	return d.zoneID
}

func (s *service) Scan() {
	ctx, cancel := context.WithTimeout(context.Background(), s.config.Timeout)
	defer cancel()

	// Did the IP address change from what's in DNS?
	currentPrivateRecord, exists := s.listPrivateDNSRecord(ctx)
	if !exists {
		var err error
		err = s.createPrivateRecord(ctx)
		if err != nil {
			fmt.Println(err)
		}
		err = s.createPublicRecord(ctx)
		if err != nil {
			fmt.Println(err)
			return
		}
	} else {
		prev := currentPrivateRecord.Content()
		curr := s.fetchExternalIP(ctx)
		if prev != curr {
			fmt.Printf("ip addr change, updating: prev=%s,curr=%s", prev, curr)
			s.ipCache.SetIP(curr)
			err := s.updatePrivateRecord(ctx, currentPrivateRecord)
			if err != nil {
				return
			}
		}
	}
}

func (s *service) fetchExternalIP(ctx context.Context) string {
	client := ifconfigme.NewClient(s.config.Timeout)
	ip := client.ExternalIP(ctx)
	return ip
}

// cfRecordToRecord converts a CloudFlare response record into a local domain struct
func cfRecordToRecord(record *cloudflare.Record) *dnsRecord {
	if record == nil {
		return nil
	}
	return &dnsRecord{
		id:       record.ID,
		variant:  record.Type,
		hostname: record.Name,
		content:  record.Content,
		zoneID:   record.ZoneId,
	}
}

// listPublicDNSRecords is a helper to list all DNS records with a configured hostname
func (s *service) listPublicDNSRecord(ctx context.Context) (*dnsRecord, bool) {
	records := s.listDNSRecords(ctx, DNSCNAMERecordType, s.config.PublicHostname)
	if records == nil {
		return nil, false
	}
	return records[0], true
}

// listPrivateDNSRecords is a helper to list all DNS records with a configured hostname
func (s *service) listPrivateDNSRecord(ctx context.Context) (*dnsRecord, bool) {
	records := s.listDNSRecords(ctx, DNSARecordType, s.config.PrivateHostname)
	if records == nil {
		return nil, false
	}
	return records[0], true
}

// listDNSRecords returns a list of DNS records associated with a given hostname-variant mapping
func (s *service) listDNSRecords(ctx context.Context, variant, hostname string) []*dnsRecord {
	cf := cloudflare.NewClient(s.config.APIToken, s.config.Timeout)
	request := cloudflare.ListRecordRequest{
		Type: variant,
		Name: hostname,
	}
	response := cf.List(ctx, request, s.config.ZoneID)
	if len(response) == 0 {
		return nil
	}
	var records []*dnsRecord
	for _, record := range response {
		records = append(records, cfRecordToRecord(record))
	}
	return records
}

func (s *service) updatePrivateRecord(ctx context.Context, record *dnsRecord) error {
	record.SetContent(s.ipCache.IP())
	err := s.updateDNSRecord(ctx, record)
	if err != nil {
		return err
	}
	return nil
}

// updateDNSRecord takes and maps a hostname to an IP for a given record ID
func (s *service) updateDNSRecord(ctx context.Context, record *dnsRecord) error {
	cf := cloudflare.NewClient(s.config.APIToken, s.config.Timeout)
	request := cloudflare.UpdateRecordRequest{
		Type:    record.Variant(),
		Name:    record.Hostname(),
		Content: record.Content(),
		TTL:     0,
		Proxied: false,
	}
	_, err := cf.Update(ctx, request, s.config.ZoneID, record.ID())
	if err != nil {
		return err
	}
	return nil
}

func (s *service) createPrivateRecord(ctx context.Context) error {
	record := newDnsRecord(DNSARecordType, s.config.PrivateHostname, s.ipCache.IP())
	err := s.createDNSRecord(ctx, record)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) createPublicRecord(ctx context.Context) error {
	record := newDnsRecord(DNSCNAMERecordType, s.config.PublicHostname, s.config.PrivateHostname)
	err := s.createDNSRecord(ctx, record)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) createDNSRecord(ctx context.Context, record *dnsRecord) error {
	cf := cloudflare.NewClient(s.config.APIToken, s.config.Timeout)
	request := cloudflare.CreateRecordRequest{
		Type:     record.Variant(),
		Name:     record.Hostname(),
		Content:  record.Content(),
		TTL:      0,
		Priority: 0,
		Proxied:  false,
	}
	_, err := cf.Create(ctx, request, s.config.ZoneID)
	if err != nil {
		return err
	}

	return nil
}
