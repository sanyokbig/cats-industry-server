// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package schema

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

func easyjsonCef4e921DecodeCatsIndustryServerSchema(in *jlexer.Lexer, out *Skill) {
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
		case "skill_id":
			out.ID = uint(in.Uint())
		case "skillpoints_in_skill":
			out.Skillpoints = uint(in.Uint())
		case "trained_skill_level":
			out.TrainedLevel = uint(in.Uint())
		case "active_skill_level":
			out.ActiveLevel = uint(in.Uint())
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
func easyjsonCef4e921EncodeCatsIndustryServerSchema(out *jwriter.Writer, in Skill) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"skill_id\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Uint(uint(in.ID))
	}
	{
		const prefix string = ",\"skillpoints_in_skill\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Uint(uint(in.Skillpoints))
	}
	{
		const prefix string = ",\"trained_skill_level\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Uint(uint(in.TrainedLevel))
	}
	{
		const prefix string = ",\"active_skill_level\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Uint(uint(in.ActiveLevel))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Skill) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonCef4e921EncodeCatsIndustryServerSchema(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Skill) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonCef4e921EncodeCatsIndustryServerSchema(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Skill) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonCef4e921DecodeCatsIndustryServerSchema(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Skill) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonCef4e921DecodeCatsIndustryServerSchema(l, v)
}
func easyjsonCef4e921DecodeCatsIndustryServerSchema1(in *jlexer.Lexer, out *Payload) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		in.Skip()
	} else {
		in.Delim('{')
		if !in.IsDelim('}') {
			*out = make(Payload)
		} else {
			*out = nil
		}
		for !in.IsDelim('}') {
			key := string(in.String())
			in.WantColon()
			var v1 interface{}
			if m, ok := v1.(easyjson.Unmarshaler); ok {
				m.UnmarshalEasyJSON(in)
			} else if m, ok := v1.(json.Unmarshaler); ok {
				_ = m.UnmarshalJSON(in.Raw())
			} else {
				v1 = in.Interface()
			}
			(*out)[key] = v1
			in.WantComma()
		}
		in.Delim('}')
	}
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonCef4e921EncodeCatsIndustryServerSchema1(out *jwriter.Writer, in Payload) {
	if in == nil && (out.Flags&jwriter.NilMapAsEmpty) == 0 {
		out.RawString(`null`)
	} else {
		out.RawByte('{')
		v2First := true
		for v2Name, v2Value := range in {
			if v2First {
				v2First = false
			} else {
				out.RawByte(',')
			}
			out.String(string(v2Name))
			out.RawByte(':')
			if m, ok := v2Value.(easyjson.Marshaler); ok {
				m.MarshalEasyJSON(out)
			} else if m, ok := v2Value.(json.Marshaler); ok {
				out.Raw(m.MarshalJSON())
			} else {
				out.Raw(json.Marshal(v2Value))
			}
		}
		out.RawByte('}')
	}
}

// MarshalJSON supports json.Marshaler interface
func (v Payload) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonCef4e921EncodeCatsIndustryServerSchema1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Payload) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonCef4e921EncodeCatsIndustryServerSchema1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Payload) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonCef4e921DecodeCatsIndustryServerSchema1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Payload) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonCef4e921DecodeCatsIndustryServerSchema1(l, v)
}
func easyjsonCef4e921DecodeCatsIndustryServerSchema2(in *jlexer.Lexer, out *Message) {
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
		case "type":
			out.Type = string(in.String())
		case "payload":
			(out.Payload).UnmarshalEasyJSON(in)
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
func easyjsonCef4e921EncodeCatsIndustryServerSchema2(out *jwriter.Writer, in Message) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"type\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Type))
	}
	{
		const prefix string = ",\"payload\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		(in.Payload).MarshalEasyJSON(out)
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Message) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonCef4e921EncodeCatsIndustryServerSchema2(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Message) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonCef4e921EncodeCatsIndustryServerSchema2(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Message) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonCef4e921DecodeCatsIndustryServerSchema2(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Message) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonCef4e921DecodeCatsIndustryServerSchema2(l, v)
}
func easyjsonCef4e921DecodeCatsIndustryServerSchema3(in *jlexer.Lexer, out *CharactersList) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		in.Skip()
		*out = nil
	} else {
		in.Delim('[')
		if *out == nil {
			if !in.IsDelim(']') {
				*out = make(CharactersList, 0, 1)
			} else {
				*out = CharactersList{}
			}
		} else {
			*out = (*out)[:0]
		}
		for !in.IsDelim(']') {
			var v3 Character
			(v3).UnmarshalEasyJSON(in)
			*out = append(*out, v3)
			in.WantComma()
		}
		in.Delim(']')
	}
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonCef4e921EncodeCatsIndustryServerSchema3(out *jwriter.Writer, in CharactersList) {
	if in == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
		out.RawString("null")
	} else {
		out.RawByte('[')
		for v4, v5 := range in {
			if v4 > 0 {
				out.RawByte(',')
			}
			(v5).MarshalEasyJSON(out)
		}
		out.RawByte(']')
	}
}

// MarshalJSON supports json.Marshaler interface
func (v CharactersList) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonCef4e921EncodeCatsIndustryServerSchema3(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v CharactersList) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonCef4e921EncodeCatsIndustryServerSchema3(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *CharactersList) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonCef4e921DecodeCatsIndustryServerSchema3(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *CharactersList) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonCef4e921DecodeCatsIndustryServerSchema3(l, v)
}
func easyjsonCef4e921DecodeCatsIndustryServerSchema4(in *jlexer.Lexer, out *Character) {
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
		case "id":
			out.ID = uint(in.Uint())
		case "name":
			out.Name = string(in.String())
		case "is_main":
			out.IsMain = bool(in.Bool())
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
func easyjsonCef4e921EncodeCatsIndustryServerSchema4(out *jwriter.Writer, in Character) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Uint(uint(in.ID))
	}
	{
		const prefix string = ",\"name\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Name))
	}
	{
		const prefix string = ",\"is_main\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Bool(bool(in.IsMain))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Character) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonCef4e921EncodeCatsIndustryServerSchema4(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Character) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonCef4e921EncodeCatsIndustryServerSchema4(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Character) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonCef4e921DecodeCatsIndustryServerSchema4(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Character) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonCef4e921DecodeCatsIndustryServerSchema4(l, v)
}
