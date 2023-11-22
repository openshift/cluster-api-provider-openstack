<<<<<<< HEAD
// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build arm && gc && linux

package unix

import "syscall"

// Underlying system call writes to newoffset via pointer.
// Implemented in assembly to avoid allocation.
func seek(fd int, offset int64, whence int) (newoffset int64, err syscall.Errno)
||||||| parent of b907b2097 (Add cluster-capi-operator integration)
=======
// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build arm && gc && linux
// +build arm,gc,linux

package unix

import "syscall"

// Underlying system call writes to newoffset via pointer.
// Implemented in assembly to avoid allocation.
func seek(fd int, offset int64, whence int) (newoffset int64, err syscall.Errno)
>>>>>>> b907b2097 (Add cluster-capi-operator integration)
