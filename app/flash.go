package app

type FlashType int

const (
	FlashTypeSuccess FlashType = iota
	FlashTypeInfo              = iota
	FlashTypeWarning           = iota
	FlashTypeError             = iota
)

// This is a Flash. It defines a message to be displayed at the website's next convenience.
type Flash struct {
	Header  string
	Message string
	Type    FlashType
}

// Returns the Bootstrap CSS class for this Flash's Type. Note that the class is preceded with a space.
func (flash *Flash) CSSClass() string {
	switch {
	case flash.Type == FlashTypeSuccess:
		return " alert-success"
	case flash.Type == FlashTypeInfo:
		return " alert-info"
	case flash.Type == FlashTypeWarning:
		return ""
	case flash.Type == FlashTypeError:
		return " alert-error"
	}
	return ""
}
