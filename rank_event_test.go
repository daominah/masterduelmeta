package masterduelmeta

import (
	"testing"
)

func TestCheckIsNormalRank(t *testing.T) {
	tests := []struct {
		rank     RankedType
		isNormal bool
	}{
		{Top10Rating, true},
		{MasterI, true},
		{DiamondI, true},
		{PlatinumI, true},
		{MasterV, true},
		{DiamondV, true},
		{WinStreaks, true},
		{DLvMAXDuelistCup, true},
		{DLvMAXWCS2023Qualifiers, true},
		{DLvMAXWCS2024Qualifiers, true},
		{WCSRegionalQualifiersWinStreaks, true},
		{DuelistCupStage2WinStreaks, true},

		{DLvMAXXyzCup, false},

		{ThemeChronicle, false},
		{LegendAnthology, false},
		{LegendAnthologyAcademy, false},
		{LegendAnthologyAcceleration, false},
		{ExtraZeroFestival, false},
		{RitualFestival, false},
		{FusionFestival, false},
		{SynchroFestival, false},
		{XyzFestival, false},
		{FusionXXyzFestival, false},
		{FusionXLinkFestival, false},
		{SynchroXXyzFestival, false},
		{SynchroXLinkFestival, false},
		{XyzXLinkFestival, false},
		{DuelTriangleFusionSynchroXyz, false},
		{DuelTriangleFusionSynchroLink, false},
		{DarkVsLight, false},
		{WaterAndWind, false},
		{Attribute4, false},
		{MonsterTypeFestival, false},
		{MonsterTypeFestivalKingOfIsland, false},
		{LimitOneFestival, false},
		{NRFestival, false},
		{AntiSpellFestival, false},
		{SpecialDuelLinkRegulation, false},
	}

	for _, test := range tests {
		if got := CheckIsNormalRank(test.rank); got != test.isNormal {
			t.Errorf("CheckIsNormalRank(%v) = %v, want %v", test.rank, got, test.isNormal)
		}
	}
}
