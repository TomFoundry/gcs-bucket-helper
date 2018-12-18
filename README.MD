# Athera GCS Bucket Helper

Athera GCS (Google Cloud Storage) Bucket Helper creates a bucket and service account in GCP (Google Cloud Platform), then connects the bucket to Athera.

## Build Command
go build -ldflags "-X main.googleAppID=<GCP API Credentials Client ID> -X main.googleAppPseudoSecret=<GCP API Credentials Client Secret> -X main.atheraAppID=<Athera API Credentials Client ID> -X main.voyagerURL=<Voyager URL> -X main.atheraAuthURL=<Athera Auth URL> -X main.atheraTokenURL=<Athera Token URL>"

N.B. Client Secret is not actually expected to be secret for installed apps
https://developers.google.com/identity/protocols/OAuth2InstalledApp