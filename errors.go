package gitconfig

import "errors"

// ErrInvalidKeyChar indicates that there was an invalid key character
var ErrInvalidKeyChar = errors.New("invalid key character")

// ErrMissingStartQuote indicates that there was a missing start quote
var ErrMissingStartQuote = errors.New("missing start quote")

// ErrNotBoolValue indicates fail to convert config variable to bool
var ErrNotBoolValue = errors.New("not a bool value")

// ErrNotExist indicates file or dir not exist
var ErrNotExist = errors.New("config file or dir not exist")

// ErrNotInGitDir indicates not in a git dir
var ErrNotInGitDir = errors.New("not in a git dir")
