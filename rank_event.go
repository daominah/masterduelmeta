package masterduelmeta

import (
	"strings"
)

// enum values for RankedType in normal card pool
const (
	Top10Rating RankedType = "Top 10 Rating" // top 10 global, best decks at high competitive level

	MasterI    RankedType = "Master I"
	DiamondI   RankedType = "Diamond I"   // old highest rank, deprecated on 2023-06
	PlatinumI  RankedType = "Platinum I"  // old highest rank, deprecated on 2022-05
	WinStreaks RankedType = "Win Streaks" // 5 wins in a row at the highest rank type
	MasterV    RankedType = "Master V"
	DiamondV   RankedType = "Diamond V"

	DLvMAXDuelistCup        RankedType = "Duelist Cup DLv. Max"         // similar power to Master I
	DLvMAXWCS2023Qualifiers RankedType = "WCS 2023 Qualifiers DLv. Max" // similar power to Master I
	DLvMAXWCS2024Qualifiers RankedType = "WCS 2024 Qualifiers DLv. Max" // similar power to Master I

	WCSRegionalQualifiersWinStreaks RankedType = "WCS Regional Qualifiers Win Streaks"
	DuelistCupStage2WinStreaks      RankedType = "Duelist Cup Stage 2 Win Streaks"
)

// enum values for RankedType in fun events (different limited list)
const (
	DLvMAXXyzCup RankedType = "Xyz Cup DLv. Max"

	ThemeChronicle              RankedType = "Theme Chronicle"
	LegendAnthology             RankedType = "Legend Anthology"
	LegendAnthologyAcademy      RankedType = "Legend Anthology - Academy"
	LegendAnthologyAcceleration RankedType = "Legend Anthology - Acceleration"

	ExtraZeroFestival RankedType = "Extra Zero Festival"
	RitualFestival    RankedType = "Ritual Festival"

	FusionFestival  RankedType = "Fusion Festival"
	SynchroFestival RankedType = "Synchro Festival"
	XyzFestival     RankedType = "Xyz Festival"

	FusionXXyzFestival   RankedType = "Fusion x Xyz Festival"
	FusionXLinkFestival  RankedType = "Fusion x Link Festival"
	SynchroXXyzFestival  RankedType = "Synchro x Xyz Festival"
	SynchroXLinkFestival RankedType = "Synchro x Link Festival"
	XyzXLinkFestival     RankedType = "Xyz x Link Festival"

	DuelTriangleFusionSynchroXyz  RankedType = "Duel Triangle Fusion/Synchro/Xyz"
	DuelTriangleFusionSynchroLink RankedType = "Duel Triangle Fusion/Synchro/Link"

	DarkVsLight  RankedType = "Dark VS. Light"
	WaterAndWind RankedType = "Water and Wind"
	Attribute4   RankedType = "Attribute 4"

	MonsterTypeFestival             RankedType = "Monster Type Festival"
	MonsterTypeFestivalKingOfIsland RankedType = "Monster Type Festival: King of the Island"

	LimitOneFestival RankedType = "Limit One Festival"
	NRFestival       RankedType = "N/R Festival"

	AntiSpellFestival RankedType = "Anti-Spell Festival"

	SpecialDuelLinkRegulation RankedType = "Special Duel: Link Regulation"
)

func CheckIsNormalRank(rank RankedType) bool {
	switch rank {
	case Top10Rating,
		MasterI, DiamondI, PlatinumI, MasterV, DiamondV:
		return true
	}
	r := strings.ToLower(string(rank))
	if strings.Contains(r, "win streaks") {
		return true
	}
	if strings.Contains(r, "duelist cup") {
		return true
	}
	if strings.Contains(r, "dlv. max") {
		// not include Xyz Cup DLv. Max as normal rank
		return !strings.Contains(r, "xyz cup")
	}
	return false
}
