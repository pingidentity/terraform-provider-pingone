// Copyright © 2025 Ping Identity Corporation

package helpers

import (
	"github.com/patrickcping/pingone-go-sdk-v2/risk"
)

var (
	BootstrapPredictorValues = map[string]risk.RiskPredictor{
		defaultUserRiskBehavior().RiskPredictorUserRiskBehavior.CompactName:       defaultUserRiskBehavior(),
		defaultUserBasedRiskBehavior().RiskPredictorUserRiskBehavior.CompactName:  defaultUserBasedRiskBehavior(),
		defaultIpVelocityByUser().RiskPredictorVelocity.CompactName:               defaultIpVelocityByUser(),
		defaultUserVelocityByIp().RiskPredictorVelocity.CompactName:               defaultUserVelocityByIp(),
		defaultAnonymousNetwork().RiskPredictorAnonymousNetwork.CompactName:       defaultAnonymousNetwork(),
		defaultGeoVelocity().RiskPredictorGeovelocity.CompactName:                 defaultGeoVelocity(),
		defaultIpRisk().RiskPredictorIPReputation.CompactName:                     defaultIpRisk(),
		defaultNewDevice().RiskPredictorDevice.CompactName:                        defaultNewDevice(),
		defaultUserLocationAnomaly().RiskPredictorUserLocationAnomaly.CompactName: defaultUserLocationAnomaly(),
	}
)

func defaultUserRiskBehavior() risk.RiskPredictor {

	defaultWeight := 5
	defaultScore := 50

	return risk.RiskPredictor{
		RiskPredictorUserRiskBehavior: &risk.RiskPredictorUserRiskBehavior{
			Name:        "User Risk Behavior",
			CompactName: "userRiskBehavior",
			PredictionModel: risk.RiskPredictorUserRiskBehaviorAllOfPredictionModel{
				Name: risk.ENUMUSERRISKBEHAVIORRISKMODEL_LOGIN_ANOMALY_STATISTIC,
			},
			Type: risk.ENUMPREDICTORTYPE_USER_RISK_BEHAVIOR,
			Default: &risk.RiskPredictorCommonDefault{
				Weight:    int32(defaultWeight),
				Score:     risk.PtrInt32(int32(defaultScore)),
				Evaluated: risk.PtrBool(false),
			},
		},
	}
}

func defaultUserBasedRiskBehavior() risk.RiskPredictor {

	defaultWeight := 5
	defaultScore := 75

	return risk.RiskPredictor{
		RiskPredictorUserRiskBehavior: &risk.RiskPredictorUserRiskBehavior{
			Name:        "User-Based Risk Behavior",
			CompactName: "userBasedRiskBehavior",
			PredictionModel: risk.RiskPredictorUserRiskBehaviorAllOfPredictionModel{
				Name: risk.ENUMUSERRISKBEHAVIORRISKMODEL_POINTS,
			},
			Type: risk.ENUMPREDICTORTYPE_USER_RISK_BEHAVIOR,
			Default: &risk.RiskPredictorCommonDefault{
				Weight:    int32(defaultWeight),
				Score:     risk.PtrInt32(int32(defaultScore)),
				Evaluated: risk.PtrBool(false),
			},
		},
	}
}

func defaultIpVelocityByUser() risk.RiskPredictor {

	everyQuantity := 1
	everyMinSample := 5

	useMedium := 0.96
	useHigh := 0.995

	slidingWindowQuantity := 14
	slidingWindowMinSample := 3

	fallbackHigh := 30.0
	fallbackMedium := 20.0

	defaultWeight := 5
	defaultScore := 75

	return risk.RiskPredictor{
		RiskPredictorVelocity: &risk.RiskPredictorVelocity{
			Name:        "IP Velocity",
			CompactName: "ipVelocityByUser",
			Measure:     (*risk.EnumPredictorVelocityMeasure)(risk.PtrString("DISTINCT_COUNT")),
			Of:          risk.PtrString("${event.ip}"),
			By:          []string{"${event.user.id}"},
			Every: &risk.RiskPredictorVelocityAllOfEvery{
				Unit:      (*risk.EnumPredictorUnit)(risk.PtrString("HOUR")),
				Quantity:  risk.PtrInt32(int32(everyQuantity)),
				MinSample: risk.PtrInt32(int32(everyMinSample)),
			},
			Use: &risk.RiskPredictorVelocityAllOfUse{
				Type:   (*risk.EnumPredictorVelocityUseType)(risk.PtrString("POISSON_WITH_MAX")),
				Medium: risk.PtrFloat32(float32(useMedium)),
				High:   risk.PtrFloat32(float32(useHigh)),
			},
			SlidingWindow: &risk.RiskPredictorVelocityAllOfSlidingWindow{
				Unit:      (*risk.EnumPredictorUnit)(risk.PtrString("DAY")),
				Quantity:  risk.PtrInt32(int32(slidingWindowQuantity)),
				MinSample: risk.PtrInt32(int32(slidingWindowMinSample)),
			},
			Fallback: &risk.RiskPredictorVelocityAllOfFallback{
				Strategy: (*risk.EnumPredictorVelocityFallbackStrategy)(risk.PtrString("ENVIRONMENT_MAX")),
				High:     risk.PtrFloat32(float32(fallbackHigh)),
				Medium:   risk.PtrFloat32(float32(fallbackMedium)),
			},
			Type: risk.ENUMPREDICTORTYPE_VELOCITY,
			Default: &risk.RiskPredictorCommonDefault{
				Weight:    int32(defaultWeight),
				Score:     risk.PtrInt32(int32(defaultScore)),
				Evaluated: risk.PtrBool(false),
			},
		},
	}
}

func defaultUserVelocityByIp() risk.RiskPredictor {

	everyQuantity := 1
	everyMinSample := 5

	useMedium := 0.96
	useHigh := 0.995

	slidingWindowQuantity := 14
	slidingWindowMinSample := 3

	fallbackHigh := 3500.0
	fallbackMedium := 2500.0

	defaultWeight := 5
	defaultScore := 75

	return risk.RiskPredictor{
		RiskPredictorVelocity: &risk.RiskPredictorVelocity{
			Name:        "User Velocity",
			CompactName: "userVelocityByIp",
			Measure:     (*risk.EnumPredictorVelocityMeasure)(risk.PtrString("DISTINCT_COUNT")),
			Of:          risk.PtrString("${event.user.id}"),
			By:          []string{"${event.ip}"},
			Every: &risk.RiskPredictorVelocityAllOfEvery{
				Unit:      (*risk.EnumPredictorUnit)(risk.PtrString("HOUR")),
				Quantity:  risk.PtrInt32(int32(everyQuantity)),
				MinSample: risk.PtrInt32(int32(everyMinSample)),
			},
			Use: &risk.RiskPredictorVelocityAllOfUse{
				Type:   (*risk.EnumPredictorVelocityUseType)(risk.PtrString("POISSON_WITH_MAX")),
				Medium: risk.PtrFloat32(float32(useMedium)),
				High:   risk.PtrFloat32(float32(useHigh)),
			},
			SlidingWindow: &risk.RiskPredictorVelocityAllOfSlidingWindow{
				Unit:      (*risk.EnumPredictorUnit)(risk.PtrString("DAY")),
				Quantity:  risk.PtrInt32(int32(slidingWindowQuantity)),
				MinSample: risk.PtrInt32(int32(slidingWindowMinSample)),
			},
			Fallback: &risk.RiskPredictorVelocityAllOfFallback{
				Strategy: (*risk.EnumPredictorVelocityFallbackStrategy)(risk.PtrString("ENVIRONMENT_MAX")),
				High:     risk.PtrFloat32(float32(fallbackHigh)),
				Medium:   risk.PtrFloat32(float32(fallbackMedium)),
			},
			Type: risk.ENUMPREDICTORTYPE_VELOCITY,
			Default: &risk.RiskPredictorCommonDefault{
				Weight:    int32(defaultWeight),
				Score:     risk.PtrInt32(int32(defaultScore)),
				Evaluated: risk.PtrBool(false),
			},
		},
	}
}

func defaultAnonymousNetwork() risk.RiskPredictor {

	defaultWeight := 5
	defaultScore := 50

	return risk.RiskPredictor{
		RiskPredictorAnonymousNetwork: &risk.RiskPredictorAnonymousNetwork{
			Name:        "Anonymous Network Detection",
			CompactName: "anonymousNetwork",
			WhiteList:   []string{},
			Type:        risk.ENUMPREDICTORTYPE_ANONYMOUS_NETWORK,
			Default: &risk.RiskPredictorCommonDefault{
				Weight:    int32(defaultWeight),
				Score:     risk.PtrInt32(int32(defaultScore)),
				Evaluated: risk.PtrBool(false),
			},
		},
	}
}

func defaultGeoVelocity() risk.RiskPredictor {

	defaultWeight := 5
	defaultScore := 50

	return risk.RiskPredictor{
		RiskPredictorGeovelocity: &risk.RiskPredictorGeovelocity{
			Name:        "Geovelocity Anomaly",
			CompactName: "geoVelocity",
			WhiteList:   []string{},
			Type:        risk.ENUMPREDICTORTYPE_GEO_VELOCITY,
			Default: &risk.RiskPredictorCommonDefault{
				Weight:    int32(defaultWeight),
				Score:     risk.PtrInt32(int32(defaultScore)),
				Evaluated: risk.PtrBool(false),
			},
		},
	}
}

func defaultIpRisk() risk.RiskPredictor {

	defaultWeight := 5
	defaultScore := 50

	return risk.RiskPredictor{
		RiskPredictorIPReputation: &risk.RiskPredictorIPReputation{
			Name:        "IP Reputation",
			CompactName: "ipRisk",
			WhiteList:   []string{},
			Type:        risk.ENUMPREDICTORTYPE_IP_REPUTATION,
			Default: &risk.RiskPredictorCommonDefault{
				Weight:    int32(defaultWeight),
				Score:     risk.PtrInt32(int32(defaultScore)),
				Evaluated: risk.PtrBool(false),
			},
		},
	}
}

func defaultNewDevice() risk.RiskPredictor {

	defaultWeight := 5
	defaultScore := 75

	return risk.RiskPredictor{
		RiskPredictorDevice: &risk.RiskPredictorDevice{
			Name:        "New Device",
			CompactName: "newDevice",
			Detect:      risk.ENUMPREDICTORNEWDEVICEDETECTTYPE_NEW_DEVICE,
			Type:        risk.ENUMPREDICTORTYPE_DEVICE,
			Default: &risk.RiskPredictorCommonDefault{
				Weight:    int32(defaultWeight),
				Score:     risk.PtrInt32(int32(defaultScore)),
				Evaluated: risk.PtrBool(false),
			},
		},
	}
}

func defaultUserLocationAnomaly() risk.RiskPredictor {

	days := 50

	radiusDistance := 50

	defaultWeight := 5
	defaultScore := 50

	return risk.RiskPredictor{
		RiskPredictorUserLocationAnomaly: &risk.RiskPredictorUserLocationAnomaly{
			Name:        "User Location Anomaly",
			CompactName: "userLocationAnomaly",
			Type:        risk.ENUMPREDICTORTYPE_USER_LOCATION_ANOMALY,
			Days:        int32(days),
			Radius: risk.RiskPredictorUserLocationAnomalyAllOfRadius{
				Distance: int32(radiusDistance),
				Unit:     risk.ENUMDISTANCEUNIT_KILOMETERS,
			},
			Default: &risk.RiskPredictorCommonDefault{
				Weight:    int32(defaultWeight),
				Score:     risk.PtrInt32(int32(defaultScore)),
				Evaluated: risk.PtrBool(false),
			},
		},
	}
}
