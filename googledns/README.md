# googledns DNS-01 Exec Plugin

## Usage

## Configuration

The `googledns` plugin requires a Google Cloud service account, which contains all the information required to connect to and manage domains on Google Cloud DNS.

### Creating DNS-01 TXT Records

```
cat googledns.json | \
  APIVERSION="v1" \
  COMMAND="CREATE" \
  DOMAIN="hightowerlabs.com" \
  FQDN="_acme-challenge.hightowerlabs.com." \
  TOKEN="8bGFl9SNhZzukcwdR7e52gFwq6HaEHB43LbimZQwnLg" \
  googledns
```

### Deleting DNS-01 TXT Records

```
cat googledns.json | \
  APIVERSION="v1" \
  COMMAND="DELETE" \
  DOMAIN="hightowerlabs.com" \
  FQDN="_acme-challenge.hightowerlabs.com." \
  TOKEN="8bGFl9SNhZzukcwdR7e52gFwq6HaEHB43LbimZQwnLg" \
  googledns
```
