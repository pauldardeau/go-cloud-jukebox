package jukebox

import "testing"

func TestNewArgumentParser(t *testing.T) {
	ap := NewArgumentParser()
	if ap == nil {
		t.Log("NewArgumentParser should not return nil")
		t.Fail()
	}
}

func TestAddOptionalBoolFlag(t *testing.T) {
	th := NewTestHelper(t)
	ap := NewArgumentParser()
	ap.AddOptionalBoolFlag("--debug", "turn on debugging support")
	args := make([]string, 0)

	ps := ap.ParseArgs(args)
	th.Require(ps != nil, "ParseArgs should return non-nil")
	if ps != nil {
		th.RequireFalse(ps.Contains("debug"), "props must not contain optional argument that wasn't provided")
	}

	args = append(args, "--debug")
	ps = ap.ParseArgs(args)
	th.Require(ps != nil, "ParseArgs should return non-nil")
	if ps != nil {
		th.Require(ps.Contains("debug"), "props must contain provided argument")
		th.Require(ps.GetBoolValue("debug"), "bool true must be provided for optional bool arg that was provided")
	}
}

func TestAddOptionalIntArgument(t *testing.T) {
	th := NewTestHelper(t)
	ap := NewArgumentParser()
	ap.AddOptionalIntArgument("--logLevel", "adjust logging level up or down")

	args := make([]string, 0)
	ps := ap.ParseArgs(args)
	th.Require(ps != nil, "ParseArgs must return non-nil")
	if ps != nil {
		th.RequireFalse(ps.Contains("logLevel"), "property must not exist because it wasn't provided")
	}

	args = append(args, "--logLevel")
	args = append(args, "5")
	ps = ap.ParseArgs(args)
	th.Require(ps != nil, "ParseArgs must return non-nil")
	if ps != nil {
		th.Require(ps.Contains("logLevel"), "property must exist when provided")
		value := ps.GetIntValue("logLevel")
		th.Require(value == 5, "int value must match what was provided")
	}
}

func TestAddOptionalStringArgument(t *testing.T) {
	th := NewTestHelper(t)
	ap := NewArgumentParser()
	ap.AddOptionalStringArgument("--user", "user id for command")

	args := make([]string, 0)
	ps := ap.ParseArgs(args)
	th.Require(ps != nil, "ParseArgs must return non-nil")
	if ps != nil {
		th.RequireFalse(ps.Contains("user"), "property should not exist if not provided")
	}

	args = append(args, "--user")
	args = append(args, "johndoe")
	ps = ap.ParseArgs(args)
	th.Require(ps != nil, "parse_args must return non-nil")
	if ps != nil {
		th.Require(ps.Contains("user"), "provided property must exist")
		userid := ps.GetStringValue("user")
		th.RequireStringEquals("johndoe", userid, "string property values must match")
	}
}

func TestAddRequiredArgument(t *testing.T) {
	th := NewTestHelper(t)
	ap := NewArgumentParser()
	ap.AddRequiredArgument("command", "command to execute")

	args := make([]string, 0)
	ps := ap.ParseArgs(args)
	th.Require(ps != nil, "ParseArgs must return non-nil")
	if ps != nil {
		th.RequireFalse(ps.Contains("command"), "property not provided should not exist")
	}

	args = append(args, "play")
	ps = ap.ParseArgs(args)
	th.Require(ps != nil, "ParseArgs must return non-nil")
	if ps != nil {
		th.Require(ps.Contains("command"), "provided property must exist")
		command := ps.GetStringValue("command")
		th.RequireStringEquals("play", command, "")
	}
}

func TestParseArgs(t *testing.T) {
	th := NewTestHelper(t)
	ap := NewArgumentParser()
	ap.AddOptionalBoolFlag("--debug", "provide debugging support")
	ap.AddOptionalIntArgument("--logLevel", "adjust logging level up or down")
	ap.AddOptionalStringArgument("--user", "user issuing command")
	ap.AddRequiredArgument("command", "command to execute")

	args := make([]string, 0)
	args = append(args, "--logLevel")
	args = append(args, "6")
	args = append(args, "--user")
	args = append(args, "tomjones")
	args = append(args, "--debug")
	args = append(args, "play")

	ps := ap.ParseArgs(args)
	th.Require(ps != nil, "parse_args should return non-nil")
	if ps != nil {
		th.Require(ps.Contains("logLevel"), "logLevel should exist")
		th.Require(ps.Contains("user"), "user should exist")
		th.Require(ps.Contains("debug"), "debug should exist")
		th.Require(ps.Contains("command"), "command should exist")
		logLevel := ps.GetIntValue("logLevel")
		user := ps.GetStringValue("user")
		debug := ps.GetBoolValue("debug")
		command := ps.GetStringValue("command")
		th.Require(logLevel == 6, "")
		th.RequireStringEquals("tomjones", user, "")
		th.Require(debug, "")
		th.RequireStringEquals("play", command, "")
	}
}
