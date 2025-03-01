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

package mergeoffset

import (
	"github.com/matrixorigin/matrixone/pkg/sql/colexec"

	"github.com/matrixorigin/matrixone/pkg/vm/process"
)

type container struct {
	colexec.ReceiverOperator
	seen uint64
}

type Argument struct {
	// Offset records the offset number of mergeOffset operator
	Offset uint64
	// ctr contains the attributes needn't do serialization work
	ctr *container
}

func (arg *Argument) Free(proc *process.Process, pipelineFailed bool, err error) {
	if arg.ctr != nil {
		arg.ctr.FreeMergeTypeOperator(pipelineFailed)
	}
}
