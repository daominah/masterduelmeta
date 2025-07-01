package masterduelmeta

// some archetypes that need to check deck list to determine the real purpose of the deck
const (
	BarrierStatue    Archetype = "Barrier Statue"
	Branded          Archetype = "Branded"
	Bystial          Archetype = "Bystial"
	Centurion        Archetype = "Centur-Ion"
	DarkMagician     Archetype = "Dark Magician"
	Dragoon          Archetype = "Dragoon"
	FireKing         Archetype = "Fire King"
	Horus            Archetype = "Horus"
	InfernobleKnight Archetype = "Infernoble Knight"
	LairOfDarkness   Archetype = "Lair of Darkness"
	Rank8Go2nd       Archetype = "Rank 8 go 2nd"
	SnakeEye         Archetype = "Snake-Eye"
	Tearlaments      Archetype = "Tearlaments"
)

// NormalizeDeckTypeName returns the main archetype used in the deck,
// if 2 archetypes are in the deckTypeName, usually only return one of them.
//
// Sometimes, a deck list check follows by this function is needed,
// e.g. "Snake-Eye Fire King" is not sure if it is "Fire King" or "Snake-Eye",
// needs to check if the deck list contains "Legendary Fire King Ponix",
// the logic is implemented in the func Deck.Archetype
func NormalizeDeckTypeName(deckTypeName string) Archetype {
	switch deckTypeName {
	case "Adventure Combo", "Adventure Tenyi":
		return "Auroradon Combo"

	case "Stun":
		return "Barrier Statue"

	case "Branded", "Branded Tearlaments", "Despia":
		return "Branded"

	case "Centur-Ion", "Centurion", "Stardust Bystial":
		return "Centur-Ion"

	case "Danger Dark World", "Dark World":
		return "Dark World"

	case "Dinos":
		return "Dinosaur"

	case "Dogmatika", "Invoked Dogmatika":
		return "Dogmatika"

	case "Dragon Link", "Buster Blader Dragon Link":
		return "Dragon Link"

	case "Eldlich", "Zombie World Eldlich":
		return "Eldlich"

	case "Fiendsmith Control":
		return "Fiendsmith"

	case "Fire King", "Fire King Tri-Brigade", "Snake-Eye Fire King":
		return "Fire King"

	case "Gravekeeper's":
		return "Gravekeeper"

	case "HEROs", "Evil HEROs":
		return "HERO"

	case "Infernoble Knight", "Warrior Combo": // "Flame Swordsman" is a variant
		return "Infernoble Knight"

	case "Invoked", "Invoked Dogmatika Shaddoll", "Invoked Shaddoll":
		return "Invoke"

	case "Live☆Twin", "Live☆Twin Spright":
		return "Live☆Twin"

	case "Lyrilusc Tri-Brigade":
		return "Lyrilusc"

	case "Earth Machine":
		return "Machina"

	case "Madolche", "Eldlich Madolche", "Madolche Tri-Brigade":
		return "Madolche"

	case "Magical Musket", "Magical Musketeer":
		return "Magical Musketeer"

	case "@Ignister", "Code Talker", "Mathmech":
		return "Mathmech"

	case "Phantom Knight Orcust":
		return "Orcust"

	case "Pendulum Magician", "Supreme King", "Supreme King Melodious":
		return "Pendulum"

	case "Adventure Prank-Kids", "Prank-Kids":
		return "Prank-Kids"

	case "Adventure Phantom Knights", "Phantom Knights":
		return "Phantom Knights"

	case "Aroma", "Ragnaraika", "Rikka", "Sunavalon", "Plants":
		return "Plant"

	case "Predaplant", "Branded Predaplant":
		return "Predaplant"

	case "8-Axis Blind Second", "Blind Second":
		return "Rank 8 go 2nd"

	case "Resonators":
		return "Red Dragon Archfiend"

	case "Rescue-ACE", "Rescue Ace", "Snake-Eye Rescue Ace":
		return "Rescue-ACE"

	case "Ritual Beasts":
		return "Ritual Beast"

	case "Bystial Runick":
		return "Runick"

	case "Snake-Eye", "Snake-Eye Fiendsmith":
		return "Snake-Eye"

	case "Spright", "Runick Spright", "Tri-Brigade Spright":
		return "Spright"

	case "Swordsoul", "Swordsoul Tenyi", "Swordsoul Despia":
		return "Swordsoul"

	case "Synchrons":
		return "Synchron"

	case "Galaxy Tachyon":
		return "Tachyon"

	case "T.G.":
		return "Tech Genus"

	case "Traptrix", "Traptrix Dogmatika":
		return "Traptrix"

	case "Tri-Brigade Zoodiac", "Zoodiac Tri-Brigade":
		return "Tri-Brigade"

	case "Umi Control":
		return "Umi"

	case "Unchained", "Live☆Twin Unchained":
		return "Unchained"

	case "White Forest", "White Forest Azamina", "Azamina":
		return "White Forest"

	case "Yubel", "Yubel Fiendsmith":
		return "Yubel"

	default:
		return Archetype(deckTypeName)
	}
}
