// Copyright 2009 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Mach-O header data structures
// http://developer.apple.com/mac/library/documentation/DeveloperTools/Conceptual/MachORuntime/Reference/reference.html

package macho

import "strconv"

// A FileHeader represents a Mach-O file header.
type FileHeader struct {
	Magic  uint32
	Cpu    Cpu
	SubCpu uint32
	Type   Type
	Ncmd   uint32
	Cmdsz  uint32
	Flags  uint32
}

const (
	fileHeaderSize32 = 7 * 4
	fileHeaderSize64 = 8 * 4
)

const (
	Magic32 uint32 = 0xfeedface
	Magic64 uint32 = 0xfeedfacf
)

// A Type is a Mach-O file type, either an object or an executable.
type Type uint32

const (
	TypeObj  Type = 1
	TypeExec Type = 2
)

// A Cpu is a Mach-O cpu type.
type Cpu uint32

const (
	Cpu386   Cpu = 7
	CpuAmd64 Cpu = Cpu386 + 1<<24
)

var cpuStrings = []intName{
	intName{uint32(Cpu386), "Cpu386"},
	intName{uint32(CpuAmd64), "CpuAmd64"},
}

func (i Cpu) String() string   { return stringName(uint32(i), cpuStrings, false) }
func (i Cpu) GoString() string { return stringName(uint32(i), cpuStrings, true) }

// A LoadCmd is a Mach-O load command.
type LoadCmd uint32

const (
	LoadCmdSegment    LoadCmd = 1
	LoadCmdSegment64  LoadCmd = 25
	LoadCmdThread     LoadCmd = 4
	LoadCmdUnixThread LoadCmd = 5 // thread+stack
)

var cmdStrings = []intName{
	intName{uint32(LoadCmdSegment), "LoadCmdSegment"},
	intName{uint32(LoadCmdSegment64), "LoadCmdSegment64"},
	intName{uint32(LoadCmdThread), "LoadCmdThread"},
	intName{uint32(LoadCmdUnixThread), "LoadCmdUnixThread"},
}

func (i LoadCmd) String() string   { return stringName(uint32(i), cmdStrings, false) }
func (i LoadCmd) GoString() string { return stringName(uint32(i), cmdStrings, true) }

// A Segment64 is a 64-bit Mach-O segment load command.
type Segment64 struct {
	Cmd     LoadCmd
	Len     uint32
	Name    [16]byte
	Addr    uint64
	Memsz   uint64
	Offset  uint64
	Filesz  uint64
	Maxprot uint32
	Prot    uint32
	Nsect   uint32
	Flag    uint32
}

// A Segment32 is a 32-bit Mach-O segment load command.
type Segment32 struct {
	Cmd     LoadCmd
	Len     uint32
	Name    [16]byte
	Addr    uint32
	Memsz   uint32
	Offset  uint32
	Filesz  uint32
	Maxprot uint32
	Prot    uint32
	Nsect   uint32
	Flag    uint32
}

// A Section32 is a 32-bit Mach-O section header.
type Section32 struct {
	Name     [16]byte
	Seg      [16]byte
	Addr     uint32
	Size     uint32
	Offset   uint32
	Align    uint32
	Reloff   uint32
	Nreloc   uint32
	Flags    uint32
	Reserve1 uint32
	Reserve2 uint32
}

// A Section32 is a 64-bit Mach-O section header.
type Section64 struct {
	Name     [16]byte
	Seg      [16]byte
	Addr     uint64
	Size     uint64
	Offset   uint32
	Align    uint32
	Reloff   uint32
	Nreloc   uint32
	Flags    uint32
	Reserve1 uint32
	Reserve2 uint32
	Reserve3 uint32
}

// A Thread is a Mach-O thread state command.
type Thread struct {
	Cmd  LoadCmd
	Len  uint32
	Type uint32
	Data []uint32
}

// Regs386 is the Mach-O 386 register structure.
type Regs386 struct {
	AX    uint32
	BX    uint32
	CX    uint32
	DX    uint32
	DI    uint32
	SI    uint32
	BP    uint32
	SP    uint32
	SS    uint32
	FLAGS uint32
	IP    uint32
	CS    uint32
	DS    uint32
	ES    uint32
	FS    uint32
	GS    uint32
}

// RegsAMD64 is the Mach-O AMD64 register structure.
type RegsAMD64 struct {
	AX    uint64
	BX    uint64
	CX    uint64
	DX    uint64
	DI    uint64
	SI    uint64
	BP    uint64
	SP    uint64
	R8    uint64
	R9    uint64
	R10   uint64
	R11   uint64
	R12   uint64
	R13   uint64
	R14   uint64
	R15   uint64
	IP    uint64
	FLAGS uint64
	CS    uint64
	FS    uint64
	GS    uint64
}

type intName struct {
	i uint32
	s string
}

func stringName(i uint32, names []intName, goSyntax bool) string {
	for _, n := range names {
		if n.i == i {
			if goSyntax {
				return "macho." + n.s
			}
			return n.s
		}
	}
	return strconv.Uitoa64(uint64(i))
}

func flagName(i uint32, names []intName, goSyntax bool) string {
	s := ""
	for _, n := range names {
		if n.i&i == n.i {
			if len(s) > 0 {
				s += "+"
			}
			if goSyntax {
				s += "macho."
			}
			s += n.s
			i -= n.i
		}
	}
	if len(s) == 0 {
		return "0x" + strconv.Uitob64(uint64(i), 16)
	}
	if i != 0 {
		s += "+0x" + strconv.Uitob64(uint64(i), 16)
	}
	return s
}
