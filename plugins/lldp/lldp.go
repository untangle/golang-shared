package lldp

import (
	"encoding/json"
	"os/exec"

	"github.com/untangle/discoverd/services/discovery"
	"github.com/untangle/golang-shared/services/logger"
	"github.com/untangle/golang-shared/structs/protocolbuffers/Discoverd"
)

// Start starts the LLDP collector
func Start() {
	logger.Info("Starting LLDP collector plugin\n")
	discovery.RegisterCollector(LldpcallBackHandler)
}

// Stop stops LLDP collector
func Stop() {
}

// LldpcallBackHandler is the callback handler for the LLDP collector
func LldpcallBackHandler(commands []discovery.Command) {
	logger.Debug("LLDP neighbors callback handler: Received %d commands\n", len(commands))

	// run neighbors command
	cmd := exec.Command("lldpcli", "-f", "json", "show", "n", "-d")
	output, _ := cmd.CombinedOutput()

	logger.Debug(string(output))

	// parse json output data
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(output), &result); err != nil {
			logger.Err("Unable to unmarshal json: %s\n", err)
	}

	// extract lldp -> interface slice
	interfaces, _ := result["lldp"].(map[string]interface{})["interface"]

	// iterate over interface items
	for _, intf := range interfaces.([]interface{}) {
		// mac is used as id for discovery entry update
		var mac = ""

		// LLDP proto
		var sysName string = ""
		var sysDesc string = ""
		var inventoryHWRev string = ""
		var inventorySoftRev string = ""
		var inventorySerial string = ""
		var inventoryAssetTag string = ""
		var inventoryModel string = ""
		var inventoryVendor string = ""

		for _, intf_map := range intf.(map[string]interface{}) {

			// chassis
			chassis, _ := intf_map.(map[string]interface{})["chassis"]

			if (chassis != nil) {
				for chassis_key, chassis_map := range chassis.(map[string]interface{}) {
					// chassis_key is the sysName
					sysName = chassis_key

					ch := chassis_map.(map[string]interface{})
					sysDesc, _ = ch["descr"].(string)

					chassis_id, _ := ch["id"].(map[string]interface{})
					chassis_id_type, _ := chassis_id["type"].(string)
					chassis_id_value, _ := chassis_id["value"].(string)

					// if chassis id type is "mac", populate mac value
					if (chassis_id_type == "mac") {
							mac = chassis_id_value
					}
				}
			}

			// LLPD-MED
			lldp_med, _ := intf_map.(map[string]interface{})["lldp-med"]
			if (lldp_med != nil) {
				inventory, _ := lldp_med.(map[string]interface {})["inventory"]

				if (inventory != nil) {
					for inv_key, inv_map := range inventory.(map[string]interface{}) {
						switch inv_key {
							case "hardware":
								inventoryHWRev = inv_map.(string)
							case "software":
								inventorySoftRev = inv_map.(string)
							case "serial":
								inventorySerial = inv_map.(string)
							case "model":
								inventoryModel = inv_map.(string)
							case "manufacturer":
								inventoryVendor = inv_map.(string)
						}
					}
				}
			}

			entry := discovery.DeviceEntry{}
			entry.Init()
			entry.Data.Lldp = &Discoverd.LLDP{}

			entry.Data.Lldp.SysName = sysName
			entry.Data.Lldp.SysDesc = sysDesc

			entry.Data.Lldp.InventoryHWRev = inventoryHWRev
			entry.Data.Lldp.InventorySoftRev = inventorySoftRev
			entry.Data.Lldp.InventorySerial = inventorySerial
			entry.Data.Lldp.InventoryAssetTag = inventoryAssetTag
			entry.Data.Lldp.InventoryModel = inventoryModel
			entry.Data.Lldp.InventoryVendor = inventoryVendor

			if (mac != "") {
				discovery.UpdateDiscoveryEntry(mac, entry)
			}
		}
	}
}
