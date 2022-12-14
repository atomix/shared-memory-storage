// SPDX-FileCopyrightText: 2022-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package sharedmemory

import (
	"context"
	counterv1 "github.com/atomix/runtime/primitives/pkg/counter/v1"
	countermapv1 "github.com/atomix/runtime/primitives/pkg/countermap/v1"
	electionv1 "github.com/atomix/runtime/primitives/pkg/election/v1"
	indexedmapv1 "github.com/atomix/runtime/primitives/pkg/indexedmap/v1"
	lockv1 "github.com/atomix/runtime/primitives/pkg/lock/v1"
	mapv1 "github.com/atomix/runtime/primitives/pkg/map/v1"
	multimapv1 "github.com/atomix/runtime/primitives/pkg/multimap/v1"
	setv1 "github.com/atomix/runtime/primitives/pkg/set/v1"
	valuev1 "github.com/atomix/runtime/primitives/pkg/value/v1"
	"github.com/atomix/runtime/sdk/pkg/protocol"
	"github.com/atomix/runtime/sdk/pkg/protocol/node"
	"github.com/atomix/runtime/sdk/pkg/protocol/statemachine"
	streams "github.com/atomix/runtime/sdk/pkg/stream"
	"sync"
)

func NewProtocol() node.Protocol {
	return &memoryProtocol{
		partitions: map[protocol.PartitionID]node.Partition{
			1: node.NewPartition(1, newExecutor()),
			2: node.NewPartition(2, newExecutor()),
			3: node.NewPartition(3, newExecutor()),
		},
	}
}

type memoryProtocol struct {
	partitions map[protocol.PartitionID]node.Partition
}

func (p *memoryProtocol) Partitions() []node.Partition {
	partitions := make([]node.Partition, 0, len(p.partitions))
	for _, partition := range p.partitions {
		partitions = append(partitions, partition)
	}
	return partitions
}

func (p *memoryProtocol) Partition(partitionID protocol.PartitionID) (node.Partition, bool) {
	partition, ok := p.partitions[partitionID]
	return partition, ok
}

func newExecutor() node.Executor {
	registry := statemachine.NewPrimitiveTypeRegistry()
	counterv1.RegisterStateMachine(registry)
	countermapv1.RegisterStateMachine(registry)
	electionv1.RegisterStateMachine(registry)
	indexedmapv1.RegisterStateMachine(registry)
	lockv1.RegisterStateMachine(registry)
	mapv1.RegisterStateMachine(registry)
	multimapv1.RegisterStateMachine(registry)
	setv1.RegisterStateMachine(registry)
	valuev1.RegisterStateMachine(registry)
	return &memoryExecutor{
		sm: statemachine.NewStateMachine(registry),
	}
}

type memoryExecutor struct {
	sm statemachine.StateMachine
	mu sync.RWMutex
}

func (e *memoryExecutor) Propose(ctx context.Context, proposal *protocol.ProposalInput, stream streams.WriteStream[*protocol.ProposalOutput]) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.sm.Propose(proposal, stream)
	return nil
}

func (e *memoryExecutor) Query(ctx context.Context, query *protocol.QueryInput, stream streams.WriteStream[*protocol.QueryOutput]) error {
	e.mu.RLock()
	defer e.mu.RUnlock()
	e.sm.Query(query, stream)
	return nil
}
