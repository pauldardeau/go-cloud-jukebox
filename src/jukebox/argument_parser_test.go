package jukebox

import "testing"

func TestNewArgumentParser(t *testing.T) {
	ap := NewArgumentParser()
	if ap == nil {
		t.Log("NewArgumentParser should not return nil")
		t.Fail()
	}
}

func Test_addOption(t *testing.T) {
}

func TestAddOptionalBoolFlag(t *testing.T) {
}

func TestAddOptionalIntArgument(t *testing.T) {
}

func TestAddOptionalStringArgument(t *testing.T) {
}

func TestAddRequiredArgument(t *testing.T) {
}

func TestParseArgs(t *testing.T) {
}
