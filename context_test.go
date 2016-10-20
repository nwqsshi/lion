package lion

import (
	"testing"

	"context"
)

// mss is an alias for map[string]string.
type mss map[string]string

func (p *ctx) toMap() mss {
	m := mss{}
	for i := range p.keys {
		m[p.keys[i]] = p.values[i]
	}
	return m
}

func TestContextAddParam(t *testing.T) {
	c := NewContext()
	c.parent = context.WithValue(context.TODO(), "parentKey", "parentVal")
	c.AddParam("key", "val")
	if len(c.keys) != 1 {
		t.Errorf("Length of keys should be 1 but got %s", red(len(c.keys)))
	}

	val := c.Value("key")
	if val == nil {
		t.Errorf("Context Value() should not be nil")
	}
	str, ok := val.(string)
	if !ok {
		t.Error("Context value is not a string")
	}

	if str != "val" {
		t.Error("Wrong value for Context.Value()")
	}

	parentVal := c.Value("parentKey")
	if parentVal == nil {
		t.Errorf("Context parent Value() should not be nil")
	}
	str, ok = parentVal.(string)
	if !ok {
		t.Error("Context parent value is not a string")
	}

	if str != "parentVal" {
		t.Error("Wrong value for parent Context.Value()")
	}
}

func TestContextC(t *testing.T) {
	c := C(context.TODO())
	ctx := c.(*ctx)
	if ctx.parent != context.TODO() {
		t.Error("Context C: Parent should be context.TODO()")
	}
}

func TestGetParamInNestedContext(t *testing.T) {
	c := NewContextWithParent(context.Background())
	c.AddParam("id", "myid")

	base := context.WithValue(context.Background(), ctxKey, c)
	nc := context.WithValue(base, "db", "mydb")
	nc = context.WithValue(nc, "t", "t")

	pc := nc.Value(ctxKey)
	if pc == nil {
		t.Errorf("Should not be nil")
	}

	if _, ok := pc.(*ctx); !ok {
		t.Errorf("Should be a *ctx")
	}

	id := Param(nc, "id")
	if id != "myid" {
		t.Errorf("id should be equal to 'myid' but got '%s'", id)
	}

	nonexistant := Param(nc, "nonexistant")
	if nonexistant != "" {
		t.Errorf("Should be empty but got '%s'", nonexistant)
	}

	val := nc.Value("db")
	v, ok := val.(string)
	if !ok {
		t.Errorf("Should be a string")
	}

	if v != "mydb" {
		t.Errorf("Value should be equal to 'mydb' but got '%s'", v)
	}
}
