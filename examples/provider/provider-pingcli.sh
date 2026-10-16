# Log in to PingOne with Ping CLI first, storing the authentication session in the OS keychain:
pingcli pingone auth login

# Then enable Ping CLI authentication (or set `use_pingcli_auth` directly in
# HCL):
export PINGONE_REGION_CODE="NA"
