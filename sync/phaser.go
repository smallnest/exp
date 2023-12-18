package sync

import (
	"sync"
	"sync/atomic"
)

// Phaser is a reusable synchronization barrier, similar in functionality to java Phaser.
type Phaser struct {
	parties    atomic.Int32
	arrived    atomic.Int32
	phase      atomic.Int32
	barrier    *sync.Cond
	terminated atomic.Int32
}

// NewPhaser creates a new Phaser instance.
func NewPhaser(parties int32) *Phaser {
	var p Phaser
	p.parties.Store(parties)
	p.barrier = sync.NewCond(&sync.Mutex{})
	return &p
}

// Join adds a new party to this phaser.
// Just like java.util.concurrent.Phaser's register() method.
func (p *Phaser) Join() int32 {
	p.barrier.L.Lock()
	defer p.barrier.L.Unlock()

	return p.parties.Add(1)
}

// BulkJoin adds a number of new parties to this phaser.
// Just like java.util.concurrent.Phaser's bulkRegister(int parties) method.
func (p *Phaser) BulkJoin(parties int32) int32 {
	p.barrier.L.Lock()
	defer p.barrier.L.Unlock()

	return p.parties.Add(parties)
}

// Arrive arrives at this phaser, without waiting for others to arrive.
func (p *Phaser) Arrive() int32 {
	p.barrier.L.Lock()
	defer p.barrier.L.Unlock()

	currentArrived := p.arrived.Add(1)
	if currentArrived == p.parties.Load() { // all arrived
		p.phase.Add(1)
		p.arrived.Store(0)
		p.barrier.Broadcast()
	}

	return p.phase.Load()
}

// Wait awaits the phase of this phaser to advance from the given phase value,
// returning immediately if the current phase is not equal to the given phase value or this phaser is terminated.
// Just like java.util.concurrent.Phaser's awaitAdvance(int phase) method.
func (p *Phaser) Wait(phase int) int32 {
	p.barrier.L.Lock()
	defer p.barrier.L.Unlock()

	for int32(phase) == p.phase.Load() && p.terminated.Load() == 0 {
		p.barrier.Wait()
	}

	return p.phase.Load()
}

// ArriveAndWait arrives at this phaser and waits others.
// Just like java.util.concurrent.Phaser's arriveAndAwaitAdvance() method.
func (p *Phaser) ArriveAndWait() int32 {
	p.barrier.L.Lock()
	defer p.barrier.L.Unlock()

	phase := p.phase.Load()
	currentArrived := p.arrived.Add(1)
	if currentArrived == p.parties.Load() { // all arrived
		p.phase.Add(1)
		p.arrived.Store(0)
		p.barrier.Broadcast()
	} else {
		// wait for phase to change in current phase
		for phase == p.phase.Load() && p.terminated.Load() == 0 {
			p.barrier.Wait()
		}
	}

	return p.phase.Load()
}

// ArriveAndLeave arrives at this phaser and leaves from it without waiting for others to arrive.
// Just like java.util.concurrent.Phaser's arriveAndDeregister() method.
func (p *Phaser) ArriveAndLeave() int32 {
	p.barrier.L.Lock()
	defer p.barrier.L.Unlock()

	phase := p.phase.Load()
	currentArrived := p.arrived.Add(1)
	if currentArrived == p.parties.Load() { // all arrived
		p.phase.Add(1)
		p.arrived.Store(0)
		p.barrier.Broadcast()
	} else {
		// wait for phase to change in current phase
		for phase == p.phase.Load() && p.terminated.Load() == 0 {
			p.barrier.Wait()
		}
	}

	p.leave()

	return p.phase.Load()
}

func (p *Phaser) leave() int32 {
	// leave this phaser
	parties := p.parties.Add(-1)
	if parties == 0 { // is the last one, terminate this phaser
		p.parties.Store(0)
		p.phase.Store(0)
		p.arrived.Store(0)
		p.terminated.Store(1)
		p.barrier.Broadcast()
	}

	return parties
}

// Leave leaves from this phaser without waiting for others to arrive.
// Just like java.util.concurrent.Phaser's deregister() method.
func (p *Phaser) Leave() int32 {
	p.barrier.L.Lock()
	defer p.barrier.L.Unlock()

	return p.leave()
}

// Phase returns the current phase number.
func (p *Phaser) Phase() int32 {
	return p.phase.Load()
}

// Phase returns the current phase number.
func (p *Phaser) Arrived() int32 {
	return p.arrived.Load()
}

// Parties returns the number of parties joined in this phaser.
func (p *Phaser) Parties() int32 {
	return p.parties.Load()
}

// ForceTermination forces this phaser to enter termination state.
func (p *Phaser) ForceTermination() {
	p.terminated.Store(1)
}

// IsTerminated returns true if this phaser has been terminated.
func (p *Phaser) IsTerminated() bool {
	return p.terminated.Load() == 1
}
