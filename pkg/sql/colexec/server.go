// Copyright 2021 Matrix Origin
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package colexec

import (
	"github.com/google/uuid"
	"github.com/matrixorigin/matrixone/pkg/logservice"
	"github.com/matrixorigin/matrixone/pkg/objectio"
	"github.com/matrixorigin/matrixone/pkg/vm/process"
)

var Srv *Server

const (
	TxnWorkSpaceIdType = 1
	CnBlockIdType      = 2
)

func NewServer(client logservice.CNHAKeeperClient) *Server {
	if Srv != nil {
		return Srv
	}
	Srv = &Server{
		mp:            make(map[uint64]*process.WaitRegister),
		hakeeper:      client,
		uuidCsChanMap: UuidProcMap{mp: make(map[uuid.UUID]uuidProcMapItem, 1024)},
		cnSegmentMap:  CnSegmentMap{mp: make(map[objectio.Segmentid]int32, 1024)},
	}
	return Srv
}

func (srv *Server) GetConnector(id uint64) *process.WaitRegister {
	srv.Lock()
	defer srv.Unlock()
	defer func() { delete(srv.mp, id) }()
	return srv.mp[id]
}

func (srv *Server) RegistConnector(reg *process.WaitRegister) uint64 {
	srv.Lock()
	defer srv.Unlock()
	srv.mp[srv.id] = reg
	defer func() { srv.id++ }()
	return srv.id
}

// GetProcByUuid try to get process from map, and decrease referenceCount by 1 when get succeed.
// if forcedDelete is true, it will do action to avoid other goroutine to put a new item into map when get failed.
func (srv *Server) GetProcByUuid(u uuid.UUID, forcedDelete bool) (*process.Process, bool) {
	srv.uuidCsChanMap.Lock()
	defer srv.uuidCsChanMap.Unlock()
	p, ok := srv.uuidCsChanMap.mp[u]
	if !ok {
		if forcedDelete {
			srv.uuidCsChanMap.mp[u] = uuidProcMapItem{proc: nil, referenceCount: 1}
		}
		return nil, false
	}
	p.referenceCount--
	if p.referenceCount == 0 {
		delete(srv.uuidCsChanMap.mp, u)
		return nil, true
	} else {
		srv.uuidCsChanMap.mp[u] = p
	}
	return p.proc, true
}

func (srv *Server) PutProcIntoUuidMap(u uuid.UUID, p *process.Process) error {
	srv.uuidCsChanMap.Lock()
	defer srv.uuidCsChanMap.Unlock()
	if _, ok := srv.uuidCsChanMap.mp[u]; ok {
		delete(srv.uuidCsChanMap.mp, u)
		return nil
	}

	srv.uuidCsChanMap.mp[u] = uuidProcMapItem{proc: p, referenceCount: 2}
	return nil
}

func (srv *Server) DeleteUuids(uuids []uuid.UUID) {
	srv.uuidCsChanMap.Lock()
	defer srv.uuidCsChanMap.Unlock()
	for i := range uuids {
		uid, ok := srv.uuidCsChanMap.mp[uuids[i]]
		if !ok {
			continue
		}

		uid.referenceCount--
		if uid.referenceCount == 0 {
			delete(srv.uuidCsChanMap.mp, uuids[i])
		} else {
			srv.uuidCsChanMap.mp[uuids[i]] = uid
		}
	}
}

func (srv *Server) PutCnSegment(sid *objectio.Segmentid, segmentType int32) {
	srv.cnSegmentMap.Lock()
	defer srv.cnSegmentMap.Unlock()
	srv.cnSegmentMap.mp[*sid] = segmentType
}

func (srv *Server) DeleteTxnSegmentIds(sids []objectio.Segmentid) {
	srv.cnSegmentMap.Lock()
	defer srv.cnSegmentMap.Unlock()
	for _, segmentName := range sids {
		delete(srv.cnSegmentMap.mp, segmentName)
	}
}

func (srv *Server) GetCnSegmentMap() map[string]int32 {
	srv.cnSegmentMap.Lock()
	defer srv.cnSegmentMap.Unlock()
	new_mp := make(map[string]int32)
	for k, v := range srv.cnSegmentMap.mp {
		new_mp[string(k[:])] = v
	}
	return new_mp
}

func (srv *Server) GetCnSegmentType(sid *objectio.Segmentid) int32 {
	srv.cnSegmentMap.Lock()
	defer srv.cnSegmentMap.Unlock()
	return srv.cnSegmentMap.mp[*sid]
}

// SegmentId is part of Id for cn2s3 directly, for more info, refer to docs about it
func (srv *Server) GenerateObject() objectio.ObjectName {
	srv.Lock()
	defer srv.Unlock()
	return objectio.BuildObjectName(objectio.NewSegmentid(), 0)
	// for future fileOffset
	// if srv.InitSegmentId {
	// 	srv.incrementSegmentId()
	// } else {
	// 	srv.getNewSegmentId()
	// 	srv.currentFileOffset = 0
	// 	srv.InitSegmentId = true
	// }
	// return objectio.BuildObjectName(srv.CNSegmentId, srv.currentFileOffset)
}

// func (srv *Server) incrementSegmentId() {
// 	if srv.currentFileOffset < math.MaxUint16 {
// 		srv.currentFileOffset++
// 	} else {
// 		srv.getNewSegmentId()
// 		srv.currentFileOffset = 0
// 	}
// }

// // for now, rowId is common between CN and DN.
// func (srv *Server) getNewSegmentId() {
// 	srv.CNSegmentId = common.NewSegmentid()
// }
