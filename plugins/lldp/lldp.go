package lldp

import (
	"encoding/json"
	"os/exec"

	"github.com/untangle/discoverd/services/discovery"
	"github.com/untangle/golang-shared/services/logger"
	"github.com/untangle/golang-shared/structs/protocolbuffers/Discoverd"
)

type Json struct {
	Lldp []Lldp `json:"lldp"`
}

type Lldp struct {
	Intf []Intf `json:"interface"`
}

type Intf struct {
	Chassis []Chassis `json:"chassis"`
	LldpMed []LldpMed `json:"lldp-med"`
}

type Chassis struct {
	Id []ChassisId                 `json:"id"`
	Name []Value                   `json:"name"`
	Desc []Value                   `json:"descr"`
	Capability []ChassisCapability `json:"capability"`
}

type ChassisId struct {
	Type  string `json:"type,omitempty"`
	Value string `json:"value,omitempty"`
}

type ChassisCapability struct {
	Type string  `json:"type,omitempty"`
	Enabled bool `json:"enabled,omitempty"`
}

type LldpMed struct {
	DeviceType []Value             `json:"device-type"`
	Capability []LldpMedCapability `json:"capability"`
	Inventory []LldpMedInventory   `json:"inventory"`
}

type LldpMedCapability struct {
	Type string    `json:"type,omitempty"`
	Available bool `json:"available,omitempty"`
}

type LldpMedInventory struct {
	Hardware []Value     `json:"hardware"`
	Software []Value     `json:"software"`
	Serial []Value       `json:"serial"`
	Manufacturer []Value `json:"manufacturer"`
	Model []Value        `json:"model"`
}

type Value struct {
	Value string `json:"value,omitempty"`
}

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
	cmd := exec.Command("lldpcli", "-f", "json0", "show", "n", "details")
	output, _ := cmd.CombinedOutput()

	// logger.Info("LLDP output: %s\n", string(output))

	// parse json output data
	var result Json
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		logger.Err("Unable to unmarshal json: %s\n", err)
	}

	// return on empty data
	if (len(result.Lldp[0].Intf) == 0) {
		logger.Debug("No LLDP neighbors found!\n")
		return
	}

	// iterate over interface items
	for _, intf := range result.Lldp[0].Intf {
		// initialize the discovery entry
		entry := discovery.DeviceEntry{}
		entry.Init()
		entry.Data.Lldp = &Discoverd.LLDP{}

		// mac is used as id for discovery entry update
		var mac = ""

		if (len(intf.Chassis) > 0) {
			chassis := intf.Chassis[0]
			// mac discovery
			if (len(chassis.Id) > 0) {
				if (chassis.Id[0].Type == "mac") {
					mac = chassis.Id[0].Value
				}
			}

			if (len(chassis.Name) > 0) {
				entry.Data.Lldp.SysName = chassis.Name[0].Value
			}
			if (len(chassis.Desc) > 0) {
				entry.Data.Lldp.SysDesc = chassis.Desc[0].Value
			}

			// chasis capabilities
			if (len(chassis.Capability) > 0) {
				for _, val := range chassis.Capability {
					cap := &Discoverd.LLDPCapabilities{}
					cap.Capability = val.Type
					cap.Enabled = val.Enabled
					entry.Data.Lldp.ChassisCapabilities = append(entry.Data.Lldp.ChassisCapabilities, cap)
				}
			}
		}

		if (len(intf.LldpMed) > 0) {
			lldpmed := intf.LldpMed[0]

			// LLDP-MED inventory
			if (len(lldpmed.Inventory) > 0) {
				inv := lldpmed.Inventory[0]

				if (len(inv.Hardware) > 0) {
					entry.Data.Lldp.InventoryHWRev = inv.Hardware[0].Value
				}
				if (len(inv.Software) > 0) {
					entry.Data.Lldp.InventorySoftRev = inv.Software[0].Value
				}
				if (len(inv.Serial) > 0) {
					entry.Data.Lldp.InventorySerial = inv.Serial[0].Value
				}
				if (len(inv.Model) > 0) {
					entry.Data.Lldp.InventoryModel = inv.Model[0].Value
				}
				if (len(inv.Manufacturer) > 0) {
					entry.Data.Lldp.InventoryVendor = inv.Manufacturer[0].Value
				}
			}

			// LLDP-MED capabilities
			if (len(lldpmed.Capability) > 0) {
				for _, val := range lldpmed.Capability {
					cap := &Discoverd.LLDPCapabilities{}
					cap.Capability = val.Type
					cap.Enabled = val.Available
					entry.Data.Lldp.MedCapabilities = append(entry.Data.Lldp.MedCapabilities, cap)
				}
			}
		}

		if (mac != "") {
			discovery.UpdateDiscoveryEntry(mac, entry)
		}
	}
}
