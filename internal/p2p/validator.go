package p2p

import (
	"bytes"
	"pnode/internal/apperror"
	"pnode/internal/prefix"
	"strings"
)

type Validator struct {
	b []byte
}

func (v *Validator) Bytes() []byte {
	return v.b
}

func (v *Validator) Validate() (d *CommandDTO, err error) {
	d = &CommandDTO{}

	if len(v.Bytes()) < prefix.LenMinimum {
		err = apperror.InvalidPrefixLength
		return
	}
	pref := v.b[0:3]
	cmdBody := v.b[4:]

	if bytes.Equal(pref, []byte(prefix.StringCMD)) {
		d.Type = prefix.TypeCMD
	} else if bytes.Equal(pref, []byte(prefix.StringP2P)) {
		d.Type = prefix.TypeP2P
	} else {
		err = apperror.InvalidPrefix
		return
	}

	split := strings.Split(string(cmdBody), " ")
	if len(split) < 1 {
		err = apperror.InvalidCommandBody
		return
	}

	d.Name = split[0]
	d.Args = split[1:]

	return
}

func NewValidator(b []byte) *Validator {
	return &Validator{b: b}
}
