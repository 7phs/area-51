package config

import (
	"fmt"
)

var (
	_ error = NotExists{}
	_ error = PermissionDenied{}
	_ error = IsDir{}
	_ error = Unexpected{}
	_ error = EmptyParam{}
)

type NotExists struct {
	path string
}

func ErrNotExists(path string) NotExists {
	return NotExists{path: path}
}

func (e NotExists) Error() string {
	return fmt.Sprintf("not found: '%s'", e.path)
}

type PermissionDenied struct {
	path string
}

func ErrPermissionDenied(path string) PermissionDenied {
	return PermissionDenied{path: path}
}

func (e PermissionDenied) Error() string {
	return fmt.Sprintf("permission denied: '%s'", e.path)
}

type IsDir struct {
	path string
}

func ErrIsDir(path string) IsDir {
	return IsDir{path: path}
}

func (e IsDir) Error() string {
	return fmt.Sprintf("it is a dir, but a file is expected: '%s'", e.path)
}

type EqualPath struct {
	onePath     string
	anotherPath string
}

func ErrEqualPath(onePath, anotherPath string) EqualPath {
	return EqualPath{
		onePath:     onePath,
		anotherPath: anotherPath,
	}
}

func (e EqualPath) Error() string {
	return fmt.Sprintf("paths are equal: '%s' and '%s'", e.onePath, e.anotherPath)
}

type Unexpected struct {
	err  error
	path string
}

func ErrUnexpected(path string, err error) Unexpected {
	return Unexpected{
		path: path,
		err:  err,
	}
}

func (e Unexpected) Error() string {
	return fmt.Sprintf("unexpected error: '%s' [%v]", e.path, e.err)
}

func (e Unexpected) Unwrap() error {
	return e.err
}

type EmptyParam struct {
	paramName string
}

func ErrEmptyParam(paramName string) EmptyParam {
	return EmptyParam{paramName: paramName}
}

func (e EmptyParam) Error() string {
	return fmt.Sprintf("empty param: '%s'", e.paramName)
}
