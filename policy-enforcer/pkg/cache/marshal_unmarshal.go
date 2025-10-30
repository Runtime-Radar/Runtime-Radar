package cache

import (
	"bytes"
	"encoding/gob"

	"github.com/gobwas/glob/match"
)

func init() {
	// To be able to marshal glob.Glob interface, we have to register
	// all of its implementations, and this list must be updated if
	// there is any change on github.com/gobwas/glob side:
	gob.Register(match.Any{})
	gob.Register(match.AnyOf{})
	gob.Register(match.BTree{})
	gob.Register(match.Contains{})
	gob.Register(match.EveryOf{})
	gob.Register(match.List{})
	gob.Register(match.Max{})
	gob.Register(match.Min{})
	gob.Register(match.Nothing{})
	gob.Register(match.Prefix{})
	gob.Register(match.PrefixAny{})
	gob.Register(match.PrefixSuffix{})
	gob.Register(match.Range{})
	gob.Register(match.Row{})
	gob.Register(match.Single{})
	gob.Register(match.Suffix{})
	gob.Register(match.SuffixAny{})
	gob.Register(match.Super{})
	gob.Register(match.Text{})
}

// marshal serializes arbitrary data to []byte.
func marshal(v interface{}) ([]byte, error) {
	return marshalGob(v)
}

// unmarshal unserializes data to a given destination.
func unmarshal(data []byte, v interface{}) error {
	return unmarshalGob(data, v)
}

func marshalGob(v interface{}) ([]byte, error) {
	buf := &bytes.Buffer{}
	enc := gob.NewEncoder(buf)
	if err := enc.Encode(v); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func unmarshalGob(data []byte, v interface{}) error {
	r := bytes.NewReader(data)
	dec := gob.NewDecoder(r)

	return dec.Decode(v)
}
