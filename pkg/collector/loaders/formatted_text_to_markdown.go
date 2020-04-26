package loaders

import (
	"encoding/json"
	"fmt"
	"github.com/Arman92/go-tdlib"
)

func FormattedTextToMarkdown(ft *tdlib.FormattedText) (string, []error) {
	t := ft.Text
	var lenCorr int32 = 0 // коррекция для учета изменений длины от предыдущих замен
	errs := make([]error, 0, 1)
	for _, e := range ft.Entities {
		switch e.Type.GetTextEntityTypeEnum() {
		case tdlib.TextEntityTypeBoldType:
			head := utf8Substr(t, 0, e.Offset + lenCorr)
			text := utf8Substr(t, e.Offset + lenCorr, e.Length)
			tail := utf8Substr(t, e.Offset + e.Length + lenCorr, -1)
			t = fmt.Sprintf("%s**%s**%s", head, text, tail)
			lenCorr += int32(len([]rune("****")))
		case tdlib.TextEntityTypeItalicType:
			head := utf8Substr(t, 0, e.Offset + lenCorr)
			text := utf8Substr(t, e.Offset + lenCorr, e.Length)
			tail := utf8Substr(t, e.Offset + e.Length + lenCorr, -1)
			t = fmt.Sprintf("%s*%s*%s", head, text, tail)
			lenCorr += int32(len([]rune("**")))
		case tdlib.TextEntityTypeUnderlineType:
			head := utf8Substr(t, 0, e.Offset + lenCorr)
			text := utf8Substr(t, e.Offset + lenCorr, e.Length)
			tail := utf8Substr(t, e.Offset + e.Length + lenCorr, -1)
			t = fmt.Sprintf("%s~%s~%s", head, text, tail)
			lenCorr += int32(len([]rune("~~")))
		case tdlib.TextEntityTypeStrikethroughType:
			head := utf8Substr(t, 0, e.Offset + lenCorr)
			text := utf8Substr(t, e.Offset + lenCorr, e.Length)
			tail := utf8Substr(t, e.Offset + e.Length + lenCorr, -1)
			t = fmt.Sprintf("%s~~%s~~%s", head, text, tail)
			lenCorr += int32(len([]rune("~~~~")))
		case tdlib.TextEntityTypeCodeType:
			head := utf8Substr(t, 0, e.Offset + lenCorr)
			text := utf8Substr(t, e.Offset + lenCorr, e.Length)
			tail := utf8Substr(t, e.Offset + e.Length + lenCorr, -1)
			t = fmt.Sprintf("%s`%s`%s", head, text, tail)
			lenCorr += int32(len([]rune("``")))
		case tdlib.TextEntityTypeEmailAddressType:
			head := utf8Substr(t, 0, e.Offset + lenCorr)
			address := utf8Substr(t, e.Offset + lenCorr, e.Length)
			tail := utf8Substr(t, e.Offset + e.Length + lenCorr, -1)
			t = fmt.Sprintf("%s[%s](mailto:%s)%s", head, address, address, tail)
			lenCorr += int32(len(address) + len([]rune("[](mailto:)")))
		case tdlib.TextEntityTypeMentionType:
			head := utf8Substr(t, 0, e.Offset + lenCorr)
			nick := utf8Substr(t, e.Offset + lenCorr + 1, e.Length)
			tail := utf8Substr(t, e.Offset + e.Length + lenCorr, -1)
			t = fmt.Sprintf("%s[@%s](https://t.me/%s)%s", head, nick, nick, tail)
			lenCorr += int32(len(nick) + len([]rune("[@](https://t.me/)")))
		case tdlib.TextEntityTypePreType:
			head := utf8Substr(t, 0, e.Offset + lenCorr)
			text := utf8Substr(t, e.Offset + lenCorr, e.Length)
			tail := utf8Substr(t, e.Offset + e.Length + lenCorr, -1)
			t = fmt.Sprintf("%s\n```\n%s\n```\n%s", head, text, tail)
			lenCorr += int32(len([]rune("\n```\n\n```\n")))
		case tdlib.TextEntityTypeTextUrlType: // foo -> [foo](example.com)
			tt := e.Type.(*tdlib.TextEntityTypeTextUrl)
			head := utf8Substr(t, 0, e.Offset + lenCorr)
			linkText := utf8Substr(t, e.Offset + lenCorr, e.Length)
			tail := utf8Substr(t, e.Offset+e.Length + lenCorr, -1)
			t = fmt.Sprintf("%s[%s](%s)%s", head, linkText, tt.Url, tail)
			lenCorr += int32(len([]rune(tt.Url)) + len([]rune("[]()")))
		case tdlib.TextEntityTypeUrlType: // example.com -> [example.com](example.com)
			head := utf8Substr(t, 0, e.Offset + lenCorr)
			linkText := utf8Substr(t, e.Offset + lenCorr, e.Length)
			tail := utf8Substr(t, e.Offset+e.Length + lenCorr, -1)
			t = fmt.Sprintf("%s[%s](%s)%s", head, linkText, linkText, tail)
			lenCorr += int32(len([]rune(linkText)) + len([]rune("[]()")))
		case tdlib.TextEntityTypeHashtagType:
		case tdlib.TextEntityTypePhoneNumberType:
		default:
			js, _ := json.MarshalIndent(e, "", "  ")
			errs = append(errs, fmt.Errorf("unsupported text entity type `%s`: %s", e.Type.GetTextEntityTypeEnum(), string(js)))
		}
	}
	return t, errs
}

func utf8Substr(s string, offset int32, length int32) (res string) {
	if length == -1 {
		length = int32(len(s))
	}
	var cc int32 = 0
	for _, ch := range s {
		if cc == (offset + length) {
			break
		}
		if cc >= offset {
			res += string(ch)
		}
		cc++
	}
	return res
}
