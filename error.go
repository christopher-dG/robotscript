package robotscript

import "github.com/pkg/errors"

// wrongOptType returns an error for an incorrect option value type.
func wrongOptType(command, expected, option string, observed interface{}) error {
	return errors.Errorf(
		"%v: expected %v value to '%v' option (value = %v)",
		command, expected, option, observed,
	)
}

// wrongListEntryType returns an error for an incorrect option value list entry type.
func wrongListEntryType(command, expected, option string, observed interface{}) error {
	return errors.Errorf(
		"%v: expected %v entries in '%v' option (value = %v)",
		command, expected, option, observed,
	)
}

// unrecognizedOpt returns an error for an unknown option.
func unrecognizedOpt(command, option string) error {
	return errors.Errorf("%v: unrecognized option '%v'", command, option)
}

// missingOpt returns an error for a missing required option.
func missingOpt(command, option string) error {
	return errors.Errorf("%v: missing required option '%v'", command, option)
}

// invalidOpt returns an error for a disallowed option value.
func invalidOpt(command, option string, val interface{}) error {
	return errors.Errorf("%v: invalid value for option '%v' (value = %v)", command, option, val)
}

// unrecognizedCmd returns an error for an unknown command.
func unrecognizedCmd(command string) error {
	return errors.Errorf("unrecognized command '%v'", command)
}
