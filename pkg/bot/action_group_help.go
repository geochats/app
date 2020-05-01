package bot

import "github.com/Arman92/go-tdlib"

func (b *Bot) ActionShowHelp(msg *tdlib.Message, single bool) error {
	if single {
		return b.sendText(
			msg,
			"Чтобы добавить одиночный пикет, нужно всего-то:\n\n" +
				"- быть готовым к тому, что на карте будет ваше имя и юзернейм\n" +
				"- указать его место, командой\n" +
				"      <pre>%s широта, долгота</pre>\n" +
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
	return b.sendText(
		msg,
		"Чтобы добавить групповой митинг, нужно:\n\n" +
			"- указать его место, командой\n" +
			"      <code>%s широта, долгота</code>\n" +
			"- отправить боту текст воззвания, в формате <a href=\"https://github.com/adam-p/markdown-here/wiki/Markdown-Cheatsheet\">Markdown</a>\n" +
			"      <code>%s текст</code>\n" +
			"- включить это все, командой\n" +
			"      <code>%s</code>" +
			"\n\n" +
			" Синтаксис текста:\n" +
			"  **жирный** \n" +
			"\n\n" +
			"Все это нужно делать из группы митинга. Если вы не знаете широты и долготы, отправьте локацию (работает только на телефонах), и бот вернет вам координаты.\n" +
			"Чтобы прекратить это все, отправьте <code>%s</code>.\n" +
			"\n\n" +
			"Чтобы создать одиночный пикет, напишите <code>%s</code> боту @%s напрямую.",
		locationCommand,
		textCommand,
		publishCommand,
		hideCommand,
		startCommand,
		b.me.Username,
	)
}
