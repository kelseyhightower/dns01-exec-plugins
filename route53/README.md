# route53 DNS-01 Exec Plugin

## Usage

## Configuration

The `route53` plugin requires a ZONEID as environment variable.

### Creating DNS-01 TXT Records

```
cat dns01.json | \
  APIVERSION="v1" \
  COMMAND="CREATE" \
  DOMAIN="hightowerlabs.com" \
  FQDN="_acme-challenge.hightowerlabs.com." \
  TOKEN="8bGFl9SNhZzukcwdR7e52gFwq6HaEHB43LbimZQwnLg" \
  ZONEID="Z23BO2W3S4CRS4" \
  route53
```

### Deleting DNS-01 TXT Records

```
cat dns01.json | \
  APIVERSION="v1" \
  COMMAND="DELETE" \
  DOMAIN="hightowerlabs.com" \
  FQDN="_acme-challenge.hightowerlabs.com." \
  TOKEN="8bGFl9SNhZzukcwdR7e52gFwq6HaEHB43LbimZQwnLg" \
  ZONEID="Z23BO2W3S4CRS4" \
  route53
```
