package parse

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/spf13/cobra"
)

type StreamEntry struct {
	Channel   *string `json:"channel"`
	Feed      *string `json:"feed"`
	Title     string  `json:"title"`
	URL       *string `json:"url"`
	Quality   *string `json:"quality"`
	UserAgent *string `json:"user_agent"`
	Referrer  *string `json:"referrer"`
}

type LogoEntry struct {
	Channel *string  `json:"channel"`
	Feed    *string  `json:"feed"`
	Tags    []string `json:"tags"`
	Width   float64  `json:"width"`
	Height  float64  `json:"height"`
	Format  string   `json:"format"`
	URL     string   `json:"url"`
}

type CategoryEntry struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	Categories []string `json:"categories"`
}

func ParseCommand(app *pocketbase.PocketBase) *cobra.Command {
	return &cobra.Command{
		Use:   "parse",
		Short: "Parse streams.json and import to PocketBase",
		Run: func(cmd *cobra.Command, args []string) {
			if err := runParse(app); err != nil {
				log.Fatal(err)
			}
		},
	}
}

func runParse(app *pocketbase.PocketBase) error {
	// Read streams JSON file
	jsonPath := filepath.Join("pkg", "json", "streams.json")
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		return fmt.Errorf("failed to read streams JSON file: %w", err)
	}

	var streams []StreamEntry
	if err := json.Unmarshal(data, &streams); err != nil {
		return fmt.Errorf("failed to parse streams JSON: %w", err)
	}

	log.Printf("Found %d streams in JSON file\n", len(streams))

	// Read logos JSON file
	logosPath := filepath.Join("pkg", "json", "logos.json")
	logosData, err := os.ReadFile(logosPath)
	if err != nil {
		return fmt.Errorf("failed to read logos JSON file: %w", err)
	}

	var logos []LogoEntry
	if err := json.Unmarshal(logosData, &logos); err != nil {
		return fmt.Errorf("failed to parse logos JSON: %w", err)
	}

	log.Printf("Found %d logos in JSON file\n", len(logos))

	// Build logo map (lowercase channel name -> logo entry)
	logoMap := make(map[string]*LogoEntry)
	for i := range logos {
		if logos[i].Channel != nil && *logos[i].Channel != "" {
			logoMap[strings.ToLower(*logos[i].Channel)] = &logos[i]
		}
	}

	log.Printf("Built logo map with %d entries\n", len(logoMap))

	// Read categories JSON file
	categoriesPath := filepath.Join("pkg", "json", "categories.json")
	categoriesData, err := os.ReadFile(categoriesPath)
	if err != nil {
		return fmt.Errorf("failed to read categories JSON file: %w", err)
	}

	var categories []CategoryEntry
	if err := json.Unmarshal(categoriesData, &categories); err != nil {
		return fmt.Errorf("failed to parse categories JSON: %w", err)
	}

	log.Printf("Found %d categories in JSON file\n", len(categories))

	// Build category map (lowercase channel id -> category entry)
	categoryMap := make(map[string]*CategoryEntry)
	for i := range categories {
		if categories[i].ID != "" {
			categoryMap[strings.ToLower(categories[i].ID)] = &categories[i]
		}
	}

	log.Printf("Built category map with %d entries\n", len(categoryMap))

	// Get collections
	channelsCollection, err := app.FindCollectionByNameOrId("channels")
	if err != nil {
		return fmt.Errorf("failed to find channels collection: %w", err)
	}

	qualitiesCollection, err := app.FindCollectionByNameOrId("qualities")
	if err != nil {
		return fmt.Errorf("failed to find qualities collection: %w", err)
	}

	countriesCollection, err := app.FindCollectionByNameOrId("countries")
	if err != nil {
		return fmt.Errorf("failed to find countries collection: %w", err)
	}

	languagesCollection, err := app.FindCollectionByNameOrId("languages")
	if err != nil {
		return fmt.Errorf("failed to find languages collection: %w", err)
	}

	logosCollection, err := app.FindCollectionByNameOrId("logos")
	if err != nil {
		return fmt.Errorf("failed to find logos collection: %w", err)
	}

	categoriesCollection, err := app.FindCollectionByNameOrId("categories")
	if err != nil {
		return fmt.Errorf("failed to find categories collection: %w", err)
	}

	// Get or create the "empty" category for channels without categories
	emptyCategoryID, err := getOrCreateEmptyCategory(app, categoriesCollection)
	if err != nil {
		return fmt.Errorf("failed to get/create empty category: %w", err)
	}
	log.Printf("Using empty category ID: %s for channels without categories\n", emptyCategoryID)

	imported := 0
	skipped := 0
	logosImported := 0
	logosSkipped := 0

	for _, stream := range streams {
		// Skip if any required field is missing
		if stream.Channel == nil || *stream.Channel == "" {
			skipped++
			continue
		}
		if stream.URL == nil || *stream.URL == "" {
			skipped++
			continue
		}
		if stream.Quality == nil || *stream.Quality == "" {
			skipped++
			continue
		}
		if stream.Title == "" {
			skipped++
			continue
		}

		// Get or create quality record
		qualityID, err := getOrCreateQuality(app, qualitiesCollection, *stream.Quality)
		if err != nil {
			log.Printf("Warning: failed to get/create quality for %s: %v\n", stream.Title, err)
			skipped++
			continue
		}

		// Extract country from channel domain
		countryName := extractCountryFromChannel(*stream.Channel)
		var countryID string
		if countryName != "" {
			countryID, err = getOrCreateCountry(app, countriesCollection, countryName)
			if err != nil {
				log.Printf("Warning: failed to get/create country for %s: %v\n", stream.Title, err)
			}
		}

		// Extract language from title or channel domain
		languageName := extractLanguage(stream.Title, *stream.Channel)
		var languageID string
		if languageName != "" {
			languageID, err = getOrCreateLanguage(app, languagesCollection, languageName)
			if err != nil {
				log.Printf("Warning: failed to get/create language for %s: %v\n", stream.Title, err)
			}
		}

		// Get category for this channel
		var categoryID string
		channelNameLower := strings.ToLower(*stream.Channel)
		if categoryEntry, found := categoryMap[channelNameLower]; found && len(categoryEntry.Categories) > 0 {
			categoryID, err = getOrCreateCategory(app, categoriesCollection, categoryEntry.Categories)
			if err != nil {
				log.Printf("Warning: failed to get/create category for %s: %v\n", stream.Title, err)
				categoryID = emptyCategoryID
			}
		} else {
			// Use empty category for channels without categories
			categoryID = emptyCategoryID
		}

		// Create channel record
		channel := core.NewRecord(channelsCollection)
		channel.Set("channel", *stream.Channel)
		channel.Set("title", stream.Title)
		channel.Set("url", *stream.URL)
		channel.Set("quality", qualityID)
		if countryID != "" {
			channel.Set("country", countryID)
		}
		if languageID != "" {
			channel.Set("language", languageID)
		}
		if categoryID != "" {
			channel.Set("category", categoryID)
		}

		if err := app.Save(channel); err != nil {
			log.Printf("Warning: failed to save channel %s: %v\n", stream.Title, err)
			skipped++
			continue
		}

		imported++

		// Try to import logo for this channel
		if logoEntry, found := logoMap[channelNameLower]; found {
			logoRecord := core.NewRecord(logosCollection)
			logoRecord.Set("channel", channel.Id)
			logoRecord.Set("logo_url", logoEntry.URL)
			logoRecord.Set("width", logoEntry.Width)
			logoRecord.Set("height", logoEntry.Height)

			if err := app.Save(logoRecord); err != nil {
				log.Printf("Warning: failed to save logo for channel %s: %v\n", stream.Title, err)
				logosSkipped++
			} else {
				logosImported++
				
				// Update channel with logo relation
				channel.Set("logo", logoRecord.Id)
				if err := app.Save(channel); err != nil {
					log.Printf("Warning: failed to update channel %s with logo relation: %v\n", stream.Title, err)
				}
			}
		} else {
			logosSkipped++
		}

		if imported%100 == 0 {
			log.Printf("Imported %d channels, %d logos...\n", imported, logosImported)
		}
	}

	log.Printf("\nParse complete!\n")
	log.Printf("Imported: %d channels\n", imported)
	log.Printf("Skipped: %d channels (missing required fields)\n", skipped)
	log.Printf("Imported: %d logos\n", logosImported)
	log.Printf("Skipped: %d logos (not found or errors)\n", logosSkipped)

	return nil
}

func getOrCreateQuality(app *pocketbase.PocketBase, collection *core.Collection, qualityValue string) (string, error) {
	// Try to find existing quality
	records, err := app.FindRecordsByFilter(
		collection.Name,
		"quality = {:quality}",
		"-created",
		1,
		0,
		map[string]any{"quality": qualityValue},
	)
	if err != nil {
		return "", err
	}

	if len(records) > 0 {
		return records[0].Id, nil
	}

	// Create new quality record
	quality := core.NewRecord(collection)
	quality.Set("quality", qualityValue)

	if err := app.Save(quality); err != nil {
		return "", err
	}

	return quality.Id, nil
}

func getOrCreateCountry(app *pocketbase.PocketBase, collection *core.Collection, countryName string) (string, error) {
	// Try to find existing country
	records, err := app.FindRecordsByFilter(
		collection.Name,
		"name = {:name}",
		"-created",
		1,
		0,
		map[string]any{"name": countryName},
	)
	if err != nil {
		return "", err
	}

	if len(records) > 0 {
		return records[0].Id, nil
	}

	// Create new country record
	country := core.NewRecord(collection)
	country.Set("name", countryName)

	if err := app.Save(country); err != nil {
		return "", err
	}

	return country.Id, nil
}

func getOrCreateLanguage(app *pocketbase.PocketBase, collection *core.Collection, languageName string) (string, error) {
	// Try to find existing language
	records, err := app.FindRecordsByFilter(
		collection.Name,
		"name = {:name}",
		"-created",
		1,
		0,
		map[string]any{"name": languageName},
	)
	if err != nil {
		return "", err
	}

	if len(records) > 0 {
		return records[0].Id, nil
	}

	// Create new language record
	language := core.NewRecord(collection)
	language.Set("name", languageName)

	if err := app.Save(language); err != nil {
		return "", err
	}

	return language.Id, nil
}

func getOrCreateEmptyCategory(app *pocketbase.PocketBase, collection *core.Collection) (string, error) {
	// Try to find existing empty category (all fields empty or null)
	records, err := app.FindRecordsByFilter(
		collection.Name,
		"name_1 = '' && name_2 = '' && name_3 = ''",
		"-created",
		1,
		0,
		nil,
	)
	if err == nil && len(records) > 0 {
		return records[0].Id, nil
	}

	// Create new empty category record
	category := core.NewRecord(collection)
	category.Set("name_1", "")
	category.Set("name_2", "")
	category.Set("name_3", "")

	if err := app.Save(category); err != nil {
		return "", err
	}

	return category.Id, nil
}

func getOrCreateCategory(app *pocketbase.PocketBase, collection *core.Collection, categories []string) (string, error) {
	// Get up to 3 category names
	name1 := ""
	name2 := ""
	name3 := ""
	
	if len(categories) > 0 {
		name1 = categories[0]
	}
	if len(categories) > 1 {
		name2 = categories[1]
	}
	if len(categories) > 2 {
		name3 = categories[2]
	}

	// Build filter to find exact match - handle empty strings properly
	filterParts := []string{}
	params := map[string]any{}
	
	if name1 != "" {
		filterParts = append(filterParts, "name_1 = {:name1}")
		params["name1"] = name1
	} else {
		filterParts = append(filterParts, "(name_1 = '' || name_1 = null)")
	}
	
	if name2 != "" {
		filterParts = append(filterParts, "name_2 = {:name2}")
		params["name2"] = name2
	} else {
		filterParts = append(filterParts, "(name_2 = '' || name_2 = null)")
	}
	
	if name3 != "" {
		filterParts = append(filterParts, "name_3 = {:name3}")
		params["name3"] = name3
	} else {
		filterParts = append(filterParts, "(name_3 = '' || name_3 = null)")
	}
	
	filter := strings.Join(filterParts, " && ")

	records, err := app.FindRecordsByFilter(
		collection.Name,
		filter,
		"-created",
		1,
		0,
		params,
	)
	if err == nil && len(records) > 0 {
		return records[0].Id, nil
	}

	// Create new category record
	category := core.NewRecord(collection)
	category.Set("name_1", name1)
	category.Set("name_2", name2)
	category.Set("name_3", name3)

	if err := app.Save(category); err != nil {
		return "", err
	}

	return category.Id, nil
}

// extractCountryFromChannel extracts country name from channel domain
func extractCountryFromChannel(channel string) string {
	// Extract domain extension (e.g., .uz, .ru, .cn)
	parts := strings.Split(channel, ".")
	if len(parts) < 2 {
		return ""
	}

	domain := strings.ToLower(parts[len(parts)-1])

	// Map domain to country name - COMPREHENSIVE LIST
	countryMap := map[string]string{
		// A
		"ad": "Andorra",
		"ae": "United Arab Emirates",
		"af": "Afghanistan",
		"ag": "Antigua and Barbuda",
		"ai": "Anguilla",
		"al": "Albania",
		"am": "Armenia",
		"ao": "Angola",
		"ar": "Argentina",
		"as": "American Samoa",
		"at": "Austria",
		"au": "Australia",
		"aw": "Aruba",
		"az": "Azerbaijan",
		// B
		"ba": "Bosnia and Herzegovina",
		"bb": "Barbados",
		"bd": "Bangladesh",
		"be": "Belgium",
		"bf": "Burkina Faso",
		"bg": "Bulgaria",
		"bh": "Bahrain",
		"bi": "Burundi",
		"bj": "Benin",
		"bm": "Bermuda",
		"bn": "Brunei",
		"bo": "Bolivia",
		"br": "Brazil",
		"bs": "Bahamas",
		"bt": "Bhutan",
		"bw": "Botswana",
		"by": "Belarus",
		"bz": "Belize",
		// C
		"ca": "Canada",
		"cd": "Democratic Republic of the Congo",
		"cf": "Central African Republic",
		"cg": "Republic of the Congo",
		"ch": "Switzerland",
		"ci": "Ivory Coast",
		"ck": "Cook Islands",
		"cl": "Chile",
		"cm": "Cameroon",
		"cn": "China",
		"co": "Colombia",
		"cr": "Costa Rica",
		"cu": "Cuba",
		"cv": "Cape Verde",
		"cw": "Curaçao",
		"cy": "Cyprus",
		"cz": "Czech Republic",
		// D
		"de": "Germany",
		"dj": "Djibouti",
		"dk": "Denmark",
		"dm": "Dominica",
		"do": "Dominican Republic",
		"dz": "Algeria",
		// E
		"ec": "Ecuador",
		"ee": "Estonia",
		"eg": "Egypt",
		"er": "Eritrea",
		"es": "Spain",
		"et": "Ethiopia",
		// F
		"fi": "Finland",
		"fj": "Fiji",
		"fk": "Falkland Islands",
		"fm": "Micronesia",
		"fo": "Faroe Islands",
		"fr": "France",
		// G
		"ga": "Gabon",
		"gb": "United Kingdom",
		"gd": "Grenada",
		"ge": "Georgia",
		"gh": "Ghana",
		"gi": "Gibraltar",
		"gl": "Greenland",
		"gm": "Gambia",
		"gn": "Guinea",
		"gq": "Equatorial Guinea",
		"gr": "Greece",
		"gt": "Guatemala",
		"gu": "Guam",
		"gw": "Guinea-Bissau",
		"gy": "Guyana",
		// H
		"hk": "Hong Kong",
		"hn": "Honduras",
		"hr": "Croatia",
		"ht": "Haiti",
		"hu": "Hungary",
		// I
		"id": "Indonesia",
		"ie": "Ireland",
		"il": "Israel",
		"in": "India",
		"iq": "Iraq",
		"ir": "Iran",
		"is": "Iceland",
		"it": "Italy",
		// J
		"jm": "Jamaica",
		"jo": "Jordan",
		"jp": "Japan",
		// K
		"ke": "Kenya",
		"kg": "Kyrgyzstan",
		"kh": "Cambodia",
		"ki": "Kiribati",
		"km": "Comoros",
		"kn": "Saint Kitts and Nevis",
		"kp": "North Korea",
		"kr": "South Korea",
		"kw": "Kuwait",
		"ky": "Cayman Islands",
		"kz": "Kazakhstan",
		// L
		"la": "Laos",
		"lb": "Lebanon",
		"lc": "Saint Lucia",
		"li": "Liechtenstein",
		"lk": "Sri Lanka",
		"lr": "Liberia",
		"ls": "Lesotho",
		"lt": "Lithuania",
		"lu": "Luxembourg",
		"lv": "Latvia",
		"ly": "Libya",
		// M
		"ma": "Morocco",
		"mc": "Monaco",
		"md": "Moldova",
		"me": "Montenegro",
		"mg": "Madagascar",
		"mh": "Marshall Islands",
		"mk": "North Macedonia",
		"ml": "Mali",
		"mm": "Myanmar",
		"mn": "Mongolia",
		"mo": "Macau",
		"mr": "Mauritania",
		"ms": "Montserrat",
		"mt": "Malta",
		"mu": "Mauritius",
		"mv": "Maldives",
		"mw": "Malawi",
		"mx": "Mexico",
		"my": "Malaysia",
		"mz": "Mozambique",
		// N
		"na": "Namibia",
		"nc": "New Caledonia",
		"ne": "Niger",
		"ng": "Nigeria",
		"ni": "Nicaragua",
		"nl": "Netherlands",
		"no": "Norway",
		"np": "Nepal",
		"nr": "Nauru",
		"nu": "Niue",
		"nz": "New Zealand",
		// O
		"om": "Oman",
		// P
		"pa": "Panama",
		"pe": "Peru",
		"pf": "French Polynesia",
		"pg": "Papua New Guinea",
		"ph": "Philippines",
		"pk": "Pakistan",
		"pl": "Poland",
		"pm": "Saint Pierre and Miquelon",
		"pr": "Puerto Rico",
		"ps": "Palestine",
		"pt": "Portugal",
		"pw": "Palau",
		"py": "Paraguay",
		// Q
		"qa": "Qatar",
		// R
		"ro": "Romania",
		"rs": "Serbia",
		"ru": "Russia",
		"rw": "Rwanda",
		// S
		"sa": "Saudi Arabia",
		"sb": "Solomon Islands",
		"sc": "Seychelles",
		"sd": "Sudan",
		"se": "Sweden",
		"sg": "Singapore",
		"si": "Slovenia",
		"sk": "Slovakia",
		"sl": "Sierra Leone",
		"sm": "San Marino",
		"sn": "Senegal",
		"so": "Somalia",
		"sr": "Suriname",
		"ss": "South Sudan",
		"st": "São Tomé and Príncipe",
		"sv": "El Salvador",
		"sy": "Syria",
		"sz": "Eswatini",
		// T
		"tc": "Turks and Caicos Islands",
		"td": "Chad",
		"tg": "Togo",
		"th": "Thailand",
		"tj": "Tajikistan",
		"tl": "Timor-Leste",
		"tm": "Turkmenistan",
		"tn": "Tunisia",
		"to": "Tonga",
		"tr": "Turkey",
		"tt": "Trinidad and Tobago",
		"tv": "Tuvalu",
		"tw": "Taiwan",
		"tz": "Tanzania",
		// U
		"ua": "Ukraine",
		"ug": "Uganda",
		"uk": "United Kingdom",
		"us": "United States",
		"uy": "Uruguay",
		"uz": "Uzbekistan",
		// V
		"va": "Vatican City",
		"vc": "Saint Vincent and the Grenadines",
		"ve": "Venezuela",
		"vg": "British Virgin Islands",
		"vi": "US Virgin Islands",
		"vn": "Vietnam",
		"vu": "Vanuatu",
		// W
		"ws": "Samoa",
		// Y
		"ye": "Yemen",
		// Z
		"za": "South Africa",
		"zm": "Zambia",
		"zw": "Zimbabwe",
	}

	if country, ok := countryMap[domain]; ok {
		return country
	}

	return ""
}

// extractLanguage extracts language from title or channel domain
func extractLanguage(title, channel string) string {
	// First, try to extract language from title
	// Common language patterns in titles
	languagePatterns := map[string]*regexp.Regexp{
		"Albanian":   regexp.MustCompile(`(?i)\bAlbanian|Shqip\b`),
		"Arabic":     regexp.MustCompile(`(?i)\bArabic|العربية\b`),
		"Armenian":   regexp.MustCompile(`(?i)\bArmenian|Հայերեն\b`),
		"Azerbaijani": regexp.MustCompile(`(?i)\bAzerbaijani|Azərbaycan\b`),
		"Basque":     regexp.MustCompile(`(?i)\bBasque|Euskara\b`),
		"Belarusian": regexp.MustCompile(`(?i)\bBelarusian|Беларуская\b`),
		"Bengali":    regexp.MustCompile(`(?i)\bBengali|বাংলা\b`),
		"Bosnian":    regexp.MustCompile(`(?i)\bBosnian|Bosanski\b`),
		"Bulgarian":  regexp.MustCompile(`(?i)\bBulgarian|Български\b`),
		"Catalan":    regexp.MustCompile(`(?i)\bCatalan|Català\b`),
		"Chinese":    regexp.MustCompile(`(?i)\bChinese|中文|普通话|國語\b`),
		"Croatian":   regexp.MustCompile(`(?i)\bCroatian|Hrvatski\b`),
		"Czech":      regexp.MustCompile(`(?i)\bCzech|Čeština\b`),
		"Danish":     regexp.MustCompile(`(?i)\bDanish|Dansk\b`),
		"Dutch":      regexp.MustCompile(`(?i)\bDutch|Nederlands\b`),
		"English":    regexp.MustCompile(`(?i)\bEnglish\b`),
		"Estonian":   regexp.MustCompile(`(?i)\bEstonian|Eesti\b`),
		"Filipino":   regexp.MustCompile(`(?i)\bFilipino|Tagalog\b`),
		"Finnish":    regexp.MustCompile(`(?i)\bFinnish|Suomi\b`),
		"French":     regexp.MustCompile(`(?i)\bFrench|Français\b`),
		"Georgian":   regexp.MustCompile(`(?i)\bGeorgian|ქართული\b`),
		"German":     regexp.MustCompile(`(?i)\bGerman|Deutsch\b`),
		"Greek":      regexp.MustCompile(`(?i)\bGreek|Ελληνικά\b`),
		"Hebrew":     regexp.MustCompile(`(?i)\bHebrew|עברית\b`),
		"Hindi":      regexp.MustCompile(`(?i)\bHindi|हिन्दी\b`),
		"Hungarian":  regexp.MustCompile(`(?i)\bHungarian|Magyar\b`),
		"Icelandic":  regexp.MustCompile(`(?i)\bIcelandic|Íslenska\b`),
		"Indonesian": regexp.MustCompile(`(?i)\bIndonesian|Bahasa Indonesia\b`),
		"Italian":    regexp.MustCompile(`(?i)\bItalian|Italiano\b`),
		"Japanese":   regexp.MustCompile(`(?i)\bJapan(ese)?|日本語\b`),
		"Kazakh":     regexp.MustCompile(`(?i)\bKazakh|Қазақ\b`),
		"Korean":     regexp.MustCompile(`(?i)\bKorean|한국어\b`),
		"Latvian":    regexp.MustCompile(`(?i)\bLatvian|Latviešu\b`),
		"Lithuanian": regexp.MustCompile(`(?i)\bLithuanian|Lietuvių\b`),
		"Macedonian": regexp.MustCompile(`(?i)\bMacedonian|Македонски\b`),
		"Malay":      regexp.MustCompile(`(?i)\bMalay|Bahasa Melayu\b`),
		"Mongolian":  regexp.MustCompile(`(?i)\bMongolian|Монгол\b`),
		"Norwegian":  regexp.MustCompile(`(?i)\bNorwegian|Norsk\b`),
		"Persian":    regexp.MustCompile(`(?i)\bPersian|Farsi|فارسی\b`),
		"Polish":     regexp.MustCompile(`(?i)\bPolish|Polski\b`),
		"Portuguese": regexp.MustCompile(`(?i)\bPortuguese|Português\b`),
		"Romanian":   regexp.MustCompile(`(?i)\bRomanian|Română\b`),
		"Russian":    regexp.MustCompile(`(?i)\bRussian|Русский\b`),
		"Serbian":    regexp.MustCompile(`(?i)\bSerbian|Српски|Srpski\b`),
		"Slovak":     regexp.MustCompile(`(?i)\bSlovak|Slovenčina\b`),
		"Slovenian":  regexp.MustCompile(`(?i)\bSlovenian|Slovenščina\b`),
		"Spanish":    regexp.MustCompile(`(?i)\bSpanish|Español\b`),
		"Swedish":    regexp.MustCompile(`(?i)\bSwedish|Svenska\b`),
		"Thai":       regexp.MustCompile(`(?i)\bThai|ไทย\b`),
		"Turkish":    regexp.MustCompile(`(?i)\bTurkish|Türkçe\b`),
		"Ukrainian":  regexp.MustCompile(`(?i)\bUkrainian|Українська\b`),
		"Urdu":       regexp.MustCompile(`(?i)\bUrdu|اردو\b`),
		"Uzbek":      regexp.MustCompile(`(?i)\bUzbek|Oʻzbekcha\b`),
		"Vietnamese": regexp.MustCompile(`(?i)\bVietnamese|Tiếng Việt\b`),
	}

	// Check if title contains language name
	for lang, pattern := range languagePatterns {
		if pattern.MatchString(title) {
			return lang
		}
	}

	// If no language in title, use domain extension
	parts := strings.Split(channel, ".")
	if len(parts) < 2 {
		return "English" // Default to English
	}

	domain := strings.ToLower(parts[len(parts)-1])

	// Map domain to language - COMPREHENSIVE LIST
	languageMap := map[string]string{
		// A
		"ad": "Catalan",
		"ae": "Arabic",
		"af": "Pashto",
		"al": "Albanian",
		"am": "Armenian",
		"ao": "Portuguese",
		"ar": "Spanish",
		"at": "German",
		"au": "English",
		"az": "Azerbaijani",
		// B
		"ba": "Bosnian",
		"bb": "English",
		"bd": "Bengali",
		"be": "Dutch",
		"bg": "Bulgarian",
		"bh": "Arabic",
		"bn": "Malay",
		"bo": "Spanish",
		"br": "Portuguese",
		"by": "Belarusian",
		// C
		"ca": "English",
		"ch": "German",
		"cl": "Spanish",
		"cn": "Chinese",
		"co": "Spanish",
		"cr": "Spanish",
		"cu": "Spanish",
		"cy": "Greek",
		"cz": "Czech",
		// D
		"de": "German",
		"dk": "Danish",
		"do": "Spanish",
		"dz": "Arabic",
		// E
		"ec": "Spanish",
		"ee": "Estonian",
		"eg": "Arabic",
		"es": "Spanish",
		"et": "Amharic",
		// F
		"fi": "Finnish",
		"fj": "English",
		"fr": "French",
		// G
		"gb": "English",
		"ge": "Georgian",
		"gh": "English",
		"gr": "Greek",
		"gt": "Spanish",
		// H
		"hk": "Chinese",
		"hn": "Spanish",
		"hr": "Croatian",
		"ht": "French",
		"hu": "Hungarian",
		// I
		"id": "Indonesian",
		"ie": "English",
		"il": "Hebrew",
		"in": "Hindi",
		"iq": "Arabic",
		"ir": "Persian",
		"is": "Icelandic",
		"it": "Italian",
		// J
		"jm": "English",
		"jo": "Arabic",
		"jp": "Japanese",
		// K
		"ke": "Swahili",
		"kg": "Kyrgyz",
		"kh": "Khmer",
		"kn": "English",
		"kr": "Korean",
		"kw": "Arabic",
		"kz": "Kazakh",
		// L
		"la": "Lao",
		"lb": "Arabic",
		"lk": "Sinhala",
		"lt": "Lithuanian",
		"lu": "Luxembourgish",
		"lv": "Latvian",
		"ly": "Arabic",
		// M
		"ma": "Arabic",
		"mc": "French",
		"md": "Romanian",
		"me": "Montenegrin",
		"mk": "Macedonian",
		"mm": "Burmese",
		"mn": "Mongolian",
		"mo": "Chinese",
		"mt": "Maltese",
		"mx": "Spanish",
		"my": "Malay",
		"mz": "Portuguese",
		// N
		"ng": "English",
		"ni": "Spanish",
		"nl": "Dutch",
		"no": "Norwegian",
		"np": "Nepali",
		"nz": "English",
		// O
		"om": "Arabic",
		// P
		"pa": "Spanish",
		"pe": "Spanish",
		"ph": "Filipino",
		"pk": "Urdu",
		"pl": "Polish",
		"pr": "Spanish",
		"ps": "Arabic",
		"pt": "Portuguese",
		"py": "Spanish",
		// Q
		"qa": "Arabic",
		// R
		"ro": "Romanian",
		"rs": "Serbian",
		"ru": "Russian",
		// S
		"sa": "Arabic",
		"sd": "Arabic",
		"se": "Swedish",
		"sg": "English",
		"si": "Slovenian",
		"sk": "Slovak",
		"sn": "French",
		"so": "Somali",
		"sv": "Spanish",
		"sy": "Arabic",
		// T
		"th": "Thai",
		"tj": "Tajik",
		"tm": "Turkmen",
		"tn": "Arabic",
		"tr": "Turkish",
		"tt": "English",
		"tw": "Chinese",
		"tz": "Swahili",
		// U
		"ua": "Ukrainian",
		"ug": "English",
		"uk": "English",
		"us": "English",
		"uy": "Spanish",
		"uz": "Uzbek",
		// V
		"va": "Italian",
		"ve": "Spanish",
		"vn": "Vietnamese",
		// Y
		"ye": "Arabic",
		// Z
		"za": "English",
	}

	if language, ok := languageMap[domain]; ok {
		return language
	}

	return "English" // Default to English
}
