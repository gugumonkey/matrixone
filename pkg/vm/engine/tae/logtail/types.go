// Copyright 2021 Matrix Origin
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package logtail

import (
	pkgcatalog "github.com/matrixorigin/matrixone/pkg/catalog"
	"github.com/matrixorigin/matrixone/pkg/container/types"
	"github.com/matrixorigin/matrixone/pkg/vm/engine/tae/catalog"
	"github.com/matrixorigin/matrixone/pkg/vm/engine/tae/txn/txnbase"
)

type Scope = int

const (
	// changes for mo_databases
	ScopeDatabases Scope = iota + 1
	// changes for mo_tables
	ScopeTables
	// changes for mo_columns
	ScopeColumns
	// changes for user tables
	ScopeUserTables
)

const (
	SnapshotAttr_TID                            = catalog.SnapshotAttr_TID
	SnapshotAttr_DBID                           = catalog.SnapshotAttr_DBID
	SegmentAttr_ID                              = catalog.SegmentAttr_ID
	SegmentAttr_CreateAt                        = catalog.SegmentAttr_CreateAt
	SegmentAttr_SegNode                         = catalog.SegmentAttr_SegNode
	SnapshotAttr_BlockMaxRow                    = catalog.SnapshotAttr_BlockMaxRow
	SnapshotAttr_SegmentMaxBlock                = catalog.SnapshotAttr_SegmentMaxBlock
	SnapshotMetaAttr_BlockInsertBatchStart      = "block_insert_batch_start"
	SnapshotMetaAttr_BlockInsertBatchEnd        = "block_insert_batch_end"
	SnapshotMetaAttr_BlockInsertBatchLocation   = "block_insert_batch_location"
	SnapshotMetaAttr_BlockDeleteBatchStart      = "block_delete_batch_start"
	SnapshotMetaAttr_BlockDeleteBatchEnd        = "block_delete_batch_end"
	SnapshotMetaAttr_BlockDeleteBatchLocation   = "block_delete_batch_location"
	SnapshotMetaAttr_BlockCNInsertBatchLocation = "block_cn_insert_batch_location"
	SnapshotMetaAttr_SegDeleteBatchStart        = "seg_delete_batch_start"
	SnapshotMetaAttr_SegDeleteBatchEnd          = "seg_delete_batch_end"
	SnapshotMetaAttr_SegDeleteBatchLocation     = "seg_delete_batch_location"
	CheckpointMetaAttr_BlockLocation            = "checkpoint_meta_block_location"
	CheckpointMetaAttr_SchemaType               = "checkpoint_meta_schema_type"

	AccountIDDbNameTblName = catalog.AccountIDDbNameTblName
	AccountIDDbName        = catalog.AccountIDDbName

	SnapshotAttr_SchemaExtra = catalog.SnapshotAttr_SchemaExtra
)

var (
	// for blk meta response
	BlkMetaSchema    *catalog.Schema // latest version
	BlkMetaSchema_V1 *catalog.Schema // previous version
	DelSchema        *catalog.Schema
	SegSchema        *catalog.Schema
	TxnNodeSchema    *catalog.Schema
	DBTNSchema       *catalog.Schema
	TblTNSchema      *catalog.Schema
	SegTNSchema      *catalog.Schema
	BlkTNSchema      *catalog.Schema
	MetaSchema_V1    *catalog.Schema
	MetaSchema       *catalog.Schema
	DBDelSchema      *catalog.Schema
	TblDelSchema     *catalog.Schema
	ColumnDelSchema  *catalog.Schema
	TNMetaSchema     *catalog.Schema

	DBSpecialDeleteSchema  *catalog.Schema
	TBLSpecialDeleteSchema *catalog.Schema
)

var (
	DBSpecialDeleteAttr = []string{
		pkgcatalog.SystemDBAttr_ID,
		AccountIDDbName,
	}
	DBSpecialDeleteTypes = []types.Type{
		types.New(types.T_uint64, 0, 0),
		types.New(types.T_varchar, 0, 0),
	}
	TBLSpecialDeleteAttr = []string{
		pkgcatalog.SystemRelAttr_ID,
		AccountIDDbNameTblName,
	}
	TBLSpecialDeleteTypes = []types.Type{
		types.New(types.T_uint64, 0, 0),
		types.New(types.T_varchar, 0, 0),
	}
	SegmentSchemaAttr = []string{
		SegmentAttr_ID,
		SegmentAttr_CreateAt,
		SegmentAttr_SegNode,
	}
	SegmentSchemaTypes = []types.Type{
		types.New(types.T_uuid, 0, 0),
		types.New(types.T_TS, 0, 0),
		types.New(types.T_blob, 0, 0),
	}
	TxnNodeSchemaAttr = []string{
		txnbase.SnapshotAttr_LogIndex_LSN,
		txnbase.SnapshotAttr_StartTS,
		txnbase.SnapshotAttr_PrepareTS,
		txnbase.SnapshotAttr_CommitTS,
		txnbase.SnapshotAttr_LogIndex_CSN,
		txnbase.SnapshotAttr_LogIndex_Size,
	}
	TxnNodeSchemaTypes = []types.Type{
		types.New(types.T_uint64, 0, 0),
		types.New(types.T_TS, 0, 0),
		types.New(types.T_TS, 0, 0),
		types.New(types.T_TS, 0, 0),
		types.New(types.T_uint32, 0, 0),
		types.New(types.T_uint32, 0, 0),
	}
	DBTNSchemaAttr = []string{
		txnbase.SnapshotAttr_LogIndex_LSN,
		txnbase.SnapshotAttr_StartTS,
		txnbase.SnapshotAttr_PrepareTS,
		txnbase.SnapshotAttr_CommitTS,
		txnbase.SnapshotAttr_LogIndex_CSN,
		txnbase.SnapshotAttr_LogIndex_Size,
		SnapshotAttr_DBID,
		SnapshotAttr_TID,
	}
	DBTNSchemaType = []types.Type{
		types.New(types.T_uint64, 0, 0),
		types.New(types.T_TS, 0, 0),
		types.New(types.T_TS, 0, 0),
		types.New(types.T_TS, 0, 0),
		types.New(types.T_uint32, 0, 0),
		types.New(types.T_uint32, 0, 0),
		types.New(types.T_uint64, 0, 0),
		types.New(types.T_uint64, 0, 0),
	}
	TblTNSchemaAttr = []string{
		txnbase.SnapshotAttr_LogIndex_LSN,
		txnbase.SnapshotAttr_StartTS,
		txnbase.SnapshotAttr_PrepareTS,
		txnbase.SnapshotAttr_CommitTS,
		txnbase.SnapshotAttr_LogIndex_CSN,
		txnbase.SnapshotAttr_LogIndex_Size,
		SnapshotAttr_DBID,
		SnapshotAttr_TID,
		SnapshotAttr_BlockMaxRow,
		SnapshotAttr_SegmentMaxBlock,
		SnapshotAttr_SchemaExtra,
	}
	TblTNSchemaType = []types.Type{
		types.New(types.T_uint64, 0, 0),
		types.New(types.T_TS, 0, 0),
		types.New(types.T_TS, 0, 0),
		types.New(types.T_TS, 0, 0),
		types.New(types.T_uint32, 0, 0),
		types.New(types.T_uint32, 0, 0),
		types.New(types.T_uint64, 0, 0),
		types.New(types.T_uint64, 0, 0),
		types.New(types.T_uint32, 0, 0),
		types.New(types.T_uint16, 0, 0),
		types.New(types.T_varchar, 0, 0),
	}
	SegmentTNSchemaAttr = []string{
		txnbase.SnapshotAttr_LogIndex_LSN,
		txnbase.SnapshotAttr_StartTS,
		txnbase.SnapshotAttr_PrepareTS,
		txnbase.SnapshotAttr_CommitTS,
		txnbase.SnapshotAttr_LogIndex_CSN,
		txnbase.SnapshotAttr_LogIndex_Size,
		SnapshotAttr_DBID,
		SnapshotAttr_TID,
	}
	SegmentTNSchemaTypes = []types.Type{
		types.New(types.T_uint64, 0, 0),
		types.New(types.T_TS, 0, 0),
		types.New(types.T_TS, 0, 0),
		types.New(types.T_TS, 0, 0),
		types.New(types.T_uint32, 0, 0),
		types.New(types.T_uint32, 0, 0),
		types.New(types.T_uint64, 0, 0),
		types.New(types.T_uint64, 0, 0),
	}
	BlockTNSchemaAttr = []string{
		txnbase.SnapshotAttr_LogIndex_LSN,
		txnbase.SnapshotAttr_StartTS,
		txnbase.SnapshotAttr_PrepareTS,
		txnbase.SnapshotAttr_CommitTS,
		txnbase.SnapshotAttr_LogIndex_CSN,
		txnbase.SnapshotAttr_LogIndex_Size,
		SnapshotAttr_DBID,
		SnapshotAttr_TID,
		pkgcatalog.BlockMeta_MetaLoc,
		pkgcatalog.BlockMeta_DeltaLoc,
	}
	BlockTNSchemaTypes = []types.Type{
		types.New(types.T_uint64, 0, 0),
		types.New(types.T_TS, 0, 0),
		types.New(types.T_TS, 0, 0),
		types.New(types.T_TS, 0, 0),
		types.New(types.T_uint32, 0, 0),
		types.New(types.T_uint32, 0, 0),
		types.New(types.T_uint64, 0, 0),
		types.New(types.T_uint64, 0, 0),
		types.New(types.T_varchar, types.MaxVarcharLen, 0),
		types.New(types.T_varchar, types.MaxVarcharLen, 0),
	}
	MetaSchemaAttr_V1 = []string{
		SnapshotAttr_TID,
		SnapshotMetaAttr_BlockInsertBatchStart,
		SnapshotMetaAttr_BlockInsertBatchEnd,
		SnapshotMetaAttr_BlockDeleteBatchStart,
		SnapshotMetaAttr_BlockDeleteBatchEnd,
		SnapshotMetaAttr_SegDeleteBatchStart,
		SnapshotMetaAttr_SegDeleteBatchEnd,
	}
	MetaShcemaTypes_V1 = []types.Type{
		types.New(types.T_uint64, 0, 0),
		types.New(types.T_int32, 0, 0),
		types.New(types.T_int32, 0, 0),
		types.New(types.T_int32, 0, 0),
		types.New(types.T_int32, 0, 0),
		types.New(types.T_int32, 0, 0),
		types.New(types.T_int32, 0, 0),
	}
	MetaSchemaAttr = []string{
		SnapshotAttr_TID,
		SnapshotMetaAttr_BlockInsertBatchLocation,
		SnapshotMetaAttr_BlockCNInsertBatchLocation,
		SnapshotMetaAttr_BlockDeleteBatchLocation,
		SnapshotMetaAttr_SegDeleteBatchLocation,
	}
	MetaShcemaTypes = []types.Type{
		types.New(types.T_uint64, 0, 0),
		types.New(types.T_varchar, types.MaxVarcharLen, 0),
		types.New(types.T_varchar, types.MaxVarcharLen, 0),
		types.New(types.T_varchar, types.MaxVarcharLen, 0),
		types.New(types.T_varchar, types.MaxVarcharLen, 0),
	}
	DBDelSchemaAttr = []string{
		pkgcatalog.SystemDBAttr_ID,
	}
	DBDelSchemaTypes = []types.Type{
		types.T_uint64.ToType(),
	}
	TblDelSchemaAttr = []string{
		pkgcatalog.SystemRelAttr_ID,
	}
	TblDelSchemaTypes = []types.Type{
		types.T_uint64.ToType(),
	}
	ColumnDelSchemaAttr = []string{
		pkgcatalog.SystemColAttr_UniqName,
	}
	ColumnDelSchemaTypes = []types.Type{
		types.T_varchar.ToType(),
	}
	TNMetaSchemaAttr = []string{
		CheckpointMetaAttr_BlockLocation,
		CheckpointMetaAttr_SchemaType,
	}
	TNMetaShcemaTypes = []types.Type{
		types.New(types.T_varchar, types.MaxVarcharLen, 0),
		types.New(types.T_uint16, 0, 0),
	}

	BaseAttr = []string{
		catalog.AttrRowID,
		catalog.AttrCommitTs,
	}
	BaseTypes = []types.Type{
		types.T_Rowid.ToType(),
		types.T_TS.ToType(),
	}
)

var (
	skippedDB = map[uint64]bool{
		pkgcatalog.MO_CATALOG_ID: true,
	}

	skippedTbl = map[uint64]bool{
		pkgcatalog.MO_DATABASE_ID: true,
		pkgcatalog.MO_TABLES_ID:   true,
		pkgcatalog.MO_COLUMNS_ID:  true,
	}
)

func shouldIgnoreDBInLogtail(id uint64) bool {
	return skippedDB[id]
}

func shouldIgnoreTblInLogtail(id uint64) bool {
	return skippedTbl[id]
}

func init() {

	DBSpecialDeleteSchema = catalog.NewEmptySchema("db_special_delete")

	for i, colname := range DBSpecialDeleteAttr {
		if i == 0 {
			if err := DBSpecialDeleteSchema.AppendPKCol(colname, DBSpecialDeleteTypes[i], 0); err != nil {
				panic(err)
			}
		} else {
			if err := DBSpecialDeleteSchema.AppendCol(colname, DBSpecialDeleteTypes[i]); err != nil {
				panic(err)
			}
		}
	}
	if err := DBSpecialDeleteSchema.Finalize(true); err != nil { // no phyaddr column
		panic(err)
	}

	TBLSpecialDeleteSchema = catalog.NewEmptySchema("tbl_special_delete")

	for i, colname := range TBLSpecialDeleteAttr {
		if i == 0 {
			if err := TBLSpecialDeleteSchema.AppendPKCol(colname, TBLSpecialDeleteTypes[i], 0); err != nil {
				panic(err)
			}
		} else {
			if err := TBLSpecialDeleteSchema.AppendCol(colname, TBLSpecialDeleteTypes[i]); err != nil {
				panic(err)
			}
		}
	}
	if err := TBLSpecialDeleteSchema.Finalize(true); err != nil { // no phyaddr column
		panic(err)
	}

	BlkMetaSchema = catalog.NewEmptySchema("blkMeta")

	for i, colname := range pkgcatalog.MoTableMetaSchema {
		if i == 0 {
			if err := BlkMetaSchema.AppendPKCol(colname, pkgcatalog.MoTableMetaTypes[i], 0); err != nil {
				panic(err)
			}
		} else {
			if err := BlkMetaSchema.AppendCol(colname, pkgcatalog.MoTableMetaTypes[i]); err != nil {
				panic(err)
			}
		}
	}
	if err := BlkMetaSchema.Finalize(true); err != nil { // no phyaddr column
		panic(err)
	}

	BlkMetaSchema_V1 = catalog.NewEmptySchema("blkMetaV1")

	for i, colname := range pkgcatalog.MoTableMetaSchemaV1 {
		if i == 0 {
			if err := BlkMetaSchema_V1.AppendPKCol(colname, pkgcatalog.MoTableMetaTypesV1[i], 0); err != nil {
				panic(err)
			}
		} else {
			if err := BlkMetaSchema_V1.AppendCol(colname, pkgcatalog.MoTableMetaTypesV1[i]); err != nil {
				panic(err)
			}
		}
	}
	if err := BlkMetaSchema_V1.Finalize(true); err != nil { // no phyaddr column
		panic(err)
	}

	// empty schema, no finalize, makeRespBatchFromSchema will add necessary colunms
	DelSchema = catalog.NewEmptySchema("del")

	SegSchema = catalog.NewEmptySchema("segment")
	for i, colname := range SegmentSchemaAttr {
		if i == 0 {
			if err := SegSchema.AppendPKCol(colname, SegmentSchemaTypes[i], 0); err != nil {
				panic(err)
			}
		} else {
			if err := SegSchema.AppendCol(colname, SegmentSchemaTypes[i]); err != nil {
				panic(err)
			}
		}
	}

	TxnNodeSchema = catalog.NewEmptySchema("txn_node")
	for i, colname := range TxnNodeSchemaAttr {
		if i == 0 {
			if err := TxnNodeSchema.AppendPKCol(colname, TxnNodeSchemaTypes[i], 0); err != nil {
				panic(err)
			}
		} else {
			if err := TxnNodeSchema.AppendCol(colname, TxnNodeSchemaTypes[i]); err != nil {
				panic(err)
			}
		}
	}

	DBTNSchema = catalog.NewEmptySchema("db_dn")
	for i, colname := range DBTNSchemaAttr {
		if i == 0 {
			if err := DBTNSchema.AppendPKCol(colname, DBTNSchemaType[i], 0); err != nil {
				panic(err)
			}
		} else {
			if err := DBTNSchema.AppendCol(colname, DBTNSchemaType[i]); err != nil {
				panic(err)
			}
		}
	}

	TblTNSchema = catalog.NewEmptySchema("table_dn")
	for i, colname := range TblTNSchemaAttr {
		if i == 0 {
			if err := TblTNSchema.AppendPKCol(colname, TblTNSchemaType[i], 0); err != nil {
				panic(err)
			}
		} else {
			if err := TblTNSchema.AppendCol(colname, TblTNSchemaType[i]); err != nil {
				panic(err)
			}
		}
	}

	SegTNSchema = catalog.NewEmptySchema("segment_dn")
	for i, colname := range SegmentTNSchemaAttr {
		if i == 0 {
			if err := SegTNSchema.AppendPKCol(colname, SegmentTNSchemaTypes[i], 0); err != nil {
				panic(err)
			}
		} else {
			if err := SegTNSchema.AppendCol(colname, SegmentTNSchemaTypes[i]); err != nil {
				panic(err)
			}
		}
	}

	BlkTNSchema = catalog.NewEmptySchema("block_dn")
	for i, colname := range BlockTNSchemaAttr {
		if i == 0 {
			if err := BlkTNSchema.AppendPKCol(colname, BlockTNSchemaTypes[i], 0); err != nil {
				panic(err)
			}
		} else {
			if err := BlkTNSchema.AppendCol(colname, BlockTNSchemaTypes[i]); err != nil {
				panic(err)
			}
		}
	}

	MetaSchema = catalog.NewEmptySchema("meta")
	for i, colname := range MetaSchemaAttr {
		if i == 0 {
			if err := MetaSchema.AppendPKCol(colname, MetaShcemaTypes[i], 0); err != nil {
				panic(err)
			}
		} else {
			if err := MetaSchema.AppendCol(colname, MetaShcemaTypes[i]); err != nil {
				panic(err)
			}
		}
	}

	MetaSchema_V1 = catalog.NewEmptySchema("meta")
	for i, colname := range MetaSchemaAttr_V1 {
		if i == 0 {
			if err := MetaSchema_V1.AppendPKCol(colname, MetaShcemaTypes_V1[i], 0); err != nil {
				panic(err)
			}
		} else {
			if err := MetaSchema_V1.AppendCol(colname, MetaShcemaTypes_V1[i]); err != nil {
				panic(err)
			}
		}
	}

	DBDelSchema = catalog.NewEmptySchema("meta")
	for i, colname := range DBDelSchemaAttr {
		if i == 0 {
			if err := DBDelSchema.AppendPKCol(colname, DBDelSchemaTypes[i], 0); err != nil {
				panic(err)
			}
		} else {
			if err := DBDelSchema.AppendCol(colname, DBDelSchemaTypes[i]); err != nil {
				panic(err)
			}
		}
	}

	TblDelSchema = catalog.NewEmptySchema("meta")
	for i, colname := range TblDelSchemaAttr {
		if i == 0 {
			if err := TblDelSchema.AppendPKCol(colname, TblDelSchemaTypes[i], 0); err != nil {
				panic(err)
			}
		} else {
			if err := TblDelSchema.AppendCol(colname, TblDelSchemaTypes[i]); err != nil {
				panic(err)
			}
		}
	}

	ColumnDelSchema = catalog.NewEmptySchema("meta")
	for i, colname := range ColumnDelSchemaAttr {
		if i == 0 {
			if err := ColumnDelSchema.AppendPKCol(colname, ColumnDelSchemaTypes[i], 0); err != nil {
				panic(err)
			}
		} else {
			if err := ColumnDelSchema.AppendCol(colname, ColumnDelSchemaTypes[i]); err != nil {
				panic(err)
			}
		}
	}

	TNMetaSchema = catalog.NewEmptySchema("meta")
	for i, colname := range TNMetaSchemaAttr {
		if i == 0 {
			if err := TNMetaSchema.AppendPKCol(colname, TNMetaShcemaTypes[i], 0); err != nil {
				panic(err)
			}
		} else {
			if err := TNMetaSchema.AppendCol(colname, TNMetaShcemaTypes[i]); err != nil {
				panic(err)
			}
		}
	}
}
