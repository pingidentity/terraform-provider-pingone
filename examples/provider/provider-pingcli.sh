# Log in to PingOne with Ping CLI first, storing the authentication session in the OS keychain:
pingcli pingone auth login

# Then enable Ping CLI authentication (or set `use_pingcli_auth = true` in the provider configuration):
export PINGONE_TF_USE_PINGCLI_AUTH=true
export PINGONE_REGION_CODE="NA"
