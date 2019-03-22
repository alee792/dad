package dad_test

import (
	"context"
	"io"
	"reflect"
	"strings"
	"testing"

	"github.com/alee792/dad/pkg/dad"
	"github.com/google/go-cmp/cmp"
)

// fields define common testing components.
type fields struct {
	chain *dad.Chain
	cfg   dad.Config
}

func TestNewChain(t *testing.T) {
	c := dad.NewChain(nil, dad.Config{})
	if c.Config.Order < 1 {
		t.Fatalf("expected non-zero, positive default")
	}
}

func TestChain_Read(t *testing.T) {
	type args struct {
		ctx context.Context
		s   io.Reader
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   map[string][]string
	}{
		{
			name:   "duplicate 1-grams",
			fields: defaultFields(t),
			args: args{
				ctx: context.TODO(),
				s:   strings.NewReader("one two one three\none two one three"),
			},
			want: map[string][]string{
				"one":   []string{"two", "three", "two", "three"},
				"two":   []string{"one", "one"},
				"three": []string(nil),
			},
		},
		{
			name:   "duplicate 1-gram, tricky finale",
			fields: defaultFields(t),
			args: args{
				ctx: context.TODO(),
				s:   strings.NewReader("one two one three one\none two one three one"),
			},
			want: map[string][]string{
				"one":   []string{"two", "three", "two", "three"},
				"two":   []string{"one", "one"},
				"three": []string{"one", "one"},
			},
		},
		{
			name: "duplicate 2-gram",
			fields: func() fields {
				f := defaultFields(t)
				f.chain.Config.Order = 2
				return f
			}(),
			args: args{
				ctx: context.TODO(),
				s:   strings.NewReader("one two one three\none two one three"),
			},
			want: map[string][]string{
				"one two":   []string{"one three", "one three"},
				"two one":   []string{"three", "three"},
				"one three": []string(nil),
			},
		},
		{
			name:   "unique 1-gram, tricky finale",
			fields: defaultFields(t),
			args: args{
				ctx: context.TODO(),
				s:   strings.NewReader("one two one three one\none four one three one"),
			},
			want: map[string][]string{
				"one":   []string{"two", "three", "four", "three"},
				"two":   []string{"one"},
				"three": []string{"one", "one"},
				"four":  []string{"one"},
			},
		},
		{
			name: "unique 2-gram",
			fields: func() fields {
				f := defaultFields(t)
				f.chain.Config.Order = 2
				return f
			}(),
			args: args{
				ctx: context.TODO(),
				s:   strings.NewReader("one two one three\none three one four"),
			},
			want: map[string][]string{
				"one two":   []string{"one three"},
				"two one":   []string{"three"},
				"one three": []string{"one four"},
				"three one": []string{"four"},
				"one four":  []string(nil),
			},
		},
		{
			name:   "ctx cancel",
			fields: defaultFields(t),
			args: args{
				ctx: func() context.Context {
					ctx := context.Background()
					ctx, cancel := context.WithCancel(ctx)
					cancel()
					return ctx
				}(),
				s: strings.NewReader("one two one three\none two one three"),
			},
			want: make(map[string][]string),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields.chain.Read(tt.args.ctx, tt.args.s)
		})
		if !reflect.DeepEqual(tt.want, tt.fields.chain.Grams()) {
			t.Errorf(cmp.Diff(tt.want, tt.fields.chain.Grams()))
		}
	}
}

func TestChain_ReadSentence(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   map[string][]string
	}{
		{
			name:   "1-gram",
			fields: defaultFields(t),
			args: args{
				s: "one two one three",
			},
			want: map[string][]string{
				"one":   []string{"two", "three"},
				"two":   []string{"one"},
				"three": []string(nil),
			},
		},
		{
			name:   "1-gram, tricky finale",
			fields: defaultFields(t),
			args: args{
				s: "one two one three one",
			},
			want: map[string][]string{
				"one":   []string{"two", "three"},
				"two":   []string{"one"},
				"three": []string{"one"},
			},
		},
		{
			name: "2-gram",
			fields: func() fields {
				f := defaultFields(t)
				f.chain.Config.Order = 2
				return f
			}(),
			args: args{
				s: "one two one three",
			},
			want: map[string][]string{
				"one two":   []string{"one three"},
				"two one":   []string{"three"},
				"one three": []string(nil),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields.chain.ReadSentence(tt.args.s)
		})
		if !reflect.DeepEqual(tt.want, tt.fields.chain.Grams()) {
			t.Errorf(cmp.Diff(tt.want, tt.fields.chain.Grams()))
		}
	}
}

func TestChain_Generate(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		corpus io.Reader
	}{
		{
			name:   "1-gram",
			fields: defaultFields(t),
			args: args{
				ctx: context.TODO(),
			},
			corpus: defaultCorpus(t),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields.chain.Read(tt.args.ctx, tt.corpus)
			var joke string
			for joke == "" {
				joke = tt.fields.chain.Generate(tt.args.ctx)
			}
			t.Log(joke)
		})
	}
}

func TestChain_GenerateWithMax(t *testing.T) {
	type args struct {
		ctx context.Context
		max int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		corpus io.Reader
	}{
		{
			name:   "tinygram",
			fields: defaultFields(t),
			args: args{
				ctx: context.TODO(),
				max: 1,
			},
			corpus: defaultCorpus(t),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields.chain.Read(tt.args.ctx, tt.corpus)
			var joke string
			for joke == "" {
				joke = tt.fields.chain.GenerateWithMax(tt.args.ctx, tt.args.max)
			}
			gotWords := len(strings.Fields(joke))
			if gotWords > tt.args.max {
				t.Errorf("expected joke with <= %d words, got %d", tt.args.max, gotWords)
			}
			t.Log(joke)
		})
	}
}

func defaultFields(t *testing.T) fields {
	f := fields{
		chain: defaultChain(t),
		cfg:   defaultConfig(t),
	}
	return f
}

func defaultChain(t *testing.T) *dad.Chain {
	return dad.NewChain(nil, defaultConfig(t))
}

func defaultConfig(t *testing.T) dad.Config {
	return dad.Config{
		Order: 1,
	}
}

func defaultCorpus(t *testing.T) io.Reader {
	return strings.NewReader(`
	Why are fish so smart? Because they live in schools!
	Why does a chicken coop only have two doors? Because if it had four doors it would be a chicken sedan.
	What did Romans use to cut pizza before the rolling cutter was invented? Lil Caesars
	Did you hear about the chameleon who couldn't change color? They had a reptile dysfunction.
	Did you know that protons have mass? I didn't even know they were catholic.
	What do you call a group of disorganized cats? A cat-tastrophe.
?Dfjk adj 01978&$&7 818 ^&((#()))
	`)
}
