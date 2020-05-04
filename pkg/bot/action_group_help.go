package bot

import "github.com/Arman92/go-tdlib"

func (b *Bot) ActionShowHelp(msg *tdlib.Message, single bool) error {
	if single {
		return b.sendText(
			msg,
			"Чтобы добавить одиночный пикет:\n"+
				"- будьте готовы к тому, что на карте будет ваше имя и юзернейм\n"+
				"- укажите его место, командой:\n"+
				"<pre><code>  %s широта, долгота</code></pre>\n"+
				"- отправьте боту текст воззвания, в формате <a href=\"https://github.com/adam-p/markdown-here/wiki/Markdown-Cheatsheet\">Markdown</a>:\n"+
				"<pre><code>  %s текст</code></pre>\n"+
				"- включите это все, командой:\n"+
				"<pre><code>  %s</code></pre>\n"+
				"\n"+
				"Если вы не знаете широты и долготы, отправьте локацию (работает только на телефонах), и бот вернет вам координаты.\n"+
				"Чтобы прекратить это все, отправьте <code>%s</code>.\n"+
				"\n"+
				"Чтобы добавить групповой митинг, вы должны быть администратором в публичной группе митинга. "+
				"Добавьте бота в группу и отправьте ему <code>%s@%s</code>.",
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
		"Чтобы добавить групповой митинг:\n\n"+
			"- укажите его место, командой:\n"+
			"<pre><code>  %s широта, долгота</code></pre>\n"+
			"- отправьте боту текст воззвания, в формате <a href=\"https://github.com/adam-p/markdown-here/wiki/Markdown-Cheatsheet\">Markdown</a>:\n"+
			"<pre><code>  %s текст</code></pre>\n"+
			"- включите это все, командой:\n"+
			"<pre><code>  %s</code></pre>\n"+
			"\n"+
			"Чтобы скрыть митинг, отправьте боту"+
			"<pre><code>  %s</code></pre>"+
			"\n"+
			"Чтобы создать одиночный пикет, напишите <code>%s</code> боту @%s напрямую.",
		locationCommand,
		textCommand,
		publishCommand,
		hideCommand,
		startCommand,
		b.me.Username,
	)
}
