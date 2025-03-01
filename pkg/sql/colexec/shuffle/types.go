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

package shuffle

import (
	"github.com/matrixorigin/matrixone/pkg/container/batch"
	"github.com/matrixorigin/matrixone/pkg/vm/process"
)

const shuffleBatchSize = 8192

const (
	input = iota
	outPutNotEnding
	outPutEnding
)

type Argument struct {
	ctr           *container
	ShuffleColIdx int32
	ShuffleType   int32
	ShuffleColMin int64
	ShuffleColMax int64
	AliveRegCnt   int32
}

type container struct {
	state        int
	sels         [][]int32
	shuffledBats []*batch.Batch
}

func (arg *Argument) Free(proc *process.Process, pipelineFailed bool, err error) {
	for i := range arg.ctr.shuffledBats {
		if arg.ctr.shuffledBats[i] != nil {
			arg.ctr.shuffledBats[i].Clean(proc.Mp())
			arg.ctr.shuffledBats[i] = nil
		}
	}
	arg.ctr.sels = nil
}
