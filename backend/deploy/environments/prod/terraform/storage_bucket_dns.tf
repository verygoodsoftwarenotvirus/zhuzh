# CNAME records for domain-named GCS buckets (HTTP website hosting).
# Points media and userdata subdomains to GCS.
# proxied = true: Cloudflare terminates TLS and proxies to GCS over HTTP.
resource "cloudflare_dns_record" "media_storage_bucket" {
  zone_id = var.CLOUDFLARE_ZONE_ID
  name    = "media"
  type    = "CNAME"
  content = "c.storage.googleapis.com"
  proxied = true
  ttl     = 1
  comment = "GCS bucket ${local.media_domain} (HTTPS via Cloudflare)"
}

resource "cloudflare_dns_record" "userdata_storage_bucket" {
  zone_id = var.CLOUDFLARE_ZONE_ID
  name    = "userdata"
  type    = "CNAME"
  content = "c.storage.googleapis.com"
  proxied = true
  ttl     = 1
  comment = "GCS bucket ${local.userdata_domain} (HTTPS via Cloudflare)"
}
