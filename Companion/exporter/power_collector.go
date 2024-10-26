package exporter

import (
	"log"
	"strconv"
)

type PowerInfo struct {
	CircuitGroupId float64 `json:"CircuitGroupID"`
	PowerConsumed  float64 `json:"PowerConsumed"`
}

type PowerCollector struct {
	endpoint string
}

type PowerDetails struct {
	CircuitGroupId      float64 `json:"CircuitGroupID"`
	PowerConsumed       float64 `json:"PowerConsumed"`
	PowerCapacity       float64 `json:"PowerCapacity"`
	PowerMaxConsumed    float64 `json:"PowerMaxConsumed"`
	BatteryDifferential float64 `json:"BatteryDifferential"`
	BatteryPercent      float64 `json:"BatteryPercent"`
	BatteryCapacity     float64 `json:"BatteryCapacity"`
	BatteryTimeEmpty    string  `json:"BatteryTimeEmpty"`
	BatteryTimeFull     string  `json:"BatteryTimeFull"`
	FuseTriggered       bool    `json:"FuseTriggered"`
}

func NewPowerCollector(endpoint string) *PowerCollector {
	return &PowerCollector{
		endpoint: endpoint,
	}
}

func (c *PowerCollector) Collect(frmAddress string, sessionName string) {
	details := []PowerDetails{}
	err := retrieveData(frmAddress+c.endpoint, &details)
	if err != nil {
		log.Printf("error reading power statistics from FRM: %s\n", err)
		return
	}

	for _, d := range details {
		circuitId := strconv.FormatFloat(d.CircuitGroupId, 'f', -1, 64)
		PowerConsumed.WithLabelValues(circuitId, frmAddress, sessionName).Set(d.PowerConsumed)
		PowerCapacity.WithLabelValues(circuitId, frmAddress, sessionName).Set(d.PowerCapacity)
		PowerMaxConsumed.WithLabelValues(circuitId, frmAddress, sessionName).Set(d.PowerMaxConsumed)
		BatteryDifferential.WithLabelValues(circuitId, frmAddress, sessionName).Set(d.BatteryDifferential)
		BatteryPercent.WithLabelValues(circuitId, frmAddress, sessionName).Set(d.BatteryPercent)
		BatteryCapacity.WithLabelValues(circuitId, frmAddress, sessionName).Set(d.BatteryCapacity)
		batterySecondsRemaining := parseTimeSeconds(d.BatteryTimeEmpty)
		if batterySecondsRemaining != nil {
			BatterySecondsEmpty.WithLabelValues(circuitId, frmAddress, sessionName).Set(*batterySecondsRemaining)
		}
		batterySecondsFull := parseTimeSeconds(d.BatteryTimeFull)
		if batterySecondsFull != nil {
			BatterySecondsFull.WithLabelValues(circuitId, frmAddress, sessionName).Set(*batterySecondsFull)
		}
		fuseTriggered := parseBool(d.FuseTriggered)
		FuseTriggered.WithLabelValues(circuitId, frmAddress, sessionName).Set(fuseTriggered)
	}
}
