package bot

import (
	"github.com/Arman92/go-tdlib"
)

func (b *Bot) ActionSingleShowHelp(msg *tdlib.Message) error {
	return b.sendText(
		msg,
		"Чтобы добавить одиночный пикет, нужно всего-то:\n\n" +
			"- быть готовым к тому, что на карте будет ваше имя и юзернейм\n" +
			"- указать его место, командой\n" +
			"      <code>%s широта, долгота</code>\n" +
			"- отправить боту текст воззвания, в формате <a href=\"https://github.com/adam-p/markdown-here/wiki/Markdown-Cheatsheet\">Markdown</a>\n" +
			"      <code>%s текст</code>\n" +
			"- включить это все, командой\n" +
			"      <code>%s</code>" +
			"\n\n" +
			"Если вы не знаете широты и долготы, отправьте локацию (работает только на телефонах), и бот вернет вам координаты.\n" +
			"Чтобы прекратить это все, отправьте <code>%s</code>.\n" +
			"\n\n" +
			"Чтобы добавить групповой митинг, вы должны быть администратором в публичной группе митинга. " +
			"Добавьте бота в группу и отправьте ему <code>%s@%s</code>",
		locationCommand,
		textCommand,
		publishCommand,
		hideCommand,
		startCommand,
		b.me.Username,
	)
}