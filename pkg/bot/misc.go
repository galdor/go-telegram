package bot

import "net/url"

func BotURI(botName string) *url.URL {
	return &url.URL{
		Scheme: "https",
		Host:   "t.me",
		Path:   "/" + url.PathEscape(botName),
	}
}

func BotStartURI(botName, argument string) *url.URL {
	var query url.Values
	query.Add("start", argument)

	uri := BotURI(botName)
	uri.RawQuery = query.Encode()

	return uri
}

func BotGroupStartURI(botName, argument string) *url.URL {
	var query url.Values
	query.Add("startgroup", argument)

	uri := BotURI(botName)
	uri.RawQuery = query.Encode()

	return uri
}
