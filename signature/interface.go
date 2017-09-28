package signature

import "github.com/carlcui/expressive/typing"
import "fmt"

type Signature struct {
	Params []typing.Typing
	Result typing.Typing
}

func (signature *Signature) Accepts(params ...typing.Typing) bool {
	if len(signature.Params) != len(params) {
		return false
	}

	for i, signatureParam := range signature.Params {
		param := params[i]

		if !signatureParam.Equals(param) {
			return false
		}
	}

	return true
}

// Mapping is from key to signatures. The key can only be an operator or a string.
type Mapping map[interface{}][]*Signature

var keyToSignatures Mapping

func HasSignature(key interface{}, params ...typing.Typing) bool {
	resultType := ResultTyping(key, params...)

	if resultType == typing.ERROR_TYPE {
		return false
	}

	return true
}

func ResultTyping(key interface{}, params ...typing.Typing) typing.Typing {
	checkValidKey(key)

	signatures, ok := keyToSignatures[key]

	if !ok {
		return typing.ERROR_TYPE
	}

	for _, signature := range signatures {
		if signature.Accepts(params...) {
			return signature.Result
		}
	}

	return typing.ERROR_TYPE
}

func checkValidKey(key interface{}) {
	switch key.(type) {
	case Operator:
		return
	case string:
		return
	default:
	}

	err := fmt.Errorf("passing invalid key: %v", key)
	panic(err.Error())
}

func CreateSignature(result typing.Typing, params ...typing.Typing) *Signature {
	return &Signature{
		Params: params,
		Result: result,
	}
}

func init() {
	keyToSignatures = make(Mapping)

	addBuiltInSignatures()
}
