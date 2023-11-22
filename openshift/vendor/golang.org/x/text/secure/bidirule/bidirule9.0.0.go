<<<<<<< HEAD
// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !go1.10

package bidirule

func (t *Transformer) isFinal() bool {
	if !t.isRTL() {
		return true
	}
	return t.state == ruleLTRFinal || t.state == ruleRTLFinal || t.state == ruleInitial
}
||||||| parent of b907b2097 (Add cluster-capi-operator integration)
=======
// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !go1.10
// +build !go1.10

package bidirule

func (t *Transformer) isFinal() bool {
	if !t.isRTL() {
		return true
	}
	return t.state == ruleLTRFinal || t.state == ruleRTLFinal || t.state == ruleInitial
}
>>>>>>> b907b2097 (Add cluster-capi-operator integration)
