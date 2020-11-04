package ipfix

// decodeState represents the state while decoding IPFIX traffic flow.
type decodeState struct {
	templates        []templateRecord
	optionsTemplates []optionsTemplateRecord
}
