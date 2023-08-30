package NitroAdjudicator

import "github.com/statechannels/go-nitro/channel/state"

// ConvertBindingsSignatureToSignature converts the signature type returned from abigien bindings to a state.Signature
func ConvertBindingsSignatureToSignature(s INitroTypesSignature) state.Signature {
	return state.Signature{
		R: s.R[:],
		S: s.S[:],
		V: s.V,
	}
}

// ConvertBindingsSignatureToSignature converts a slice of the signature type returned from abigien bindings to a []state.Signature
func ConvertBindingsSignaturesToSignatures(ss []INitroTypesSignature) []state.Signature {
	sigs := make([]state.Signature, 0, len(ss))
	for _, s := range ss {
		sigs = append(sigs, ConvertBindingsSignatureToSignature(s))
	}
	return sigs
}
