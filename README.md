# ddns

A small, personal project used to enable DDNS via CloudFlare.

## Goals:

- Interval-based update
- Start/stop for timer via API
- Current set public IP address via API (e.g. `/v1/ddns/ip`)

## Example Config

```json
{
  "api_token": "",
  "private_hostname": "",
  "public_hostname": "",
  "timeout_seconds": 3,
  "service_port": 8000,
  "zone_id": ""
}
```