package signature

import "github.com/carlcui/expressive/typing"

func addBuiltInSignatures() {
	keyToSignatures[ADD] = []*Signature{
		CreateSignature(typing.INT, typing.INT, typing.INT),
		CreateSignature(typing.FLOAT, typing.FLOAT, typing.FLOAT),
		CreateSignature(typing.STRING, typing.STRING, typing.STRING),
	}
	keyToSignatures[SUBTRACT] = []*Signature{
		CreateSignature(typing.INT, typing.INT, typing.INT),
		CreateSignature(typing.FLOAT, typing.FLOAT, typing.FLOAT),
	}
	keyToSignatures[MULTIPLY] = []*Signature{
		CreateSignature(typing.INT, typing.INT, typing.INT),
		CreateSignature(typing.FLOAT, typing.FLOAT, typing.FLOAT),
	}
	keyToSignatures[DIVIDE] = []*Signature{
		CreateSignature(typing.INT, typing.INT, typing.INT),
		CreateSignature(typing.FLOAT, typing.FLOAT, typing.FLOAT),
	}
	keyToSignatures[MODULO] = []*Signature{
		CreateSignature(typing.INT, typing.INT, typing.INT),
	}
	keyToSignatures[EXPONENTIATE] = []*Signature{
		CreateSignature(typing.INT, typing.INT, typing.INT),
		CreateSignature(typing.FLOAT, typing.FLOAT, typing.FLOAT),
	}
	keyToSignatures[LOGIC_AND] = []*Signature{
		CreateSignature(typing.BOOL, typing.BOOL, typing.BOOL),
	}
	keyToSignatures[LOGIC_OR] = []*Signature{
		CreateSignature(typing.BOOL, typing.BOOL, typing.BOOL),
	}
	keyToSignatures[LOGIC_NOT] = []*Signature{
		CreateSignature(typing.BOOL, typing.BOOL),
	}
	keyToSignatures[IF_ELSE] = []*Signature{
		CreateSignature(typing.BOOL, typing.BOOL, typing.BOOL, typing.BOOL),
		CreateSignature(typing.INT, typing.BOOL, typing.INT, typing.INT),
		CreateSignature(typing.FLOAT, typing.BOOL, typing.FLOAT, typing.FLOAT),
		CreateSignature(typing.CHAR, typing.BOOL, typing.CHAR, typing.CHAR),
		CreateSignature(typing.STRING, typing.BOOL, typing.STRING, typing.STRING),
	}
	keyToSignatures[GREATER] = []*Signature{
		CreateSignature(typing.BOOL, typing.INT, typing.INT),
		CreateSignature(typing.BOOL, typing.FLOAT, typing.FLOAT),
	}
	keyToSignatures[GREATER_OR_EQUAL] = []*Signature{
		CreateSignature(typing.BOOL, typing.INT, typing.INT),
		CreateSignature(typing.BOOL, typing.FLOAT, typing.FLOAT),
	}
	keyToSignatures[LESS] = []*Signature{
		CreateSignature(typing.BOOL, typing.INT, typing.INT),
		CreateSignature(typing.BOOL, typing.FLOAT, typing.FLOAT),
	}
	keyToSignatures[LESS_OR_EQUAL] = []*Signature{
		CreateSignature(typing.BOOL, typing.INT, typing.INT),
		CreateSignature(typing.BOOL, typing.FLOAT, typing.FLOAT),
	}
	keyToSignatures[SHALLOW_EQUAL] = []*Signature{
		CreateSignature(typing.BOOL, typing.INT, typing.INT),
		CreateSignature(typing.BOOL, typing.FLOAT, typing.FLOAT),
		CreateSignature(typing.BOOL, typing.CHAR, typing.CHAR),
		CreateSignature(typing.BOOL, typing.BOOL, typing.BOOL),
		CreateSignature(typing.BOOL, typing.STRING, typing.STRING),
	}
	keyToSignatures[SHALLOW_NOT_EQUAL] = []*Signature{
		CreateSignature(typing.BOOL, typing.INT, typing.INT),
		CreateSignature(typing.BOOL, typing.FLOAT, typing.FLOAT),
		CreateSignature(typing.BOOL, typing.CHAR, typing.CHAR),
		CreateSignature(typing.BOOL, typing.BOOL, typing.BOOL),
		CreateSignature(typing.BOOL, typing.STRING, typing.STRING),
	}
	keyToSignatures[DEEP_EQUAL] = []*Signature{
		CreateSignature(typing.BOOL, typing.INT, typing.INT),
		CreateSignature(typing.BOOL, typing.FLOAT, typing.FLOAT),
		CreateSignature(typing.BOOL, typing.CHAR, typing.CHAR),
		CreateSignature(typing.BOOL, typing.BOOL, typing.BOOL),
		CreateSignature(typing.BOOL, typing.STRING, typing.STRING),
	}
	keyToSignatures[DEEP_NOT_EQUAL] = []*Signature{
		CreateSignature(typing.BOOL, typing.INT, typing.INT),
		CreateSignature(typing.BOOL, typing.FLOAT, typing.FLOAT),
		CreateSignature(typing.BOOL, typing.CHAR, typing.CHAR),
		CreateSignature(typing.BOOL, typing.BOOL, typing.BOOL),
		CreateSignature(typing.BOOL, typing.STRING, typing.STRING),
	}
}
