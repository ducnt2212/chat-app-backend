package middleware

import "net/http"

type Constructor func(http.Handler) http.Handler

type Chain struct {
	constructors []Constructor
}

func NewChain(constructos ...Constructor) Chain {
	return Chain{constructors: append([]Constructor(nil), constructos...)}
}

func (chain Chain) Then(handler http.Handler) http.Handler {
	for i := range chain.constructors {
		handler = chain.constructors[len(chain.constructors)-1-i](handler)
	}

	return handler
}

func (chain Chain) ThenFunc(handlerFunc http.HandlerFunc) http.Handler {
	return chain.Then(handlerFunc)
}

func (chain Chain) Append(constructors ...Constructor) Chain {
	newConstructors := make([]Constructor, 0, len(chain.constructors)+len(constructors))
	newConstructors = append(newConstructors, chain.constructors...)
	newConstructors = append(newConstructors, constructors...)

	return NewChain(newConstructors...)
}

func (chain Chain) Extend(otherChain Chain) Chain {
	return chain.Append(otherChain.constructors...)
}
