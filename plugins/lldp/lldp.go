package lldp

import (
	"encoding/json"
	"os/exec"

	"github.com/untangle/discoverd/services/discovery"
	disc "github.com/untangle/golang-shared/services/discovery"
	"github.com/untangle/golang-shared/services/logger"
	"github.com/untangle/golang-shared/structs/protocolbuffers/Discoverd"
)

type jsonData struct {
	Lldp []lldp `json:"lldp"`
}

type lldp struct {
	Intf []intf `json:"interface"`
}

type intf struct {
	Chassis []chassis `json:"chassis"`
	LldpMed []lldpMed `json:"lldp-med"`
}

type chassis struct {
	ID         []identity          `json:"id"`
	Name       []value             `json:"name"`
	Desc       []value             `json:"descr"`
	Capability []chassisCapability `json:"capability"`
}

type identity struct {
	Type  string `json:"type,omitempty"`
	Value string `json:"value,omitempty"`
}

type chassisCapability struct {
	Type    string `json:"type,omitempty"`
	Enabled bool   `json:"enabled,omitempty"`
}

type lldpMed struct {
	DeviceType []value             `json:"device-type"`
	Capability []lldpMedCapability `json:"capability"`
	Inventory  []lldpMedInventory  `json:"inventory"`
}

type lldpMedCapability struct {
	Type      string `json:"type,omitempty"`
	Available bool   `json:"available,omitempty"`
}

type lldpMedInventory struct {
	Hardware     []value `json:"hardware"`
	Software     []value `json:"software"`
	Serial       []value `json:"serial"`
	Manufacturer []value `json:"manufacturer"`
	Model        []value `json:"model"`
}

type value struct {
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
	var result jsonData
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		logger.Err("Unable to unmarshal json: %s\n", err)
	}

	// return on empty data
	if len(result.Lldp[0].Intf) == 0 {
		logger.Debug("No LLDP neighbors found!\n")
		return
	}

	// iterate over interface items
	for _, intf := range result.Lldp[0].Intf {
		// initialize the discovery entry
		entry := disc.DeviceEntry{}
		entry.Init()
		entry.Data.Lldp = &Discoverd.LLDP{}

		// mac is used as id for discovery entry update
		var mac = ""

		if len(intf.Chassis) > 0 {
			chassis := intf.Chassis[0]
			// mac discovery
			if len(chassis.ID) > 0 {
				if chassis.ID[0].Type == "mac" {
					mac = chassis.ID[0].Value
				}
			}

			if len(chassis.Name) > 0 {
				entry.Data.Lldp.SysName = chassis.Name[0].Value
			}
			if len(chassis.Desc) > 0 {
				entry.Data.Lldp.SysDesc = chassis.Desc[0].Value
			}

			// chasis capabilities
			if len(chassis.Capability) > 0 {
				for _, val := range chassis.Capability {
					cap := &Discoverd.LLDPCapabilities{}
					cap.Capability = val.Type
					cap.Enabled = val.Enabled
					entry.Data.Lldp.ChassisCapabilities = append(entry.Data.Lldp.ChassisCapabilities, cap)
				}
			}
		}

		if len(intf.LldpMed) > 0 {
			lldpmed := intf.LldpMed[0]

			// LLDP-MED inventory
			if len(lldpmed.Inventory) > 0 {
				inv := lldpmed.Inventory[0]

				if len(inv.Hardware) > 0 {
					entry.Data.Lldp.InventoryHWRev = inv.Hardware[0].Value
				}
				if len(inv.Software) > 0 {
					entry.Data.Lldp.InventorySoftRev = inv.Software[0].Value
				}
				if len(inv.Serial) > 0 {
					entry.Data.Lldp.InventorySerial = inv.Serial[0].Value
				}
				if len(inv.Model) > 0 {
					entry.Data.Lldp.InventoryModel = inv.Model[0].Value
				}
				if len(inv.Manufacturer) > 0 {
					entry.Data.Lldp.InventoryVendor = inv.Manufacturer[0].Value
				}
			}

			// LLDP-MED capabilities
			if len(lldpmed.Capability) > 0 {
				for _, val := range lldpmed.Capability {
					cap := &Discoverd.LLDPCapabilities{}
					cap.Capability = val.Type
					cap.Enabled = val.Available
					entry.Data.Lldp.MedCapabilities = append(entry.Data.Lldp.MedCapabilities, cap)
				}
			}
		}

		if mac != "" {
			entry.Data.MacAddress = mac
			discovery.UpdateDiscoveryEntry(mac, entry)
		}
	}
}
