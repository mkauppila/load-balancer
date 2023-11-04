package config

import (
	"testing"
)

func TestParseBasicConfiguration(t *testing.T) {
	contents := []byte(`
		health_check on 5000 /path/path
		server http://localhost:4001
	`)
	conf, err := ParseConfiguration(contents)
	if conf.HealthCheck.Enabled != true {
		t.Error("Health check should be enabled")
	}
	if conf.HealthCheck.IntervalMs != 5000 {
		t.Error("Health check duration should be 5000")
	}
	if conf.HealthCheck.Path != "/path/path" {
		t.Error("Health check path should be /path/path")
	}
	if err != nil {
		t.Error("Unexpected error")
	}
}

func TestParsingHealthCheckEnabled(t *testing.T) {
	hc, _ := parseHealthCheck("health_check on")
	if hc.Enabled != true {
		t.Error("Health check should be enabled")
	}
	hc, _ = parseHealthCheck("health_check off")
	if hc.Enabled != false {
		t.Error("Health check should be disabled")
	}

	_, err := parseHealthCheck("health_check faulty")
	if err == nil || err.Error() != "Unknown option: faulty" {
		t.Error("Health check should error")
	}
}

func TestParsingHealthCheckIntervalErrorCase(t *testing.T) {
	_, err := parseHealthCheck("health_check on 50a1")
	if err == nil || err.Error() != "interval needs to be an integer" {
		t.Error("Invalid interval")
	}
}

func TestParsingHealthCheckPath(t *testing.T) {
	hc, _ := parseHealthCheck("health_check on 501 /pathie")
	if hc.Path != "/pathie" {
		t.Error("Path not parsed properly. Actual ", hc.Path)
	}
}

func TestParsePort(t *testing.T) {
	port, _ := parsePort("listen 1447")
	if port != 1447 {
		t.Error("Port not parsed properly. Actual ", port)
	}
}
