package service

import (
	"context"
	"log"
	"net"
	"strings"
)

// FraudDetector provides fraud detection capabilities
type FraudDetector struct {
	logger log.Logger
}

// NewFraudDetector creates a new fraud detector
func NewFraudDetector() *FraudDetector {
	return &FraudDetector{}
}

// AssessCheckoutRisk assesses the risk level of a checkout transaction
func (fd *FraudDetector) AssessCheckoutRisk(ctx context.Context, assessment *CheckoutRiskAssessment) float64 {
	riskScore := 0.0

	// IP-based risk assessment
	ipRisk := fd.assessIPRisk(assessment.ClientIP)
	riskScore += ipRisk * 0.3

	// Amount-based risk assessment
	amountRisk := fd.assessAmountRisk(assessment.Amount, assessment.Currency)
	riskScore += amountRisk * 0.4

	// User behavior risk assessment (placeholder)
	behaviorRisk := fd.assessUserBehaviorRisk(ctx, assessment.UserID)
	riskScore += behaviorRisk * 0.3

	// Ensure score is between 0 and 1
	if riskScore > 1.0 {
		riskScore = 1.0
	}

	return riskScore
}

// assessIPRisk evaluates risk based on IP address
func (fd *FraudDetector) assessIPRisk(clientIP string) float64 {
	if clientIP == "" {
		return 0.5 // Medium risk for missing IP
	}

	// Parse IP address
	ip := net.ParseIP(clientIP)
	if ip == nil {
		return 0.6 // Higher risk for invalid IP
	}

	// Check for localhost/private IPs (could indicate testing/development)
	if ip.IsLoopback() || ip.IsPrivate() {
		return 0.1 // Low risk for internal IPs
	}

	// Check for known VPN/proxy patterns (simplified)
	if fd.isLikelyVPN(clientIP) {
		return 0.7 // Higher risk for VPN/proxy
	}

	return 0.2 // Base risk for regular IPs
}

// assessAmountRisk evaluates risk based on transaction amount
func (fd *FraudDetector) assessAmountRisk(amount float64, currency string) float64 {
	// Convert to base risk scale
	var baseAmount float64

	switch strings.ToUpper(currency) {
	case "VND":
		// For VND, amounts are typically larger due to smaller currency unit
		if amount > 10000000 { // > 10 million VND (~$400 USD)
			return 0.8
		} else if amount > 5000000 { // > 5 million VND (~$200 USD)
			return 0.5
		} else if amount > 1000000 { // > 1 million VND (~$40 USD)
			return 0.2
		}
		return 0.1
	case "USD":
		baseAmount = amount
	default:
		baseAmount = amount // Assume USD equivalent for unknown currencies
	}

	// Risk assessment for USD-equivalent amounts
	if baseAmount > 500 {
		return 0.8 // High risk for very expensive purchases
	} else if baseAmount > 200 {
		return 0.5 // Medium risk
	} else if baseAmount > 50 {
		return 0.2 // Low-medium risk
	}

	return 0.1 // Low risk for small amounts
}

// assessUserBehaviorRisk evaluates risk based on user behavior patterns
func (fd *FraudDetector) assessUserBehaviorRisk(ctx context.Context, userID string) float64 {
	// Placeholder implementation
	// In a real system, this would check:
	// - Recent purchase frequency
	// - Account age
	// - Previous failed transactions
	// - Geographic patterns
	// - Device fingerprinting

	if userID == "" {
		return 0.9 // Very high risk for anonymous users
	}

	// For now, return a low base risk for authenticated users
	return 0.1
}

// isLikelyVPN checks if an IP is likely from a VPN or proxy
func (fd *FraudDetector) isLikelyVPN(ip string) bool {
	// Simplified VPN detection
	// In production, you would use a proper IP intelligence service
	// like MaxMind GeoIP2, IPQualityScore, etc.

	// Some basic patterns that might indicate VPN/proxy
	vpnPatterns := []string{
		"tor-exit",
		"proxy",
		"vpn",
		"anonymizer",
	}

	lowerIP := strings.ToLower(ip)
	for _, pattern := range vpnPatterns {
		if strings.Contains(lowerIP, pattern) {
			return true
		}
	}

	return false
}

// GetRiskLevel converts numeric risk score to descriptive level
func (fd *FraudDetector) GetRiskLevel(score float64) string {
	switch {
	case score >= 0.8:
		return "HIGH"
	case score >= 0.5:
		return "MEDIUM"
	case score >= 0.3:
		return "LOW"
	default:
		return "MINIMAL"
	}
}

// ShouldBlockTransaction determines if a transaction should be blocked
func (fd *FraudDetector) ShouldBlockTransaction(score float64) bool {
	return score >= 0.8 // Block transactions with high risk
}

// LogRiskAssessment logs the risk assessment for audit purposes
func (fd *FraudDetector) LogRiskAssessment(ctx context.Context, assessment *CheckoutRiskAssessment, score float64) {
	level := fd.GetRiskLevel(score)
	// In production, you would log this to your fraud detection system
	log.Printf("Fraud Risk Assessment: UserID=%s, IP=%s, Amount=%f %s, Risk=%.2f (%s)",
		assessment.UserID, assessment.ClientIP, assessment.Amount, assessment.Currency, score, level)
}