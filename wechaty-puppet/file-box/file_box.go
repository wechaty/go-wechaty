// file_box
// Deprecated: use filebox package

package file_box

import (
	"github.com/wechaty/go-wechaty/wechaty-puppet/filebox"
)

// FileBox ...
// Deprecated: use filebox.FileBox
type FileBox = filebox.FileBox

var (
	// FromJSON ...
	// Deprecated: use filebox.FileBox
	FromJSON = filebox.FromJSON

	// FromBase64 ...
	// Deprecated: use filebox.FromBase64
	FromBase64 = filebox.FromBase64

	// FromUrl ...
	// Deprecated: use filebox.FromUrl
	FromUrl = filebox.FromUrl

	// FromFile ...
	// Deprecated: use filebox.FromFile
	FromFile = filebox.FromFile

	// FromQRCode ...
	// Deprecated: use filebox.FromQRCode
	FromQRCode = filebox.FromQRCode

	// FromStream ...
	// Deprecated: use filebox.FromStream
	FromStream = filebox.FromStream
)
