# Log in to PingOne with Ping CLI first, storing the authentication session in the OS keychain:
pingcli pingone auth login

# Then point the provider at the Ping CLI profile (or use the provider
# attributes `config_profile` / `config_grant_type` directly in HCL):
export PINGCLI_PROFILE="myprofile"
# Optional: pin the grant type when the profile holds sessions for more than one
# (client_credentials, authorization_code or device_code):
# export PINGCLI_PINGONE_AUTHENTICATION_TYPE="device_code"
export PINGONE_REGION_CODE="NA"
