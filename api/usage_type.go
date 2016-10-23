package main

import (
	"encoding/json"
	"net/http"

	"sort"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/user"
)

type NameOrigin struct {
	Code        string `datastore:"Code"`
	Plain       string `datastore:"Plain"`
	Description string `datastore:"Description"`
}

func getAllUsages(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	user := user.Current(ctx)
	userUsages := getAllUserUsages(user.Email, ctx)
	allUsages := getNameOrigins()
	var output UsageList

	for k := range allUsages {
		if userUsages[k] != nil {
			output = append(output, userUsages[k])
		} else {
			output = append(output,
				&SettingUsage{
					NameOrigin: allUsages[k],
					Enabled:    false,
					User:       user.Email,
				})
		}
	}

	sort.Sort(output)

	b, err := json.Marshal(output)
	if err != nil {
		log.Warningf(ctx, "Unable to marshal usages to json: %v", err)
	}
	w.Write(b)
}

func getNameOrigins() map[string]NameOrigin {
	result := make(map[string]NameOrigin)

	origins := []NameOrigin{
		{"afk", "Afrikaans", "Afrikaans names are used by Afrikaans speakers in the countries of South Africa and Namibia."},
		{"afr", "African", "African names are used in various places on the continent of Africa."},
		{"aka", "Akan", "Akan names are used by the Akan people of Ghana and Ivory Coast."},
		{"alb", "Albanian", "Albanian names are used in the country of Albania, as well as Kosovo and other Albanian communities throughout the world."},
		{"alg", "Algonquin", "Algonquin names are used by the Algonquin people of Ontario and Quebec in Canada."},
		{"ame", "Native American", "These names are or were used by the various indigenous peoples who inhabited North and South America."},
		{"amh", "Amharic", "Amharic names are used in Ethiopia."},
		{"anci", "Ancient", "These names were used in various ancient regions."},
		{"apa", "Apache", "Apache names are used by the Apache peoples of the southwestern United States."},
		{"ara", "Arabic", "Arabic names are used in the Arab world, as well as some other regions within the larger Muslim world. They are not necessarily of Arabic origin, though most in fact are. Compare also Persian names and Turkish names."},
		{"arm", "Armenian", "Armenian names are used in the country of Armenia in western Asia, as well as in Armenian diaspora communities throughout the world."},
		{"astr", "Astronomy", "These names occur primarily in astronomy. They are not commonly given to real people."},
		{"aus", "Indigenous", "Australian	Indigenous Australian names are used by the indigenous people of Australia (also called Aborigines)."},
		{"aym", "Aymara", "Aymara names are used by the Aymara people of Bolivia and Peru."},
		{"aze", "Azerbaijani", "Azerbaijani names are used by the Azeri people of Azerbaijan and northern Iran."},
		{"bal", "Balinese", "Balinese names are used on the island of Bali, a part of Indonesia."},
		{"bas", "Basque", "Basque names are used in the Basque Country (northern Spain and southern France)."},
		{"bel", "Belarusian", "Belarusian names are used in the country of Belarus in eastern Europe."},
		{"ben", "Bengali", "Bengali names are used in Bangladesh and eastern India."},
		{"ber", "Berber", "Berber names are used by the Berber (Amazigh) people of North Africa."},
		{"bib", "Biblical", "(All)	These names occur in the Bible (in any language)."},
		{"bos", "Bosnian", "Bosnian names are used by the Bosniak people. For additional names, see Serbo-Croatian names, Arabic names and Turkish names."},
		{"bre", "Breton", "Breton names are used in the region of Brittany in northwest France."},
		{"bsh", "Bashkir", "Bashkir names are used in Bashkortostan in Russia."},
		{"bul", "Bulgarian", "Bulgarian names are used in the country of Bulgaria in southeastern Europe."},
		{"cat", "Catalan", "Catalan names are used in Catalonia in eastern Spain, as well as in other Catalan-speaking areas including Valencia, the Balearic Islands, and Andorra."},
		{"cau", "Caucasian", "These names are used by the various ethnic groups of the Caucasus, a region between Europe and Asia."},
		{"cela", "Ancient Celtic", "-"},
		{"celm", "Celtic Mythology", "-"},
		{"cew", "Chewa", "Chewa names are used in Malawi, Mozambique and Zambia."},
		{"che", "Chechen", "Chechen names are used in Chechnya, a federal subject of Russia."},
		{"chi", "Chinese", "Chinese names are used in China and in Chinese communities throughout the world. Note that depending on the Chinese characters used these names can have many other meanings besides those listed here."},
		{"cht", "Choctaw", "Choctaw names are used by the Choctaw people of Oklahoma and Mississippi."},
		{"cir", "Circassian", "Circassian names are used in Circassia, part of the Caucasus region of Russia."},
		{"com", "Comanche", "Comanche names are used by the Comanche people of the southern United States."},
		{"cop", "Coptic", "These names are used by Coptic Christians in Egypt."},
		{"cor", "Cornish", "Cornish names were used in southwest England in the region around Cornwall."},
		{"cre", "Cree", "Cree names are used by the Cree people of Canada."},
		{"cro", "Croatian", "Croatian names are used in the country of Croatia and other Croatian communities throughout the world."},
		{"crs", "Corsican", "Corsican names are used on the island of Corsica."},
		{"cze", "Czech", "Czech names are used in the Czech Republic in central Europe."},
		{"dan", "Danish", "Danish names are used in the country of Denmark in northern Europe."},
		{"dgs", "Dagestani", "Dagestani names are used in Dagestan, a federal subject of Russia."},
		{"dut", "Dutch", "Dutch names are used in the Netherlands and Flanders."},
		{"eng", "English", "English names are used in English-speaking countries."},
		{"esp", "Esperanto", "Esperanto names are used by speakers of the planned language Esperanto."},
		{"est", "Estonian", "Estonian names are used in the country of Estonia in northern Europe."},
		{"eth", "Ethiopian", "Ethiopian names are used in the country of Ethiopia in eastern Africa."},
		{"ewe", "Ewe", "Ewe names are used by the Ewe people of Ghana and Togo."},
		{"fae", "Faroese", "Faroese names are used on the Faroe Islands."},
		{"fairy", "Fairy", "Fun category"},
		{"fil", "Filipino", "Filipino names are used on the island nation of the Philippines."},
		{"fin", "Finnish", "Finnish names are used in the country of Finland in northern Europe."},
		{"fle", "Flemish", "Flemish names are used in Flanders (the northern half of Belgium where Dutch is spoken)."},
		{"fre", "French", "French names are used in France and other French-speaking regions."},
		{"fri", "Frisian", "Frisian names are used in Friesland in the northern Netherlands and in East and North Frisia in northwestern Germany."},
		{"gal", "Galician", "Galician names are used in Galicia in northwestern Spain."},
		{"gan", "Ganda", "Ganda names are used by the Ganda people of Uganda."},
		{"geo", "Georgian", "Georgian names are used in the country of Georgia in central Eurasia."},
		{"ger", "German", "German names are used in Germany and other German-speaking areas such as Austria and Switzerland."},
		{"goth", "Goth", "Fun category"},
		{"gre", "Greek", "Greek names are used in the country of Greece and other Greek-speaking communities throughout the world."},
		{"grea", "Ancient Greek", "-"},
		{"grem", "Greek Mythology", "-"},
		{"grn", "Greenlandic", "Greenlandic names are used by the indigenous people of Greenland."},
		{"haw", "Hawaiian", "Hawaiian names are used by the indigenous people of Hawaii."},
		{"hb", "Hillbilly", "Fun category"},
		{"hil", "Hiligaynon", "Hiligaynon names are used in the Philippines."},
		{"hippy", "Hippy", "Fun category"},
		{"hist", "History", "These names are used primarily to refer to historical persons. They are not commonly used by other people."},
		{"hmo", "Hmong", "Hmong names are used by the Hmong people in southeastern Asia."},
		{"hun", "Hungarian", "Hungarian names are used in the country of Hungary in central Europe."},
		{"ibi", "Ibibio", "Ibibio names are used by the Ibibio people of Nigeria."},
		{"ice", "Icelandic", "Icelandic names are used on the island nation of Iceland."},
		{"igb", "Igbo", "Igbo names are used by the Igbo people of Nigeria."},
		{"ind", "Indian", "Indian names are used in India and in Indian communities throughout the world."},
		{"indm", "Indian Mythology", "-"},
		{"ing", "Ingush", "Ingush names are used in Ingushetia, a federal subject of Russia."},
		{"ins", "Indonesian", "Indonesian names are used on the island nation of Indonesia in southeast Asia."},
		{"inu", "Inuit", "Inuit names are used by the Inuit people of the North American Arctic."},
		{"ira", "Iranian", "Iranian names are used in the country of Iran in southwestern Asia."},
		{"iri", "Irish", "Irish names are used on the island of Ireland as well as elsewhere in the Western World as a result of the Irish diaspora."},
		{"iro", "Iroquois", "Iroquois names are used by the Iroquois people of the United States and Canada."},
		{"ita", "Italian", "Italian names are used in Italy and other Italian-speaking regions such as southern Switzerland."},
		{"jap", "Japanese", "Japanese names are used in Japan and in Japanese communities throughout the world. Note that depending on the Japanese characters used these names can have many other meanings besides those listed here."},
		{"jav", "Javanese", "Javanese names are used on the island of Java, a part of Indonesia."},
		{"jer", "Jèrriais", "Jèrriais names are used on the island of Jersey between Britain and France."},
		{"jew", "Jewish", "These names are used by Jews. For more specific lists, see Hebrew names and Yiddish names."},
		{"kaz", "Kazakh", "Kazakh names are used in the country of Kazakhstan in central Eurasia."},
		{"khm", "Khmer", "Khmer names are used in the country of Cambodia in southeastern Asia."},
		{"kik", "Kikuyu", "Kikuyu names are used by the Kikuyu people of Kenya."},
		{"kk", "Kreatyve", "Fun category"},
		{"kor", "Korean", "Korean names are used in South and North Korea. Note that depending on the Korean characters used these names can have many other meanings besides those listed here."},
		{"kur", "Kurdish", "Kurdish names are used by the Kurdish people of the Middle East."},
		{"kyr", "Kyrgyz", "Kyrgyz names are used in the country of Kyrgyzstan in central Eurasia."},
		{"lat", "Latvian", "Latvian names are used in the country of Latvia in northern Europe."},
		{"lim", "Limburgish", "Limburgish names are used in the Limburg region, which straddles the border between Belgium, the Netherlands and Germany."},
		{"lite", "Literature", "These names occur primarily in literature. They are not commonly given to real people."},
		{"lth", "Lithuanian", "Lithuanian names are used in the country of Lithuania in northern Europe."},
		{"luh", "Luhya", "Luhya names are used by the Luhya people of western Kenya."},
		{"luo", "Luo", "Luo names are used by the Luo people of Kenya."},
		{"mac", "Macedonian", "Macedonian names are used in the country of Macedonia (FYROM) in southeastern Europe."},
		{"mal", "Maltese", "Maltese names are used on the island of Malta."},
		{"man", "Manx", "Manx names are used on the Isle of Man."},
		{"mao", "Maori", "Maori names are used by the indigenous people of New Zealand, the Maori."},
		{"map", "Mapuche", "Mapuche names are used by the Mapuche people of Chile and Argentina."},
		{"may", "Mayan", "Mayan names are used by the Maya people of Mexico and Central America."},
		{"medi", "Medieval", "These names were used in medieval times."},
		{"mlm", "Malayalam", "Malayalam names are used in southern India."},
		{"mly", "Malay", "Malay names are used in Malaysia, Indonesia, Brunei, Singapore, and Thailand."},
		{"mon", "Mongolian", "Mongolian names are used in the country of Mongolia in central Asia."},
		{"morm", "Mormon", "Mormon names occur in the Book of Mormon."},
		{"mwe", "Mwera", "Mwera names are used by the Mwera people of Tanzania."},
		{"myth", "Mythology", "These names occur in mythology and religion."},
		{"nah", "Nahuatl", "Nahuatl names are used by the Nahua peoples of Mexico and Central America."},
		{"nav", "Navajo", "Navajo names are used by the Navajo people of the southwestern United States."},
		{"nde", "Ndebele", "Ndebele names are used by the Ndebele people of South Africa and Zimbabwe."},
		{"nor", "Norwegian", "Norwegian names are used in the country of Norway in northern Europe."},
		{"nuu", "Nuu-chah-nulth", "Nuu-chah-nulth names are used on Vancouver Island, Canada."},
		{"occ", "Occitan", "Occitan names are used in southern France and parts of Spain and Italy."},
		{"oji", "Ojibwe", "Ojibwe names are used by the Ojibwe people of Canada and the United States."},
		{"oss", "Ossetian", "Ossetian names are used in Ossetia, which is a region split by Russia and Georgia."},
		{"pas", "Pashto", "Pashto names are used in Afghanistan and northwestern Pakistan."},
		{"per", "Persian", "Persian names are used in the country of Iran, in southwestern Asia, which is part of the Muslim world."},
		{"pets", "Pet", "These names are most commonly given to pets: dogs, cats, etc."},
		{"pol", "Polish", "Polish names are used in the country of Poland in central Europe."},
		{"popu", "Popular Culture", "These names occur primarily in popular culture and entertainment. They are not commonly given to real people."},
		{"por", "Portuguese", "Portuguese names are used in Portugal, Brazil and other Portuguese-speaking areas."},
		{"pun", "Punjabi", "Punjabi names are used in the Punjab region of India and Pakistan."},
		{"que", "Quechua", "Quechua names are used by the Quechua people of South America."},
		{"rap", "Rapper", "Fun category"},
		{"rmn", "Romanian", "Romanian names are used in the countries of Romania and Moldova in eastern Europe."},
		{"roma", "Ancient Roman", "-"},
		{"romm", "Roman Mythology", "-"},
		{"rus", "Russian", "Russian names are used in the country of Russia and in Russian-speaking communities throughout the world."},
		{"sam", "Sami", "Sami names are used by the Sami people who inhabit northern Scandinavia."},
		{"sca", "Scandinavian", "Scandinavian names are used in the Scandinavia region of northern Europe. For more specific lists, see Swedish names, Danish names and Norwegian names."},
		{"scam", "Norse Mythology", "-"},
		{"sco", "Scottish", "Scottish names are used in the country of Scotland as well as elsewhere in the Western World as a result of the Scottish diaspora."},
		{"sct", "Scots", "Scots names are used by speakers of the Scots language."},
		{"ser", "Serbian", "Serbian names are used in the country of Serbia in southeastern Europe."},
		{"sha", "Shawnee", "Shawnee names are used by the Shawnee people of Oklahoma."},
		{"sho", "Shona", "Shona names are used by the Shona people of Zimbabwe."},
		{"sic", "Sicilian", "These names are a subset of Italian names used more often by speakers of Sicilian. Listed separately are Italian names."},
		{"sio", "Sioux", "Sioux names are used by the Sioux people of the central United States and Canada."},
		{"sla", "Slavic", "These names are used by Slavic peoples."},
		{"slam", "Norse Mythology", "-"},
		{"slk", "Slovak", "Slovak names are used in the country of Slovakia in central Europe."},
		{"sln", "Slovene", "Slovene names are used in the country of Slovenia in central Europe."},
		{"sor", "Sorbian", "Sorbian names are used in Lusatia in eastern Germany."},
		{"sot", "Sotho", "Sotho names are used by the Sotho people of Lesotho and South Africa."},
		{"spa", "Spanish", "Spanish names are used in Spain and other Spanish-speaking countries (such as those in South America)."},
		{"swa", "Swahili", "Swahili names are used by Swahili speakers in eastern Africa."},
		{"swe", "Swedish", "Swedish names are used in the country of Sweden in northern Europe."},
		{"tag", "Tagalog", "Tagalog names are used in the Philippines."},
		{"tah", "Tahitian", "Tahitian names are used in French Polynesia (part of which is the island of Tahiti)."},
		{"taj", "Tajik", "Tajik names are used in the country of Tajikistan in central Asia."},
		{"tam", "Tamil", "Tamil names are used in southern India and Sri Lanka."},
		{"tat", "Tatar", "Tatar names are used in Tatarstan in Russia."},
		{"tel", "Telugu", "Telugu names are used in eastern India."},
		{"tha", "Thai", "Thai names are used in the country of Thailand in southeastern Asia."},
		{"theo", "Theology", "These names are used to refer to (Judeo-Christian-Islamic) deities. They are not bestowed upon real people."},
		{"tib", "Tibetan", "Tibetan names are used by the Tibetan people who live in the region of Tibet in central Asia."},
		{"tkm", "Turkmen", "Turkmen names are used in the country of Turkmenistan."},
		{"trans", "Transformer", "Fun category"},
		{"tsw", "Tswana", "Tswana names are used in Botswana."},
		{"tum", "Tumbuka", "Tumbuka names are used in Malawi, Zambia and Tanzania."},
		{"tup", "Tupí", "Tupí names are used by the Tupí people of Brazil."},
		{"tur", "Turkish", "Turkish names are used in the country of Turkey, which is situated in western Asia and southeastern Europe. Turkey is part of the larger Muslim world."},
		{"ukr", "Ukrainian", "Ukrainian names are used in the country of Ukraine in eastern Europe."},
		{"urd", "Urdu", "Urdu names are used in Pakistan."},
		{"urh", "Urhobo", "Urhobo names are used by the Urhobo people of Nigeria."},
		{"usa", "American", "American names are used in the United States."},
		{"uyg", "Uyghur", "Uyghur names are used by the Uyghur people in western China."},
		{"uzb", "Uzbek", "Uzbek names are used in the country of Uzbekistan in central Asia."},
		{"vari", "Various", "These names do not belong to any one culture. They are put here because they cannot be categorized anywhere else."},
		{"vie", "Vietnamese", "Vietnamese names are used in the country of Vietnam in southeastern Asia."},
		{"wel", "Welsh", "Welsh names are used in the country of Wales in Britain."},
		{"witch", "Witch", "Fun category"},
		{"wrest", "Wrestler", "Fun category"},
		{"xho", "Xhosa", "Xhosa names are used by the Xhosa people of South Africa."},
		{"yao", "Yao", "Yao names are used in Malawi, Tanzania and Mozambique."},
		{"yor", "Yoruba", "Yoruba names are used by the Yoruba people of Nigeria."},
		{"zap", "Zapotec", "Zapotec names are used by the Zapotec people of southern Mexico."},
		{"zul", "Zulu", "Zulu names are used by the Zulu people of South Africa."},
	}

	for origin := range origins {
		result[origins[origin].Code] = origins[origin]
	}
	return result
}
