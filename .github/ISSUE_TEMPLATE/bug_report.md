---
name: 🐛 Bug Report
about: If something isn't working as expected or documented
labels: type/bug,status/needs-triage

---

<!--- Please keep this note for the community --->

### Community Note

* Please vote on this issue by adding a 👍 [reaction](https://blog.github.com/2016-03-10-add-reactions-to-pull-requests-issues-and-comments/) to the original issue to help the community and maintainers prioritize this request
* Please do not leave "+1" or other comments that do not add relevant new information or questions, they generate extra noise for issue followers and do not help prioritize the request
* If you are interested in working on this issue or have submitted a pull request, please leave a comment

<!--- Thank you for keeping this note for the community --->

Thank you for opening an issue. Please note that we try to keep the Terraform issue tracker reserved for bug reports and feature requests. For general usage questions, please see: https://www.terraform.io/community.html.

### PingOne Terraform provider Version
<!--- Check the version you have configured in your .tf files. If you are not running the latest version of the provider, please upgrade because your issue may have already been fixed. -->

### Terraform Version
<!--- Run `terraform -v` to show the version. If you are not running the latest version of Terraform, please upgrade because your issue may have already been fixed. -->

### Affected Resource(s)
<!--- Please list the resources as a list, for example: -->
- pingone_environment
- pingone_population

<!--- If this issue appears to affect multiple resources, it may be an issue with Terraform's core, so please mention this. -->

### Terraform Configuration Files
```hcl
# Copy-paste your PingOne related Terraform configurations here - for large Terraform configs,
# please use a service like Dropbox and share a link to the ZIP file. For
# security, you can also encrypt the files using our GPG public key.

# Remember to replace any account/customer sensitive information in the configuration before submitting the issue
```

### Debug Output
<!--- Please provide your debug output with `TF_LOG=DEBUG` enabled on your `terraform plan` or `terraform apply` -->

### Panic Output
<!--- If Terraform produced a panic, please provide your debug output from the GO panic -->

### Expected Behavior
<!--- What should have happened? -->

### Actual Behavior
<!--- What actually happened? -->

### Steps to Reproduce
<!---Please list the steps required to reproduce the issue, for example: -->
1. `terraform apply`

### Important Factoids
<!--- Are there anything you'd like to share about the general setup of your PingOne account?  Please do not include sensitive information or account data -->

### References
<!--- Are there any other GitHub issues (open or closed) or Pull Requests that should be linked here? For example: -->
- GH-1234
