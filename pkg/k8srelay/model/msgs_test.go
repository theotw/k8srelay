/*
 * Copyright (c) The One True Way 2023. Apache License 2.0. The authors accept no liability, 0 nada for the use of this software.  It is offered "As IS"  Have fun with it!!
 */

package model

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestMessages(t *testing.T) {
	req := NewCallReq()
	req.AddHeader("k", "v")
	x := req.Headers["k"]
	assert.Equal(t, "v", x)

	resp := NewCallResponse()
	resp.AddHeader("k", "1")
	x = resp.Headers["k"]
	assert.Equal(t, "1", x)

}
func TestMessagesHelpers(t *testing.T) {
	subj := MakeMessageSubject("1", "2")
	assert.True(t, strings.HasSuffix(subj, "1.2"))
	subj = MakeMessageSubject("1", "")
	assert.True(t, strings.HasSuffix(subj, "1"))

	subj = MakeReplySubject("1")
	assert.True(t, strings.HasPrefix(subj, NatssyncMessagePrefix+".1"))
	subj = MakeNBReplySubject()
	assert.True(t, strings.HasPrefix(subj, NatssyncMessagePrefix+"."+CloudID))
}
