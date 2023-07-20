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

// Register adds a new party to this phaser.
func (p *Phaser) Register() int32 {
	p.barrier.L.Lock()
	defer p.barrier.L.Unlock()

	return p.parties.Add(1)
}

func (p *Phaser) BulkRegister(parties int32) int32 {
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

// AwaitAdvance awaits the phase of this phaser to advance from the given phase value,
// returning immediately if the current phase is not equal to the given phase value or this phaser is terminated.
func (p *Phaser) AwaitAdvance(phase int) int32 {
	p.barrier.L.Lock()
	defer p.barrier.L.Unlock()

	for int32(phase) == p.phase.Load() && p.terminated.Load() == 0 {
		p.barrier.Wait()
	}

	return p.phase.Load()
}

// ArriveAndAwaitAdvance arrives at this phaser and awaits others.
func (p *Phaser) ArriveAndAwaitAdvance() int32 {
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

// ArriveAndDeregister arrives at this phaser and deregisters from it without waiting for others to arrive.
func (p *Phaser) ArriveAndDeregister() int32 {
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

	p.deregister()

	return p.phase.Load()
}

func (p *Phaser) deregister() int32 {
	// deregister
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

func (p *Phaser) Deregister() int32 {
	p.barrier.L.Lock()
	defer p.barrier.L.Unlock()

	return p.deregister()
}

// Phase returns the current phase number.
func (p *Phaser) Phase() int32 {
	return p.phase.Load()
}

// Phase returns the current phase number.
func (p *Phaser) Arrived() int32 {
	return p.arrived.Load()
}

// ForceTermination forces this phaser to enter termination state.
func (p *Phaser) ForceTermination() {
	p.terminated.Store(1)
}

// IsTerminated returns true if this phaser has been terminated.
func (p *Phaser) IsTerminated() bool {
	return p.terminated.Load() == 1
}
