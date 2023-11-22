<<<<<<< HEAD
// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build gc

#include "textflag.h"

//
// System calls for amd64, Solaris are implemented in runtime/syscall_solaris.go
//

TEXT ·sysvicall6(SB),NOSPLIT,$0-88
	JMP	syscall·sysvicall6(SB)

TEXT ·rawSysvicall6(SB),NOSPLIT,$0-88
	JMP	syscall·rawSysvicall6(SB)
||||||| parent of b907b2097 (Add cluster-capi-operator integration)
=======
// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build gc
// +build gc

#include "textflag.h"

//
// System calls for amd64, Solaris are implemented in runtime/syscall_solaris.go
//

TEXT ·sysvicall6(SB),NOSPLIT,$0-88
	JMP	syscall·sysvicall6(SB)

TEXT ·rawSysvicall6(SB),NOSPLIT,$0-88
	JMP	syscall·rawSysvicall6(SB)
>>>>>>> b907b2097 (Add cluster-capi-operator integration)
