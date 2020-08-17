package goovn

import (
	"fmt"
	"github.com/ebay/libovsdb"
)

// Port_Binding table OVN SB
type DatapathBinding struct {
	UUID       string
	ExternalID map[interface{}]interface{}
	TunnelKey  int
}

func (odbi *ovndb) rowToDatapathBinding(uuid string) (*DatapathBinding, error) {
	cacheDatapathBinding, ok := odbi.cache[TableDatapathBinding][uuid]
	if !ok {
		return nil, fmt.Errorf("Datapath Binding with uuid %s not found", uuid)
	}
	datapathBinding := &DatapathBinding{
		UUID:       uuid,
		ExternalID: cacheDatapathBinding.Fields["external_ids"].(libovsdb.OvsMap).GoMap,
		TunnelKey: cacheDatapathBinding.Fields["tunnel_key"].(int),
	}
	return datapathBinding, nil
}

func (odbi *ovndb) DatapathBindingGet(uuid string) (*DatapathBinding, error) {
	odbi.cachemutex.RLock()
	defer odbi.cachemutex.RUnlock()

	_, ok := odbi.cache[TableDatapathBinding]
	if !ok {
		return nil, ErrorSchema
	}

	return odbi.rowToDatapathBinding(uuid)
}

func (odbi *ovndb) DatapathBindingGetByName(name string) ([]*DatapathBinding, error) {
	var listDatapathBinding []*DatapathBinding

	odbi.cachemutex.RLock()
	defer odbi.cachemutex.RUnlock()

	cacheDatapathBinding, ok := odbi.cache[TableDatapathBinding]
	if !ok {
		return nil, ErrorSchema
	}

	for uuid, drows := range cacheDatapathBinding {
		if external_ids, ok := drows.Fields["external_ids"].(libovsdb.OvsMap); ok {
			if external_ids.GoMap["name"] == name {
				dpBinding, err := odbi.rowToDatapathBinding(uuid)
				if err != nil {
					return nil, err
				}
				listDatapathBinding = append(listDatapathBinding, dpBinding)
			}

		}
	}
	return listDatapathBinding, nil
}

