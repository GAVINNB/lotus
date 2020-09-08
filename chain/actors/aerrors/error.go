package aerrors

import (
	"fmt"

	"github.com/filecoin-project/go-state-types/exitcode"
	"golang.org/x/xerrors"
)

func IsFatal(err ActorError) bool {
	return err != nil && err.IsFatal()
}
func RetCode(err ActorError) exitcode.ExitCode {
	if err == nil {
		return 0
	}
	return err.RetCode()
}

type internalActorError interface {
	ActorError
	FormatError(p xerrors.Printer) (next error)
	Unwrap() error
}

type guard struct{}

type ActorError interface {
	error
	IsFatal() bool
	RetCode() exitcode.ExitCode
	// guard is a sentinel to prevent code outside if the aerrors package
	// from implementing the ActorError interface. It ensures actors
	// use Runtime.Abortf and do not panic with an ActorError directly
	// and illegally using a system error code.
	guard() guard
}

type actorError struct {
	fatal   bool
	retCode exitcode.ExitCode

	msg   string
	frame xerrors.Frame
	err   error
}

func (e *actorError) IsFatal() bool {
	return e.fatal
}

func (e *actorError) RetCode() exitcode.ExitCode {
	return e.retCode
}

func (e *actorError) Error() string {
	return fmt.Sprint(e)
}
func (e *actorError) Format(s fmt.State, v rune) { xerrors.FormatError(e, s, v) }
func (e *actorError) FormatError(p xerrors.Printer) (next error) {
	p.Print(e.msg)
	if e.fatal {
		p.Print(" (FATAL)")
	} else {
		p.Printf(" (RetCode=%d)", e.retCode)
	}

	e.frame.Format(p)
	return e.err
}

func (e *actorError) Unwrap() error {
	return e.err
}

func (e *actorError) guard() guard {
	return guard{}
}

var _ internalActorError = (*actorError)(nil)
