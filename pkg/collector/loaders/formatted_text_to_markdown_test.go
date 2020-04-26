package loaders

import (
	"encoding/json"
	"github.com/Arman92/go-tdlib"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_formattedTextToMarkdown(t *testing.T) {
	tests := []struct {
		name string
		args string
		want string
	}{
		{
			name: "textEntityTypeTextUrl simple",
			args: `{
    			"@type": "formattedText",
    			"text": "один два три четыре",
    			"entities": [
        			{
						"@type": "textEntity",
						"offset": 5,
						"length": 3,
						"type": { "@type": "textEntityTypeTextUrl", "url": "https://example.com" }
					}
    			]
			}`,
			want: `один [два](https://example.com) три четыре`,
		},
		{
			name: "textEntityTypeTextUrl multiple",
			args: `{
				"text": "один два три четыре пять",
				"entities": [
					{
						"@type": "textEntity",
						"offset": 5,
						"length": 3,
						"type": { "@type": "textEntityTypeTextUrl", "url": "https://example.com/"}
					},
					{
						"@type": "textEntity",
						"offset": 13,
						"length": 6,
						"type": { "@type": "textEntityTypeTextUrl",  "url": "https://w3c.org/"}
					}
				]
			}`,
			want: `один [два](https://example.com/) три [четыре](https://w3c.org/) пять`,
		},
		{
			name: "textEntityTypeBold",
			args: `{
				"text": "Bold text",
				"entities": [
				  {
					"@type": "textEntity",
					"offset": 5,
					"length": 4,
					"type": { "@type": "textEntityTypeBold" }
				  }
				]
			}`,
			want: `Bold **text**`,
		},
		{
			name: "textEntityTypeMention",
			args: `{
				"text": "mention @korchasa",
				"entities": [
					{
						"@type": "textEntity",
						"offset": 8,
						"length": 9,
						"type": { "@type": "textEntityTypeMention" }
				    }
				]
			}`,
			want: "mention [@korchasa](https://t.me/korchasa)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ft := tdlib.FormattedText{}
			err := json.Unmarshal([]byte(tt.args), &ft)
			assert.NoError(t, err)
			got, errs := FormattedTextToMarkdown(&ft)
			for _, e := range errs {
				assert.NoError(t, e)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_Utf8Substr(t *testing.T) {
	type args struct {
		s       string
		offset  int32
		length  int32
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "from begin",
			args: args{"один два три", 0, 3},
			want: "оди",
		},
		{
			name: "in a middle",
			args: args{"один два три", 5, 3},
			want: "два",
		},
		{
			name: "to the end",
			args: args{"один два три", 9, -1},
			want: "три",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, utf8Substr(tt.args.s, tt.args.offset, tt.args.length))
		})
	}
}