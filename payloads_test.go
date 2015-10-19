// Copyright 2015 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package charm_test

import (
	"fmt"

	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"

	"gopkg.in/juju/charm.v6-unstable"
)

var _ = gc.Suite(&payloadClassSuite{})

type payloadClassSuite struct{}

func (s *payloadClassSuite) TestParsePayloadClassOkay(c *gc.C) {
	name := "my-payload"
	data := map[string]interface{}{
		"type": "docker",
	}
	payloadClass := charm.ParsePayloadClass(name, data)

	c.Check(payloadClass, jc.DeepEquals, charm.PayloadClass{
		Name: "my-payload",
		Type: "docker",
	})
}

func (s *payloadClassSuite) TestParsePayloadClassMissingName(c *gc.C) {
	name := ""
	data := map[string]interface{}{
		"type": "docker",
	}
	payloadClass := charm.ParsePayloadClass(name, data)

	c.Check(payloadClass, jc.DeepEquals, charm.PayloadClass{
		Name: "",
		Type: "docker",
	})
}

func (s *payloadClassSuite) TestParsePayloadClassEmpty(c *gc.C) {
	name := "my-payload"
	var data map[string]interface{}
	payloadClass := charm.ParsePayloadClass(name, data)

	c.Check(payloadClass, jc.DeepEquals, charm.PayloadClass{
		Name: "my-payload",
	})
}

func (s *payloadClassSuite) TestValidateFull(c *gc.C) {
	payloadClass := charm.PayloadClass{
		Name: "my-payload",
		Type: "docker",
	}
	err := payloadClass.Validate()

	c.Check(err, jc.ErrorIsNil)
}

func (s *payloadClassSuite) TestValidateValidName(c *gc.C) {
	for _, name := range []string{
		"mypayload",
		"my-payload",
		"my_payload",
		"my0-payload",
		"my-payload0",
	} {
		checkPayloadName(c, name, "")
	}
}

func (s *payloadClassSuite) TestValidateZeroValue(c *gc.C) {
	var payloadClass charm.PayloadClass
	err := payloadClass.Validate()

	c.Check(err, gc.NotNil)
}

func (s *payloadClassSuite) TestValidateMissingName(c *gc.C) {
	payloadClass := charm.PayloadClass{
		Type: "docker",
	}
	err := payloadClass.Validate()

	c.Check(err, gc.ErrorMatches, `payload class missing name`)
}

const (
	punctuation       = `!"#$%&'()*+,-./:;<=>?@[\]^_` + "`" + `{|}~`
	payloadInvalidMsg = `payload class name .* not valid`
)

func (s *payloadClassSuite) TestValidateInvalidNameStart(c *gc.C) {
	msg := payloadInvalidMsg
	for _, chr := range []byte(punctuation + `0123456789`) {
		checkPayloadName(c, fmt.Sprintf("%cmy-payload", chr), msg)
	}
}

func (s *payloadClassSuite) TestValidateInvalidNameMiddle(c *gc.C) {
	msg := payloadInvalidMsg
	for _, chr := range []byte(punctuation) {
		if chr == '-' || chr == '_' {
			continue
		}
		checkPayloadName(c, fmt.Sprintf("my%cpayload", chr), msg)
	}
}

func (s *payloadClassSuite) TestValidateInvalidNameEnd(c *gc.C) {
	msg := payloadInvalidMsg
	for _, chr := range []byte(punctuation) {
		checkPayloadName(c, fmt.Sprintf("my-payload%c", chr), msg)
	}
	checkPayloadName(c, "my-payload-", msg)
	checkPayloadName(c, "my-payload_", msg)
}

func (s *payloadClassSuite) TestValidateMissingType(c *gc.C) {
	payloadClass := charm.PayloadClass{
		Name: "my-payload",
	}
	err := payloadClass.Validate()

	c.Check(err, gc.ErrorMatches, `payload class missing type`)
}

func checkPayloadName(c *gc.C, name, msg string) {
	c.Logf("checking payload name %q", name)

	payloadClass := charm.PayloadClass{
		Name: name,
		Type: "docker",
	}
	err := payloadClass.Validate()

	if msg == "" {
		c.Check(err, jc.ErrorIsNil)
	} else {
		c.Check(err, gc.ErrorMatches, msg)
	}
}
