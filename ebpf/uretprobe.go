// Copyright 2024 The Inspektor Gadget authors
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
//
// This package is extracted from the Inspektor Gadget project (https://github.com/inspektor-gadget/inspektor-gadget/blob/c34e8b7f2a1ab9d19d905bc90534ccc5d87daf40/pkg/uprobetracer/tracer.go).

// Package ebpf provides utilities for fixed crash of uretprobes in Go programs.
// see the issues:
// - https://github.com/iovisor/bcc/issues/3034
// - https://github.com/golang/go/issues/22008
// - https://github.com/iovisor/bcc/issues/1320
package ebpf

import (
	"debug/elf"
	"errors"
	"fmt"
	"log"

	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/link"
)

// AttachUretprobe attaches a uretprobe to the function at attachPath and attachSymbol.
// The probe will be attached to all return instructions in the function.
//
// attachPath is the path to the ELF file containing the function.
// attachSymbol is the name of the function to attach the probe to.
// prog is the eBPF program to attach.
// ex is the executable to attach the probe to.
// opts are the options for the probe.
func AttachUretprobe(attachPath, attachSymbol string, prog *ebpf.Program, ex *link.Executable, opts *link.UprobeOptions) (link.Link, error) {
	funcs, err := getFunctions(attachPath, map[string]struct{}{attachSymbol: {}})
	if err != nil {
		return nil, fmt.Errorf("getting functions: %w", err)
	}

	var last link.Link
	for _, f := range funcs {
		for _, ret := range f.Returns {
			var opt link.UprobeOptions
			if opts != nil {
				opt = *opts
			}
			opt.Address = ret
			last, err = ex.Uprobe("", prog, &opt)
			if err != nil {
				return nil, fmt.Errorf("installing uprobe: %w", err)
			}
		}
	}
	return last, nil

}

// getFunctions returns a map of functions in the ELF file that are in the lookups map.
func getFunctions(filename string, lookups map[string]struct{}) (map[string]*Function, error) {
	elff, err := elf.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open ELF: %w", err)
	}
	defer elff.Close()

	functions := make(map[string]*Function)

	// find symbols
	syms, err := elff.Symbols()
	if err != nil && !errors.Is(err, elf.ErrNoSymbols) {
		return nil, fmt.Errorf("failed to get symbols: %w", err)
	}

	collectSymbols(elff, syms, functions, lookups)

	syms, err = elff.DynamicSymbols()
	if err != nil && !errors.Is(err, elf.ErrNoSymbols) {
		return nil, fmt.Errorf("failed to get dynamic symbols: %w", err)
	}

	collectSymbols(elff, syms, functions, lookups)

	for name, function := range functions {
		data := make([]byte, function.Size)
		_, err := function.Prog.ReadAt(data, int64(function.Offset-function.Prog.Off))
		if err != nil {
			return nil, fmt.Errorf("failed to read data for function %q: %w", name, err)
		}

		function.Returns, err = findReturnOffsets(function.Offset, data)
		if err != nil {
			return nil, fmt.Errorf("failed to get return offsets for function %q: %w", name, err)
		}
	}

	return functions, nil
}

// Function represents a function in an ELF file.
type Function struct {
	Offset  uint64
	Size    uint64
	Prog    *elf.Prog
	Returns []uint64
}

// collectSymbols collects symbols from the ELF file that are in the lookups map.
func collectSymbols(elff *elf.File, syms []elf.Symbol, target map[string]*Function, lookups map[string]struct{}) {
	for _, s := range syms {
		if elf.ST_TYPE(s.Info) != elf.STT_FUNC {
			continue
		}
		if _, ok := lookups[s.Name]; !ok {
			continue
		}
		address := s.Value

		log.Printf("found function %q", s.Name)
		var p *elf.Prog
		for _, prog := range elff.Progs {
			if prog.Type != elf.PT_LOAD || (prog.Flags&elf.PF_X) == 0 {
				continue
			}
			// stackoverflow.com/a/40249502
			if prog.Vaddr <= s.Value && s.Value < (prog.Vaddr+prog.Memsz) {
				address = s.Value - prog.Vaddr + prog.Off
				p = prog
				break
			}
		}
		target[s.Name] = &Function{Offset: address, Size: s.Size, Prog: p}
	}
}
