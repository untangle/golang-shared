package lldp

import (
	"encoding/json"
	"os/exec"
	"reflect"

	"github.com/untangle/discoverd/services/discovery"
	"github.com/untangle/golang-shared/services/logger"
	"github.com/untangle/golang-shared/structs/protocolbuffers/Discoverd"
)

// Start starts the LLDP collector
func Start() {
	logger.Info("Starting LLDP collector plugin\n")
	discovery.RegisterCollector(LldpcallBackHandler)

	// initial run
	LldpcallBackHandler(nil)
}

// Stop stops LLDP collector
func Stop() {
}

// LldpcallBackHandler is the callback handler for the LLDP collector
func LldpcallBackHandler(commands []discovery.Command) {
	logger.Debug("LLDP neighbors callback handler: Received %d commands\n", len(commands))

	// run neighbors command
	cmd := exec.Command("lldpcli", "-f", "json", "show", "n", "details")
	output, _ := cmd.CombinedOutput()

	// logger.Info("LLDP output: %s\n", string(output))

	// parse json output data
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		logger.Err("Unable to unmarshal json: %s\n", err)
	}

	lldp, _ := result["lldp"].(map[string]interface{})

	if (lldp["interface"] == nil) {
		logger.Debug("No LLDP neighbors found!\n")
		return
	}

	interfacesType := reflect.TypeOf(lldp["interface"])

	// create a slice (array) of interfaces
	interfaces := make([]interface{}, 0, 0)
	// when having an array of interfaces
	if (interfacesType.Kind() == reflect.Slice) {
		interfaces = lldp["interface"].([]interface{})
	}
	// when just a single interface
	if (interfacesType.Kind() == reflect.Map) {
		for _, value := range lldp["interface"].(map[string]interface {}) {
			interfaces = append(interfaces, value)
		}
	}

	// iterate over interface items
	for _, intf := range interfaces {
		// initialize the discovery entry
		entry := discovery.DeviceEntry{}
		entry.Init()
		entry.Data.Lldp = &Discoverd.LLDP{}

		// mac is used as id for discovery entry update
		var mac = ""

		// processing interfaces "eth0": { ... }
		for _, value := range intf.(map[string]interface {}) {
			// chassis
			chassis, _ := value.(map[string]interface {})["chassis"]
			if (chassis == nil) {
				continue
			}
			for key, value := range chassis.(map[string]interface {}) {
				entry.Data.Lldp.SysName = key
				ch, _ := value.(map[string]interface {})
				entry.Data.Lldp.SysDesc = ch["descr"].(string)

				chassis_id, _ := ch["id"].(map[string]interface{})
				chassis_id_type, _ := chassis_id["type"].(string)
				chassis_id_value, _ := chassis_id["value"].(string)

				// if chassis id type is "mac", populate mac value
				if (chassis_id_type == "mac") {
					mac = chassis_id_value
				}
			}

			// LLDP-MED
			lldp_med, _ := value.(map[string]interface{})["lldp-med"]
			if (lldp_med == nil) {
				continue
			}
			inventory, _ := lldp_med.(map[string]interface {})["inventory"]
			if (inventory == nil) {
				continue
			}
			for inv_key, inv_map := range inventory.(map[string]interface{}) {
				switch inv_key {
					case "hardware":
						entry.Data.Lldp.InventoryHWRev = inv_map.(string)
					case "software":
						entry.Data.Lldp.InventorySoftRev = inv_map.(string)
					case "serial":
						entry.Data.Lldp.InventorySerial = inv_map.(string)
					case "model":
						entry.Data.Lldp.InventoryModel = inv_map.(string)
					case "manufacturer":
						entry.Data.Lldp.InventoryVendor = inv_map.(string)
				}
			}
		}

		if (mac != "") {
			discovery.UpdateDiscoveryEntry(mac, entry)
		}
	}
}
