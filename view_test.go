package golivewire

import "testing"

func TestHtmlView_RenderSafe(t *testing.T) {
	tests := []struct {
		name   string
		raw    string
		expect string
	}{
		{
			name:   "return same",
			raw:    "<p>test</p>",
			expect: "<p>test</p>",
		},
		{
			name:   "multiple tag",
			raw:    "<p>1</p><p>2</p>",
			expect: "<p>1</p><p>2</p>",
		},
		{
			name:   "no tag",
			raw:    "word",
			expect: "",
		},
		{
			name:   "no tag 2",
			raw:    "word<p>test</p>",
			expect: "<p>test</p>",
		},
		{
			name:   "no tag 3",
			raw:    "<p>test</p>word",
			expect: "<p>test</p>",
		},
		{
			name:   "comment first",
			raw:    "<!-- A comment here --><div>Element</div>",
			expect: "<div>Element</div>",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			view, err := newHTMLView(test.raw)
			if err != nil {
				t.Error(err)
			}
			actual, err := view.RenderSafe()
			if err != nil {
				t.Error(err)
			}
			if string(actual) != test.expect {
				t.Errorf("html is not matched. expected:\n%s\n\nactual:\n%s", test.expect, actual)
			}
		})
	}
}

func TestHtmlView_AddWireTag(t *testing.T) {
	tests := []struct {
		name        string
		initialHtml string
		key         string
		value       string
		expected    string
	}{
		{
			name:        "id",
			initialHtml: "<p>test</p>",
			key:         "id",
			value:       "1a2b",
			expected:    "<p wire:id=\"1a2b\">test</p>",
		},
		{
			name:        "json escape",
			initialHtml: "<p>test</p>",
			key:         "initial-data",
			value:       "{\"id\":\"1a2b\"}",
			expected:    "<p wire:initial-data=\"{&#34;id&#34;:&#34;1a2b&#34;}\">test</p>",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			view, err := newHTMLView(test.initialHtml)
			if err != nil {
				t.Error(err)
			}
			view.AddWireTag(test.key, test.value)
			actual, err := view.RenderSafe()
			if err != nil {
				t.Error(err)
			}
			if string(actual) != test.expected {
				t.Errorf("html is not matched. expected:\n%s\n\nactual:\n%s", test.expected, actual)
			}
		})
	}
}
