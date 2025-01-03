# exp

![GitHub](https://img.shields.io/github/license/smallnest/exp) ![GitHub Action](https://github.com/smallnest/exp/actions/workflows/action.yaml/badge.svg) [![Go Report Card](https://goreportcard.com/badge/github.com/smallnest/exp)](https://goreportcard.com/report/github.com/smallnest/exp)  [![GoDoc](https://godoc.org/github.com/smallnest/exp?status.png)](http://godoc.org/github.com/smallnest/exp)  



Experimental packages not in std and golang.org/exp


## packages

- **sync**
  - **generic sync.Map**: modify sync.Map to support generic
  - **Phaser**: a reusable synchronization barrier, similar in functionality to java.util.concurrent.Phaser
  - **Notifier**: implement the observer pattern via channel
  - **Shard**: a sharding data structure with lock-free read and write
  - **Exchanger**: a synchronization point at which goroutines can pair and swap elements within pairs. Each goroutine presents some object on entry to the exchange method, matches with a partner goroutine, and receives its partner's object on return. An Exchanger may be viewed as a bidirectional form of a channel.
  - **atomicx**: add C++ 20 atomic wait/notify feature for go std atomic
    - Pointer[T]
    - Value
    - Uintptr
    - Bool
    - Int32
    - Int64
    - Uint32
    - Uint64
  - **Pool**: a generic pool, forked from https://github.com/mkmik/syncpool
  - **Broadcaster**: a broadcaster to broadcast message to multiple receivers

- **container**
  - **heap**: generic heap
  - **binheap**: human friendly generic heap
  - **list**: generic list
  - **ring**: generic ring
  - **skiplist**: generic skiplist based on [mauricegit/skiplist](https://github.com/mauricegit/skiplist)
  - **set**: discussion at https://github.com/golang/go/discussions/47331
    - **Set**: generic set
    - **SortedSet**: generic sorted set 
  - **maps**:
    - **OrderedMap**: an insert-order map. The main code is forked [wk8/go-ordered-map](https://github.com/wk8/go-ordered-map)
    - **AccessOrderedMap**: an access-order map.
    - **BidiMap**: a bidirectional map. 
  - **Tuple**: a collection of generic tuples.
  - **bits**: a bitset implementation

- **chanx**
  - **Batch**: batch get from channels efficiently

- **io**
  - MmapReader: a mmap reader which implements io.Reader, io.ReadAt, io.Closer and can ReadLine

- **db**
  a group of sql helper functions.

- **stat**
  - `Hist` provides a Histogram.

- **ebpf**
  - **AttachUretprobe**: a helper function to add Uretprobe in Go programs to avoid crash

- internal data structure in go std libs
  - `PoolDequeue`: a lock-free, fixed-size single-producer, multi-consumer queue
  - `PoolChain`: a lock-free, dynamically-sized single-producer, multi-consumer queue

- mock
  - `sqlmock`: a simple sql mock
  - `kafkamock`: kafka mock