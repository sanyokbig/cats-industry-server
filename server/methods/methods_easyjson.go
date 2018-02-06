// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package methods

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjson9f911648DecodeCatsIndustryServerServerMethods(in *jlexer.Lexer, out *restoreSessionPayload) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "sid":
			out.SID = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson9f911648EncodeCatsIndustryServerServerMethods(out *jwriter.Writer, in restoreSessionPayload) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"sid\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.SID))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v restoreSessionPayload) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson9f911648EncodeCatsIndustryServerServerMethods(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v restoreSessionPayload) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson9f911648EncodeCatsIndustryServerServerMethods(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *restoreSessionPayload) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson9f911648DecodeCatsIndustryServerServerMethods(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *restoreSessionPayload) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson9f911648DecodeCatsIndustryServerServerMethods(l, v)
}
func easyjson9f911648DecodeCatsIndustryServerServerMethods1(in *jlexer.Lexer, out *loginRequestPayload) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "scope_set":
			out.ScopeSet = string(in.String())
		case "sid":
			out.SID = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson9f911648EncodeCatsIndustryServerServerMethods1(out *jwriter.Writer, in loginRequestPayload) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"scope_set\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.ScopeSet))
	}
	{
		const prefix string = ",\"sid\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.SID))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v loginRequestPayload) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson9f911648EncodeCatsIndustryServerServerMethods1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v loginRequestPayload) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson9f911648EncodeCatsIndustryServerServerMethods1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *loginRequestPayload) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson9f911648DecodeCatsIndustryServerServerMethods1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *loginRequestPayload) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson9f911648DecodeCatsIndustryServerServerMethods1(l, v)
}
