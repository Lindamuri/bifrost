package local

import (
	"encoding/json"
	"github.com/tremendouscan/bifrost/internal/pkg/code"
	"github.com/tremendouscan/bifrost/pkg/resolv/V3/nginx/configuration/context"
	"github.com/tremendouscan/bifrost/pkg/resolv/V3/nginx/configuration/context_type"
	"github.com/marmotedu/errors"
	"strings"
)

type Directive struct {
	enabled bool
	Name    string
	Params  string

	fatherContext context.Context
}

func (d *Directive) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Enabled     bool                     `json:"enabled,omitempty"`
		ContextType context_type.ContextType `json:"context-type"`
		Value       string                   `json:"value"`
	}{
		Enabled:     d.IsEnabled(),
		ContextType: d.Type(),
		Value:       d.Value(),
	})
}

func (d *Directive) Insert(ctx context.Context, idx int) context.Context {
	return context.ErrContext(errors.WithCode(code.ErrV3InvalidOperation, "directive cannot insert context"))
}

func (d *Directive) Remove(idx int) context.Context {
	return context.ErrContext(errors.WithCode(code.ErrV3InvalidOperation, "directive cannot remove context"))
}

func (d *Directive) Modify(ctx context.Context, idx int) context.Context {
	return context.ErrContext(errors.WithCode(code.ErrV3InvalidOperation, "directive cannot modify context"))
}

func (d *Directive) Father() context.Context {
	return d.fatherContext
}

func (d *Directive) Child(idx int) context.Context {
	return context.ErrContext(errors.WithCode(code.ErrV3InvalidOperation, "directive has no child context"))
}

func (d *Directive) ChildrenPosSet() context.PosSet {
	return context.NewPosSet()
}

func (d *Directive) Clone() context.Context {
	return &Directive{
		enabled:       d.enabled,
		Name:          d.Name,
		Params:        d.Params,
		fatherContext: d.fatherContext,
	}
}

func (d *Directive) SetValue(v string) error {
	kv := strings.SplitN(strings.TrimSpace(v), " ", 2)
	if len(strings.TrimSpace(kv[0])) == 0 {
		return errors.WithCode(code.ErrV3InvalidOperation, "set value for directive failed, cased by: split null value")
	}
	d.Name = strings.TrimSpace(kv[0])
	if len(kv) == 2 {
		d.Params = strings.TrimSpace(kv[1])
	} else {
		d.Params = ""
	}
	return nil
}

func (d *Directive) SetFather(ctx context.Context) error {
	d.fatherContext = ctx
	return nil
}

func (d *Directive) HasChild() bool {
	return false
}

func (d *Directive) Len() int {
	return 0
}

func (d *Directive) Value() string {
	v := strings.TrimSpace(d.Name)
	if params := strings.TrimSpace(d.Params); len(params) > 0 {
		v += " " + params
	}
	return v
}

func (d *Directive) Type() context_type.ContextType {
	return context_type.TypeDirective
}

func (d *Directive) Error() error {
	return nil
}

func (d *Directive) ConfigLines(isDumping bool) ([]string, error) {
	value := d.Value() + ";"
	if !d.IsEnabled() {
		// !!disabled `directive` context will have their line breaks replaced with space characters when dumped!!
		if !isDumping {
			value = "# " + strings.ReplaceAll(value, "\n", "\n# ")
			return strings.Split(value, "\n"), nil
		}
		value = "# " + strings.ReplaceAll(value, "\n", " ")
	}
	return []string{value}, nil
}

func (d *Directive) IsEnabled() bool {
	return d.enabled
}

func (d *Directive) Enable() context.Context {
	d.enabled = true
	return d
}

func (d *Directive) Disable() context.Context {
	d.enabled = false
	return d
}

func registerDirectiveBuilder() error {
	builderMap[context_type.TypeDirective] = func(value string) context.Context {
		kv := strings.SplitN(strings.TrimSpace(value), " ", 2)
		if len(strings.TrimSpace(kv[0])) == 0 {
			return context.ErrContext(errors.New("null value"))
		}
		d := &Directive{
			enabled:       true,
			Name:          strings.TrimSpace(kv[0]),
			fatherContext: context.NullContext(),
		}
		if len(kv) == 2 {
			d.Params = strings.TrimSpace(kv[1])
		}
		return d
	}
	return nil
}

func registerDirectiveParseFunc() error {
	inStackParseFuncMap[context_type.TypeDirective] = func(data []byte, idx *int) context.Context {
		if matchIndexes := RegDirectiveWithoutValue.FindIndex(data[*idx:]); matchIndexes != nil { //nolint:nestif
			subMatch := RegDirectiveWithoutValue.FindSubmatch(data[*idx:])
			*idx += matchIndexes[len(matchIndexes)-1]
			key := string(subMatch[1])
			return NewContext(context_type.TypeDirective, key)
		}

		if matchIndexes := RegDirectiveWithValue.FindIndex(data[*idx:]); matchIndexes != nil { //nolint:nestif
			subMatch := RegDirectiveWithValue.FindSubmatch(data[*idx:])
			*idx += matchIndexes[len(matchIndexes)-1]
			name := string(subMatch[1])
			value := string(subMatch[2])
			if name == string(context_type.TypeInclude) {
				return NewContext(context_type.TypeInclude, value)
			}
			return NewContext(context_type.TypeDirective, name+" "+value)
		}
		return context.NullContext()
	}
	return nil
}
