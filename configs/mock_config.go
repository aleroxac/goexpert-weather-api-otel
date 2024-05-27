package configs

func MockConfig() *Conf {
	return &Conf{
		InputApiHttpPort:                      "8080",
		InputApiOtelServiceName:               "test",
		OrchestratorApiPort:                   "8081",
		OrchestratorApiHost:                   "test",
		OpenWeathermapApiKey:                  "test",
		OrchestratorApiServiceName:            "test",
		OpenTelemetryCollectorExporerEndpoint: "test",
	}
}
