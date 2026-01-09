// Copyright © 2025 Ping Identity Corporation

package mfa

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/patrickcping/pingone-go-sdk-v2/mfa"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest/legacysdk"
)

func MFADevicePolicyDefault_CheckDestroy(s *terraform.State) error {
	var ctx = context.Background()

	p1Client, err := legacysdk.TestClient(ctx)

	if err != nil {
		return err
	}

	apiClient := p1Client.API.MFAAPIClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pingone_mfa_device_policy_default" {
			continue
		}

		shouldContinue, err := legacysdk.CheckParentEnvironmentDestroy(ctx, p1Client.API.ManagementAPIClient, rs.Primary.Attributes["environment_id"])
		if err != nil {
			return err
		}

		if shouldContinue {
			continue
		}

		policy, _, err := apiClient.DeviceAuthenticationPolicyApi.ReadOneDeviceAuthenticationPolicy(ctx, rs.Primary.Attributes["environment_id"], rs.Primary.ID).Execute()
		if err != nil {
			return fmt.Errorf("Unable to retrieve MFA device policy: %s", err)
		}

		if policy.Authentication.GetDeviceSelection() != mfa.ENUMMFADEVICEPOLICYSELECTION_DEFAULT_TO_FIRST {
			return fmt.Errorf("Expected Authentication.DeviceSelection to be DEFAULT_TO_FIRST, got %s", policy.Authentication.GetDeviceSelection())
		}

		if policy.GetNewDeviceNotification() != mfa.ENUMMFADEVICEPOLICYNEWDEVICENOTIFICATION_NONE {
			return fmt.Errorf("Expected NewDeviceNotification to be NONE, got %s", policy.GetNewDeviceNotification())
		}

		if policy.NotificationsPolicy != nil {
			return fmt.Errorf("Expected NotificationsPolicy to be nil")
		}

		// SMS
		if policy.Sms.GetEnabled() {
			return fmt.Errorf("Expected SMS to be disabled")
		}

		if policy.Sms.GetPairingDisabled() {
			return fmt.Errorf("Expected SMS PairingDisabled to be false")
		}

		if policy.Sms.GetPromptForNicknameOnPairing() {
			return fmt.Errorf("Expected SMS PromptForNicknameOnPairing to be false")
		}

		if v := policy.Sms.Otp.LifeTime.GetDuration(); v != 30 {
			return fmt.Errorf("Expected SMS Otp.LifeTime.Duration to be 30, got %d", v)
		}

		if v := policy.Sms.Otp.LifeTime.GetTimeUnit(); v != mfa.ENUMTIMEUNIT_MINUTES {
			return fmt.Errorf("Expected SMS Otp.LifeTime.TimeUnit to be MINUTES, got %s", v)
		}

		if v := policy.Sms.Otp.Failure.GetCount(); v != 3 {
			return fmt.Errorf("Expected SMS Otp.Failure.Count to be 3, got %d", v)
		}

		if v := policy.Sms.Otp.Failure.CoolDown.GetDuration(); v != 2 {
			return fmt.Errorf("Expected SMS Otp.Failure.CoolDown.Duration to be 2, got %d", v)
		}

		if v := policy.Sms.Otp.Failure.CoolDown.GetTimeUnit(); v != mfa.ENUMTIMEUNIT_MINUTES {
			return fmt.Errorf("Expected SMS Otp.Failure.CoolDown.TimeUnit to be MINUTES, got %s", v)
		}

		if v := policy.Sms.Otp.GetOtpLength(); v != 6 {
			return fmt.Errorf("Expected SMS Otp.OtpLength to be 6, got %d", v)
		}

		// Voice
		if policy.Voice.GetEnabled() {
			return fmt.Errorf("Expected Voice to be disabled")
		}

		if policy.Voice.GetPairingDisabled() {
			return fmt.Errorf("Expected Voice PairingDisabled to be false")
		}

		if policy.Voice.GetPromptForNicknameOnPairing() {
			return fmt.Errorf("Expected Voice PromptForNicknameOnPairing to be false")
		}

		if v := policy.Voice.Otp.LifeTime.GetDuration(); v != 30 {
			return fmt.Errorf("Expected Voice Otp.LifeTime.Duration to be 30, got %d", v)
		}

		if v := policy.Voice.Otp.LifeTime.GetTimeUnit(); v != mfa.ENUMTIMEUNIT_MINUTES {
			return fmt.Errorf("Expected Voice Otp.LifeTime.TimeUnit to be MINUTES, got %s", v)
		}

		if v := policy.Voice.Otp.Failure.GetCount(); v != 3 {
			return fmt.Errorf("Expected Voice Otp.Failure.Count to be 3, got %d", v)
		}

		if v := policy.Voice.Otp.Failure.CoolDown.GetDuration(); v != 2 {
			return fmt.Errorf("Expected Voice Otp.Failure.CoolDown.Duration to be 2, got %d", v)
		}

		if v := policy.Voice.Otp.Failure.CoolDown.GetTimeUnit(); v != mfa.ENUMTIMEUNIT_MINUTES {
			return fmt.Errorf("Expected Voice Otp.Failure.CoolDown.TimeUnit to be MINUTES, got %s", v)
		}

		if v := policy.Voice.Otp.GetOtpLength(); v != 6 {
			return fmt.Errorf("Expected Voice Otp.OtpLength to be 6, got %d", v)
		}

		// Email
		if policy.Email.GetEnabled() {
			return fmt.Errorf("Expected Email to be disabled")
		}

		if policy.Email.GetPairingDisabled() {
			return fmt.Errorf("Expected Email PairingDisabled to be false")
		}

		if policy.Email.GetPromptForNicknameOnPairing() {
			return fmt.Errorf("Expected Email PromptForNicknameOnPairing to be false")
		}

		if v := policy.Email.Otp.LifeTime.GetDuration(); v != 30 {
			return fmt.Errorf("Expected Email Otp.LifeTime.Duration to be 30, got %d", v)
		}

		if v := policy.Email.Otp.LifeTime.GetTimeUnit(); v != mfa.ENUMTIMEUNIT_MINUTES {
			return fmt.Errorf("Expected Email Otp.LifeTime.TimeUnit to be MINUTES, got %s", v)
		}

		if v := policy.Email.Otp.Failure.GetCount(); v != 3 {
			return fmt.Errorf("Expected Email Otp.Failure.Count to be 3, got %d", v)
		}

		if v := policy.Email.Otp.Failure.CoolDown.GetDuration(); v != 2 {
			return fmt.Errorf("Expected Email Otp.Failure.CoolDown.Duration to be 2, got %d", v)
		}

		if v := policy.Email.Otp.Failure.CoolDown.GetTimeUnit(); v != mfa.ENUMTIMEUNIT_MINUTES {
			return fmt.Errorf("Expected Email Otp.Failure.CoolDown.TimeUnit to be MINUTES, got %s", v)
		}

		if v := policy.Email.Otp.GetOtpLength(); v != 6 {
			return fmt.Errorf("Expected Email Otp.OtpLength to be 6, got %d", v)
		}

		// Mobile
		if !policy.Mobile.GetEnabled() {
			return fmt.Errorf("Expected Mobile to be enabled")
		}

		// Mobile Applications
		if len(policy.Mobile.GetApplications()) != 1 {
			return fmt.Errorf("Expected 1 Mobile Application, got %d", len(policy.Mobile.GetApplications()))
		}

		for _, app := range policy.Mobile.GetApplications() {
			if app.GetType() == "pingIdAppConfig" {
				if !app.Push.GetEnabled() {
					return fmt.Errorf("Expected Mobile Application Push to be enabled")
				}
				if app.Push.NumberMatching.GetEnabled() {
					return fmt.Errorf("Expected Mobile Application Push NumberMatching to be disabled")
				}
				if app.Otp.GetEnabled() {
					return fmt.Errorf("Expected Mobile Application OTP to be disabled")
				}
				if app.AutoEnrollment.GetEnabled() {
					return fmt.Errorf("Expected Mobile Application AutoEnrollment to be disabled")
				}
				if app.DeviceAuthorization.GetEnabled() {
					return fmt.Errorf("Expected Mobile Application DeviceAuthorization to be disabled")
				}
				if v := app.PushTimeout.GetDuration(); v != 100 {
					return fmt.Errorf("Expected Mobile Application PushTimeout.Duration to be 100, got %d", v)
				}
				if v := app.PushTimeout.GetTimeUnit(); v != mfa.ENUMTIMEUNITPUSHTIMEOUT_SECONDS {
					return fmt.Errorf("Expected Mobile Application PushTimeout.TimeUnit to be SECONDS, got %s", v)
				}
				if v := app.PairingKeyLifetime.GetDuration(); v != 48 {
					return fmt.Errorf("Expected Mobile Application PairingKeyLifetime.Duration to be 48, got %d", v)
				}
				if v := app.PairingKeyLifetime.GetTimeUnit(); v != mfa.ENUMTIMEUNITPAIRINGKEYLIFETIME_HOURS {
					return fmt.Errorf("Expected Mobile Application PairingKeyLifetime.TimeUnit to be HOURS, got %s", v)
				}
				if v := app.PushLimit.GetCount(); v != 5 {
					return fmt.Errorf("Expected Mobile Application PushLimit.Count to be 5, got %d", v)
				}
				if v := app.PushLimit.TimePeriod.GetDuration(); v != 10 {
					return fmt.Errorf("Expected Mobile Application PushLimit.TimePeriod.Duration to be 10, got %d", v)
				}
				if v := app.PushLimit.TimePeriod.GetTimeUnit(); v != mfa.ENUMTIMEUNIT_MINUTES {
					return fmt.Errorf("Expected Mobile Application PushLimit.TimePeriod.TimeUnit to be MINUTES, got %s", v)
				}
				if v := app.PushLimit.LockDuration.GetDuration(); v != 30 {
					return fmt.Errorf("Expected Mobile Application PushLimit.LockDuration.Duration to be 30, got %d", v)
				}
				if v := app.PushLimit.LockDuration.GetTimeUnit(); v != mfa.ENUMTIMEUNIT_MINUTES {
					return fmt.Errorf("Expected Mobile Application PushLimit.LockDuration.TimeUnit to be MINUTES, got %s", v)
				}
				if app.GetPairingDisabled() {
					return fmt.Errorf("Expected Mobile Application PairingDisabled to be false")
				}
				if v := app.NewRequestDurationConfiguration.DeviceTimeout.GetDuration(); v != 25 {
					return fmt.Errorf("Expected Mobile Application NewRequestDurationConfiguration.DeviceTimeout.Duration to be 25, got %d", v)
				}
				if v := app.NewRequestDurationConfiguration.DeviceTimeout.GetTimeUnit(); v != mfa.ENUMTIMEUNITSECONDS_SECONDS {
					return fmt.Errorf("Expected Mobile Application NewRequestDurationConfiguration.DeviceTimeout.TimeUnit to be SECONDS, got %s", v)
				}
				if v := app.NewRequestDurationConfiguration.TotalTimeout.GetDuration(); v != 40 {
					return fmt.Errorf("Expected Mobile Application NewRequestDurationConfiguration.TotalTimeout.Duration to be 40, got %d", v)
				}
				if v := app.NewRequestDurationConfiguration.TotalTimeout.GetTimeUnit(); v != mfa.ENUMTIMEUNITSECONDS_SECONDS {
					return fmt.Errorf("Expected Mobile Application NewRequestDurationConfiguration.TotalTimeout.TimeUnit to be SECONDS, got %s", v)
				}
				if !app.IpPairingConfiguration.GetAnyIPAdress() {
					return fmt.Errorf("Expected Mobile Application IpPairingConfiguration.AnyIPAdress to be true")
				}
				if app.GetBiometricsEnabled() {
					return fmt.Errorf("Expected Mobile Application BiometricsEnabled to be false")
				}
			}
		}

		if policy.Mobile.GetPromptForNicknameOnPairing() {
			return fmt.Errorf("Expected Mobile PromptForNicknameOnPairing to be false")
		}

		if v := policy.Mobile.Otp.Failure.GetCount(); v != 3 {
			return fmt.Errorf("Expected Mobile Otp.Failure.Count to be 3, got %d", v)
		}

		if v := policy.Mobile.Otp.Failure.CoolDown.GetDuration(); v != 2 {
			return fmt.Errorf("Expected Mobile Otp.Failure.CoolDown.Duration to be 2, got %d", v)
		}

		if v := policy.Mobile.Otp.Failure.CoolDown.GetTimeUnit(); v != mfa.ENUMTIMEUNIT_MINUTES {
			return fmt.Errorf("Expected Mobile Otp.Failure.CoolDown.TimeUnit to be MINUTES, got %s", v)
		}

		// TOTP
		if policy.Totp.GetEnabled() {
			return fmt.Errorf("Expected TOTP to be disabled")
		}

		if policy.Totp.GetPairingDisabled() {
			return fmt.Errorf("Expected TOTP PairingDisabled to be false")
		}

		if v := policy.Totp.Otp.Failure.GetCount(); v != 3 {
			return fmt.Errorf("Expected TOTP Otp.Failure.Count to be 3, got %d", v)
		}

		if v := policy.Totp.Otp.Failure.CoolDown.GetDuration(); v != 2 {
			return fmt.Errorf("Expected TOTP Otp.Failure.CoolDown.Duration to be 2, got %d", v)
		}

		if v := policy.Totp.Otp.Failure.CoolDown.GetTimeUnit(); v != mfa.ENUMTIMEUNIT_MINUTES {
			return fmt.Errorf("Expected TOTP Otp.Failure.CoolDown.TimeUnit to be MINUTES, got %s", v)
		}

		// FIDO2
		if policy.Fido2.GetEnabled() {
			return fmt.Errorf("Expected FIDO2 to be disabled")
		}

		if policy.Fido2.GetPairingDisabled() {
			return fmt.Errorf("Expected FIDO2 PairingDisabled to be false")
		}

		if policy.Fido2.GetPromptForNicknameOnPairing() {
			return fmt.Errorf("Expected FIDO2 PromptForNicknameOnPairing to be false")
		}

		// OATH Token
		if policy.OathToken.GetEnabled() {
			return fmt.Errorf("Expected OATH Token to be disabled")
		}

		if policy.OathToken.GetPairingDisabled() {
			return fmt.Errorf("Expected OATH Token PairingDisabled to be false")
		}

		if v := policy.OathToken.Otp.Failure.GetCount(); v != 3 {
			return fmt.Errorf("Expected OATH Token Otp.Failure.Count to be 3, got %d", v)
		}

		if v := policy.OathToken.Otp.Failure.CoolDown.GetDuration(); v != 2 {
			return fmt.Errorf("Expected OATH Token Otp.Failure.CoolDown.Duration to be 2, got %d", v)
		}

		if v := policy.OathToken.Otp.Failure.CoolDown.GetTimeUnit(); v != mfa.ENUMTIMEUNIT_MINUTES {
			return fmt.Errorf("Expected OATH Token Otp.Failure.CoolDown.TimeUnit to be MINUTES, got %s", v)
		}

		// Desktop (PingID)
		if policy.Desktop != nil && policy.Desktop.GetEnabled() {
			return fmt.Errorf("Expected Desktop to be disabled")
		}

		if policy.Desktop.GetPairingDisabled() {
			return fmt.Errorf("Expected Desktop PairingDisabled to be false")
		}

		if v := policy.Desktop.Otp.Failure.GetCount(); v != 3 {
			return fmt.Errorf("Expected Desktop Otp.Failure.Count to be 3, got %d", v)
		}

		if v := policy.Desktop.Otp.Failure.CoolDown.GetDuration(); v != 2 {
			return fmt.Errorf("Expected Desktop Otp.Failure.CoolDown.Duration to be 2, got %d", v)
		}

		if v := policy.Desktop.Otp.Failure.CoolDown.GetTimeUnit(); v != mfa.ENUMTIMEUNIT_MINUTES {
			return fmt.Errorf("Expected Desktop Otp.Failure.CoolDown.TimeUnit to be MINUTES, got %s", v)
		}

		if v := policy.Desktop.PairingKeyLifetime.GetDuration(); v != 48 {
			return fmt.Errorf("Expected Desktop PairingKeyLifetime.Duration to be 48, got %d", v)
		}

		if v := policy.Desktop.PairingKeyLifetime.GetTimeUnit(); v != mfa.ENUMTIMEUNITPAIRINGKEYLIFETIME_HOURS {
			return fmt.Errorf("Expected Desktop PairingKeyLifetime.TimeUnit to be HOURS, got %s", v)
		}

		// Yubikey (PingID)
		if policy.Yubikey != nil && policy.Yubikey.GetEnabled() {
			return fmt.Errorf("Expected Yubikey to be disabled")
		}

		if policy.Yubikey.GetPairingDisabled() {
			return fmt.Errorf("Expected Yubikey PairingDisabled to be false")
		}

		if v := policy.Yubikey.Otp.Failure.GetCount(); v != 3 {
			return fmt.Errorf("Expected Yubikey Otp.Failure.Count to be 3, got %d", v)
		}

		if v := policy.Yubikey.Otp.Failure.CoolDown.GetDuration(); v != 2 {
			return fmt.Errorf("Expected Yubikey Otp.Failure.CoolDown.Duration to be 2, got %d", v)
		}

		if v := policy.Yubikey.Otp.Failure.CoolDown.GetTimeUnit(); v != mfa.ENUMTIMEUNIT_MINUTES {
			return fmt.Errorf("Expected Yubikey Otp.Failure.CoolDown.TimeUnit to be MINUTES, got %s", v)
		}

		// Remember Me
		if policy.RememberMe != nil && policy.RememberMe.Web.GetEnabled() {
			return fmt.Errorf("Expected RememberMe.Web to be disabled")
		}

		if v := policy.RememberMe.Web.LifeTime.GetDuration(); v != 30 {
			return fmt.Errorf("Expected RememberMe.Web.LifeTime.Duration to be 30, got %d", v)
		}

		if v := policy.RememberMe.Web.LifeTime.GetTimeUnit(); v != mfa.ENUMTIMEUNITREMEMBERMEWEBLIFETIME_MINUTES {
			return fmt.Errorf("Expected RememberMe.Web.LifeTime.TimeUnit to be MINUTES")
		}
	}

	return nil
}
