# Cloudflare DNS-01 Exec Plugin

## Usage

## Configuration

The `cloudflare` plugin requires a pre-existing Cloudflare account. The plugin expects the account email and API access token to be passed in the format below:

```
$ cat cloudflare.json
{
    "email": "test@test.com",
    "key": "172yeqdaysdau2ueygasd287ed8gd8asdy7ds"
}
```

### Creating DNS-01 TXT Records

```
$ cat cloudflare.json | \
  APIVERSION="v1" \
  COMMAND="CREATE" \
  DOMAIN="hightowerlabs.com" \
  FQDN="_acme-challenge.hightowerlabs.com." \
  TOKEN="8bGFl9SNhZzukcwdR7e52gFwq6HaEHB43LbimZQwnLg" \
  cloudflare
```

### Deleting DNS-01 TXT Records

```
$ cat cloudflare.json | \
  APIVERSION="v1" \
  COMMAND="DELETE" \
  DOMAIN="hightowerlabs.com" \
  FQDN="_acme-challenge.hightowerlabs.com." \
  TOKEN="8bGFl9SNhZzukcwdR7e52gFwq6HaEHB43LbimZQwnLg" \
  cloudflare
```
