/**
 * Copyright (c) 2020 eBay Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 **/

package goovn

import "github.com/ebay/libovsdb"

type NBGlobalTableRow struct {
	UUID        string
	Options     map[interface{}]interface{}
	ExternalID  map[interface{}]interface{}
	Connections []string
	SSL         string
	IPSec       bool
	NBCfg       int
	SBCfg       int
	HVCfg       int
}

func (odbi *ovndb) nbGlobalAddImp(options map[string]string) (*OvnCommand, error) {
	return odbi.addGlobalTableRowImp(options, TableNBGlobal)
}

func (odbi *ovndb) nbGlobalDelImp() (*OvnCommand, error) {
	return odbi.delGlobalTableRowImp(TableNBGlobal)
}

// ovsdb-client -v transact '["Open_vSwitch", {"op" : "update", "table" : "NB_Global", "where": [["_uuid", "==", ["uuid", "587c6ee2-93f9-4bd8-9794-f4a983d139a4"]]],
// "row":{ "options" : [ "map", [[ "bar", "baz"],["engine_test", "engine-foo"]]],}}]'

func (odbi *ovndb) nbGlobalSetOptionsImp(options map[string]string) (*OvnCommand, error) {
	return odbi.globalSetOptionsImp(options, TableNBGlobal)
}

func (odbi *ovndb) nbGlobalGetOptionsImp() (map[string]string, error) {
	return odbi.globalGetOptionsImp(TableNBGlobal)
}

func (odbi *ovndb) rowToNBGlobal(uuid string) *NBGlobalTableRow {
	cacheNBGlobal, ok := odbi.cache[TableNBGlobal][uuid]
	if !ok {
		return nil
	}

	nbGlobal := &NBGlobalTableRow{
		UUID:       uuid,
		Options:    cacheNBGlobal.Fields["options"].(libovsdb.OvsMap).GoMap,
		ExternalID: cacheNBGlobal.Fields["external_ids"].(libovsdb.OvsMap).GoMap,
		IPSec:      cacheNBGlobal.Fields["ipsec"].(bool),
		NBCfg:      cacheNBGlobal.Fields["nb_cfg"].(int),
		SBCfg:      cacheNBGlobal.Fields["sb_cfg"].(int),
		HVCfg:      cacheNBGlobal.Fields["hv_cfg"].(int),
	}
	switch cacheNBGlobal.Fields["ssl"].(type) {
	case libovsdb.UUID:
		nbGlobal.SSL = cacheNBGlobal.Fields["ssl"].(libovsdb.UUID).GoUUID
	default:
	}
	connections := cacheNBGlobal.Fields["connections"]
	switch connections.(type) {
	case string:
		nbGlobal.Connections = []string{connections.(string)}
	case libovsdb.OvsSet:
		nbGlobal.Connections = odbi.ConvertGoSetToStringArray(connections.(libovsdb.OvsSet))
	}
	return nbGlobal
}

func (odbi *ovndb) NBGlobalList() ([]*NBGlobalTableRow, error) {
	odbi.cachemutex.RLock()
	defer odbi.cachemutex.RUnlock()

	cacheNBGlobal, ok := odbi.cache[TableNBGlobal]
	if !ok {
		return nil, ErrorNotFound
	}

	nbGlobalList := make([]*NBGlobalTableRow, 0, len(cacheNBGlobal))
	i := 0
	for uuid := range cacheNBGlobal {
		nbGlobalList[i] = odbi.rowToNBGlobal(uuid)
		i++
	}

	if len(nbGlobalList) == 0 {
		return nil, ErrorNotFound
	}
	return nbGlobalList, nil
}
