# ddns

A small, personal project used to enable DDNS via CloudFlare.

## Goals:

- Interval-based update
- Start/stop for timer via API
- Current set public IP address via API (e.g. `/v1/ip`)

## Example Config

```json
{
  "api_token": "",
  "public_hostname": "",
  "private_hostname": "",
  "timeout_seconds": 3,
  "service_port": 8000,
  "zone_id": ""
}
```