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

func (odbi *ovndb) ConvertUUIDToString(row libovsdb.Row, field string) string {
	fieldValue := row.Fields[field]
	switch fieldValue.(type) {
	case libovsdb.UUID:
		return fieldValue.(libovsdb.UUID).GoUUID
	}
	return ""
}

func (odbi *ovndb) ConvertGoSetToInt(oset libovsdb.OvsSet) []int {
	var ret []int
	for _, s := range oset.GoSet {
		value, ok := s.(int)
		if ok {
			ret = append(ret, value)
		}
	}
	return ret
}

func (odbi *ovndb) ConvertToStringArray(row libovsdb.Row, field string) []string {
	fieldValue := row.Fields[field]
	switch fieldValue.(type) {
	case libovsdb.UUID:
		return []string{fieldValue.(libovsdb.UUID).GoUUID}
	case string:
		return []string{fieldValue.(string)}
	case libovsdb.OvsSet:
		return odbi.ConvertGoSetToStringArray(fieldValue.(libovsdb.OvsSet))
	}
	return nil
}

func (odbi *ovndb) ConvertToIntArray(row libovsdb.Row, field string) []int {
	fieldValue := row.Fields[field]
	switch fieldValue.(type) {
	case int:
		return []int{fieldValue.(int)}
	case libovsdb.OvsSet:
		return odbi.ConvertGoSetToInt(fieldValue.(libovsdb.OvsSet))
	}
	return nil
}

func (odbi *ovndb) rowToPortBinding(uuid string) (*PortBinding, error) {
	cachePortBinding, ok := odbi.cache[TablePortBinding][uuid]
	if !ok {
		return nil, fmt.Errorf("Port Binding with uuid %s not found", uuid)
	}
	portBinding := &PortBinding{
		UUID:           uuid,
		Chassis:        odbi.ConvertUUIDToString(cachePortBinding, "chassis"),
		Datapath:       odbi.ConvertUUIDToString(cachePortBinding, "datapath"),
		Encap:          odbi.ConvertUUIDToString(cachePortBinding, "encap"),
		ExternalID:     cachePortBinding.Fields["external_ids"].(libovsdb.OvsMap).GoMap,
		GatewayChassis: odbi.ConvertToStringArray(cachePortBinding, "gateway_chassis"),
		HaChassisGroup: odbi.ConvertUUIDToString(cachePortBinding, "ha_chassis_group"),
		LogicalPort:    cachePortBinding.Fields["logical_port"].(string),
		Mac:            odbi.ConvertToStringArray(cachePortBinding, "mac"),
		NatAddresses:   odbi.ConvertToStringArray(cachePortBinding, "nat_addresses"),
		Options:        cachePortBinding.Fields["options"].(libovsdb.OvsMap).GoMap,
		ParentPort:     odbi.ConvertUUIDToString(cachePortBinding, "parent_port"),
		TunnelKey:     cachePortBinding.Fields["tunnel_key"].(int),
		Type:          cachePortBinding.Fields["type"].(string),
		VirtualParent: "",
	}
	tagValue := odbi.ConvertToIntArray(cachePortBinding, "tag")
	portBinding.Tag = 0
	if len(tagValue) != 0 {
		portBinding.Tag = tagValue[0]
	}
	return portBinding, nil
}

func (odbi *ovndb) PortBindingList(searchMap map[string]string) ([]*PortBinding, error) {
	var listPB []*PortBinding

	odbi.cachemutex.RLock()
	defer odbi.cachemutex.RUnlock()

	cachePortBinding, ok := odbi.cache[TablePortBinding]
	if !ok {
		return nil, ErrorSchema
	}

	for uuid := range cachePortBinding {
		var err error
		var matchedKeys int

		portBinding, err := odbi.rowToPortBinding(uuid)
		if err != nil {
			continue
		}
		if searchMap != nil {
			for searchKey, searchValue := range searchMap {
				switch searchKey {
				case "chassis":
					if searchValue == portBinding.Chassis {
						matchedKeys++
					}
				case "logical_port":
					if searchValue == portBinding.LogicalPort {
						matchedKeys++
					}
				}
			}
		}
		if matchedKeys != len(searchMap) {
			continue
		}
		listPB = append(listPB, portBinding)
	}
	return listPB, nil
}
