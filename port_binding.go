package goovn

import (
	"fmt"
	"github.com/ebay/libovsdb"
)

// Port_Binding table OVN SB
type PortBinding struct {
	UUID           string
	Chassis        string
	Datapath       string
	Encap          string
	ExternalID     map[interface{}]interface{}
	GatewayChassis []string
	HaChassisGroup string
	LogicalPort    string
	Mac            []string
	NatAddresses   []string
	Options        map[interface{}]interface{}
	ParentPort     string
	Tag            int
	TunnelKey      int
	Type           string
	VirtualParent  string
}

func (odbi *ovndb) rowToPortBinding(uuid string) (*PortBinding, error) {
	cachePortBinding, ok := odbi.cache[TablePortBinding][uuid]
	if !ok {
		return nil, fmt.Errorf("Port Binding with uuid %s not found", uuid)
	}
	portBinding := &PortBinding{
		UUID:           uuid,
		Chassis:        cachePortBinding.Fields["chassis"].(string),
		Datapath:       cachePortBinding.Fields["datapath"].(string),
		Encap:          cachePortBinding.Fields["encap"].(string),
		ExternalID:     cachePortBinding.Fields["external_ids"].(libovsdb.OvsMap).GoMap,
		GatewayChassis: nil,
		HaChassisGroup: cachePortBinding.Fields["ha_chassis_group"].(string),
		LogicalPort:    cachePortBinding.Fields["logical_port"].(string),
		Mac:            nil,
		NatAddresses:   nil,
		Options:        cachePortBinding.Fields["options"].(libovsdb.OvsMap).GoMap,
		ParentPort:     cachePortBinding.Fields["parent_port"].(string),
		Tag:            cachePortBinding.Fields["tag"].(int),
		TunnelKey:      cachePortBinding.Fields["tunnel_key"].(int),
		Type:           cachePortBinding.Fields["type"].(string),
		VirtualParent:  cachePortBinding.Fields["virtual_parent"].(string),
	}
	return portBinding, nil
}

func (ovnSB *ovndb) portBindingListImp(searchMap map[string]string) ([]*PortBinding, error) {
	var listPB []*PortBinding

	ovnSB.cachemutex.RLock()
	defer ovnSB.cachemutex.RUnlock()

	cachePortBinding, ok := ovnSB.cache[TablePortBinding]
	if !ok {
		return nil, ErrorSchema
	}

	for uuid, drows := range cachePortBinding {
		var pb *PortBinding
		var err error
		var matchedKeys int

		if searchMap != nil {
			for searchKey, searchValue := range searchMap {
				if colValue, ok := drows.Fields[searchKey]; ok {
					if colValue, ok := colValue.(string); ok && colValue == searchValue {
						matchedKeys++
					}
				}
			}
		}
		if matchedKeys != len(searchMap) {
			continue
		}
		pb, err = ovnSB.rowToPortBinding(uuid)
		if err != nil {
			return nil, err
		}
		listPB = append(listPB, pb)
	}
	return listPB, nil
}
