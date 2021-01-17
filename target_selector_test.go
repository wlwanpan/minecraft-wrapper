package wrapper

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

/*////////////////////////////
//////// Suite  Setup ////////
////////////////////////////*/

type TargetSelectorSuite struct {
	suite.Suite
}

func (s *TargetSelectorSuite) SetupSuite() {}

func (s *TargetSelectorSuite) SetupTest() {}

func (s *TargetSelectorSuite) TearDownTest() {}

func (s *TargetSelectorSuite) TearDownSuite() {}

func TestTargetSelectorSuite(t *testing.T) {
	suite.Run(t, new(TargetSelectorSuite))
}

/*////////////////////////////
// TargetSelectorType  Test //
////////////////////////////*/

func (s *TargetSelectorSuite) TestTargetSelectorType() {
	type testcase struct {
		Name               string
		TargetSelectorType TargetSelectorType
		TargetOutput       string
	}

	cases := []testcase{
		{
			Name:               "AllPlayers",
			TargetSelectorType: AllPlayers,
			TargetOutput:       "@a[]",
		},
		{
			Name:               "AllEntities",
			TargetSelectorType: AllEntities,
			TargetOutput:       "@e[]",
		},
		{
			Name:               "NearestPlayer",
			TargetSelectorType: NearestPlayer,
			TargetOutput:       "@p[]",
		},
		{
			Name:               "RandomPlayer",
			TargetSelectorType: RandomPlayer,
			TargetOutput:       "@r[]",
		},
		{
			Name:               "ExecutingEntity",
			TargetSelectorType: ExecutingEntity,
			TargetOutput:       "@s[]",
		},
	}

	for _, c := range cases {
		s.Run(c.Name, func() {
			s.Assert().Equal(
				c.TargetOutput,
				TargetSelector{t: c.TargetSelectorType}.String(),
			)
		})
	}
}

/*////////////////////////////
//// WithPositional  Test ////
////////////////////////////*/

func (s *TargetSelectorSuite) TestWithPositional() {
	type testcase struct {
		Name         string
		Positional   Positional
		TargetOutput string
	}

	cases := []testcase{
		{
			Name: "X",
			Positional: Positional{
				Type:     X,
				Value:    0,
				Relative: true,
			},
			TargetOutput: "@e[x=~0.00]",
		},
		{
			Name: "Y",
			Positional: Positional{
				Type:     Y,
				Value:    0,
				Relative: false,
			},
			TargetOutput: "@e[y=0.00]",
		},
		{
			Name: "Z",
			Positional: Positional{
				Type:     Z,
				Value:    100.0001,
				Relative: true,
			},
			TargetOutput: "@e[z=~100.00]",
		},
	}

	for _, c := range cases {
		s.Run(c.Name, func() {
			s.Assert().Equal(
				c.TargetOutput,
				NewTargetSelector(AllEntities).WithPositional(c.Positional).String(),
			)
		})
	}
}

/*////////////////////////////
///// WithDistance  Test /////
////////////////////////////*/

func (s *TargetSelectorSuite) TestWithDistance() {
	type testcase struct {
		Name         string
		Distance     Distance
		TargetOutput string
	}

	cases := []testcase{
		{
			Name:         "Exact",
			Distance:     &ExactDistance{10},
			TargetOutput: "@e[distance=10]",
		},
		{
			Name:         "Range",
			Distance:     &RangeDistance{10, 20},
			TargetOutput: "@e[distance=10..20]",
		},
	}

	for _, c := range cases {
		s.Run(c.Name, func() {
			s.Assert().Equal(
				c.TargetOutput,
				NewTargetSelector(AllEntities).WithDistance(c.Distance).String(),
			)
		})
	}
}

/*////////////////////////////
////// WithVolume  Test //////
////////////////////////////*/

func (s *TargetSelectorSuite) TestWithVolume() {
	type testcase struct {
		Name         string
		Volume       Volume
		TargetOutput string
	}

	cases := []testcase{
		{
			Name: "DX",
			Volume: Volume{
				Type:  DX,
				Value: 10,
			},
			TargetOutput: "@e[dx=10.00]",
		},
		{
			Name: "DY",
			Volume: Volume{
				Type:  DY,
				Value: 10,
			},
			TargetOutput: "@e[dy=10.00]",
		},
		{
			Name: "DZ",
			Volume: Volume{
				Type:  DZ,
				Value: 10,
			},
			TargetOutput: "@e[dz=10.00]",
		},
	}

	for _, c := range cases {
		s.Run(c.Name, func() {
			s.Assert().Equal(
				c.TargetOutput,
				NewTargetSelector(AllEntities).WithVolume(c.Volume).String(),
			)
		})
	}
}

/*////////////////////////////
////// WithScore  Test //////
////////////////////////////*/

func (s *TargetSelectorSuite) TestWithScore() {
	type testcase struct {
		Name         string
		Score        []Score
		TargetOutput string
	}

	cases := []testcase{
		{
			Name: "Exact",
			Score: []Score{
				ExactScore{
					Objective: "test",
					Value:     10,
				},
			},
			TargetOutput: `@e[score={test=10}]`,
		},
		{
			Name: "Range",
			Score: []Score{
				RangeScore{
					Objective: "test",
					Min:       10,
					Max:       20,
				},
			},
			TargetOutput: "@e[score={test=10..20}]",
		},
		{
			Name: "Multi",
			Score: []Score{
				ExactScore{
					Objective: "a",
					Value:     10,
				},
				RangeScore{
					Objective: "b",
					Min:       10,
					Max:       20,
				},
			},
			TargetOutput: "@e[score={a=10,b=10..20}]",
		},
	}

	for _, c := range cases {
		s.Run(c.Name, func() {
			s.Assert().Equal(
				c.TargetOutput,
				NewTargetSelector(AllEntities).WithScore(c.Score...).String(),
			)
		})
	}
}

/*////////////////////////////
/////// WithTeam  Test ///////
////////////////////////////*/

func (s *TargetSelectorSuite) TestWithTeam() {
	type testcase struct {
		Name         string
		Team         string
		Not          bool
		TargetOutput string
	}

	cases := []testcase{
		{
			Name:         "TeamA",
			Team:         "A",
			Not:          false,
			TargetOutput: `@e[team=A]`,
		},
		{
			Name:         "NotTeamA",
			Team:         "A",
			Not:          true,
			TargetOutput: `@e[team=!A]`,
		},
		{
			Name:         "Teamless",
			Team:         "",
			Not:          false,
			TargetOutput: `@e[team=]`,
		},
		{
			Name:         "NotTeamless",
			Team:         "",
			Not:          true,
			TargetOutput: `@e[team=!]`,
		},
	}

	for _, c := range cases {
		s.Run(c.Name, func() {
			s.Assert().Equal(
				c.TargetOutput,
				NewTargetSelector(AllEntities).WithTeam(c.Team, c.Not).String(),
			)
		})
	}
}

/*////////////////////////////
/////// WithLimit Test ///////
////////////////////////////*/

func (s *TargetSelectorSuite) TestWithLimit() {
	type testcase struct {
		Name         string
		Limit        uint
		Sort         sortType
		TargetOutput string
	}

	cases := []testcase{
		{
			Name:         "WithoutSort",
			Limit:        10,
			Sort:         NoSort,
			TargetOutput: `@e[limit=10]`,
		},
		{
			Name:  "WithSort",
			Limit: 10,
			Sort:  Arbitrary,
			// TODO: Map order is not guaranteed, needs fixing
			TargetOutput: `@e[limit=10,sort=arbitrary]`,
		},
	}

	for _, c := range cases {
		s.Run(c.Name, func() {
			s.Assert().Equal(
				c.TargetOutput,
				NewTargetSelector(AllEntities).WithLimit(c.Limit, c.Sort).String(),
			)
		})
	}
}

/*////////////////////////////
//// WithExperience  Test ////
////////////////////////////*/

func (s *TargetSelectorSuite) TestWithExperience() {
	type testcase struct {
		Name         string
		Experience   Experience
		TargetOutput string
	}

	cases := []testcase{
		{
			Name:         "ExactExperience",
			Experience:   &ExactExperience{10},
			TargetOutput: `@e[level=10]`,
		},
		{
			Name:         "RangeExperience",
			Experience:   &RangeExperience{10, 20},
			TargetOutput: `@e[level=10..20]`,
		},
	}

	for _, c := range cases {
		s.Run(c.Name, func() {
			s.Assert().Equal(
				c.TargetOutput,
				NewTargetSelector(AllEntities).WithExperience(c.Experience).String(),
			)
		})
	}
}

/*////////////////////////////
///// WithGameMode  Test /////
////////////////////////////*/

func (s *TargetSelectorSuite) TestWithGameMode() {
	type testcase struct {
		Name         string
		GameMode     GameMode
		Not          bool
		TargetOutput string
	}

	cases := []testcase{
		{
			Name:         "Spectator",
			GameMode:     Spectator,
			Not:          false,
			TargetOutput: `@e[gamemode=spectator]`,
		},
		{
			Name:         "NotSpectator",
			GameMode:     Spectator,
			Not:          true,
			TargetOutput: `@e[gamemode=!spectator]`,
		},
		{
			Name:         "Adventure",
			GameMode:     Adventure,
			Not:          false,
			TargetOutput: `@e[gamemode=adventure]`,
		},
		{
			Name:         "NotAdventure",
			GameMode:     Adventure,
			Not:          true,
			TargetOutput: `@e[gamemode=!adventure]`,
		},
		{
			Name:         "Creative",
			GameMode:     Creative,
			Not:          false,
			TargetOutput: `@e[gamemode=creative]`,
		},
		{
			Name:         "NotCreative",
			GameMode:     Creative,
			Not:          true,
			TargetOutput: `@e[gamemode=!creative]`,
		},
		{
			Name:         "Survival",
			GameMode:     Survival,
			Not:          false,
			TargetOutput: `@e[gamemode=survival]`,
		},
		{
			Name:         "NotSurvival",
			GameMode:     Survival,
			Not:          true,
			TargetOutput: `@e[gamemode=!survival]`,
		},
	}

	for _, c := range cases {
		s.Run(c.Name, func() {
			s.Assert().Equal(
				c.TargetOutput,
				NewTargetSelector(AllEntities).WithGameMode(c.GameMode, c.Not).String(),
			)
		})
	}
}

/*////////////////////////////
/////// WithName  Test ///////
////////////////////////////*/

func (s *TargetSelectorSuite) TestWithName() {
	type testcase struct {
		Name         string
		PlayerName   string
		Not          bool
		TargetOutput string
	}

	cases := []testcase{
		{
			Name:         "Name",
			PlayerName:   "Steve",
			Not:          false,
			TargetOutput: `@e[name=Steve]`,
		},
		{
			Name:         "NotName",
			PlayerName:   "Steve",
			Not:          true,
			TargetOutput: `@e[name=!Steve]`,
		},
	}

	for _, c := range cases {
		s.Run(c.Name, func() {
			s.Assert().Equal(
				c.TargetOutput,
				NewTargetSelector(AllEntities).WithName(c.PlayerName, c.Not).String(),
			)
		})
	}
}

/*////////////////////////////
///// WithRotation  Test /////
////////////////////////////*/

func (s *TargetSelectorSuite) TestWithRotation() {
	type testcase struct {
		Name         string
		Rotation     Rotation
		TargetOutput string
	}

	cases := []testcase{
		{
			Name:         "Exact",
			Rotation:     &ExactRotation{XRotation, 10},
			TargetOutput: `@e[x_rotation=10.00]`,
		},
		{
			Name:         "Range",
			Rotation:     &RangeRotation{YRotation, NewFloat64(10), NewFloat64(20)},
			TargetOutput: `@e[y_rotation=10.00..20.00]`,
		},
		{
			Name:         "RangeMax",
			Rotation:     &RangeRotation{YRotation, nil, NewFloat64(20)},
			TargetOutput: `@e[y_rotation=..20.00]`,
		},
		{
			Name:         "RangeMin",
			Rotation:     &RangeRotation{YRotation, NewFloat64(10), nil},
			TargetOutput: `@e[y_rotation=10.00..]`,
		},
	}

	for _, c := range cases {
		s.Run(c.Name, func() {
			s.Assert().Equal(
				c.TargetOutput,
				NewTargetSelector(AllEntities).WithRotation(c.Rotation).String(),
			)
		})
	}
}

/*////////////////////////////
/////// WithType  Test ///////
////////////////////////////*/

func (s *TargetSelectorSuite) TestWithType() {
	type testcase struct {
		Name         string
		Types        []TypeArg
		TargetOutput string
	}

	cases := []testcase{
		{
			Name:         "None",
			Types:        nil,
			TargetOutput: `@e[]`,
		},
		{
			Name: "One",
			Types: []TypeArg{
				{
					Type: PlayerEntity,
					Not:  false,
				},
			},
			TargetOutput: `@e[type=minecraft:player]`,
		},
		{
			Name: "NotOne",
			Types: []TypeArg{
				{
					Type: PlayerEntity,
					Not:  true,
				},
			},
			TargetOutput: `@e[type=!minecraft:player]`,
		},
		{
			Name: "Multiple",
			Types: []TypeArg{
				{
					Type: PlayerEntity,
					Not:  true,
				},
				{
					Type: EnderDragonEntity,
					Not:  false,
				},
			},
			TargetOutput: `@e[type=!minecraft:player,type=minecraft:ender_dragon]`,
		},
	}

	for _, c := range cases {
		s.Run(c.Name, func() {
			s.Assert().Equal(
				c.TargetOutput,
				NewTargetSelector(AllEntities).WithType(c.Types...).String(),
			)
		})
	}
}

/*////////////////////////////
/////// WithType  Test ///////
////////////////////////////*/

func (s *TargetSelectorSuite) TestWithDataTag() {
	type testcase struct {
		Name         string
		Tags         []Tag
		TargetOutput string
	}

	cases := []testcase{
		{
			Name:         "None",
			Tags:         nil,
			TargetOutput: `@e[]`,
		},
		{
			Name: "One",
			Tags: []Tag{
				{
					Name: "A",
					Not:  false,
				},
			},
			TargetOutput: `@e[tag=A]`,
		},
		{
			Name: "NotOne",
			Tags: []Tag{
				{
					Name: "A",
					Not:  true,
				},
			},
			TargetOutput: `@e[tag=!A]`,
		},
		{
			Name: "Multiple",
			Tags: []Tag{
				{
					Name: "A",
					Not:  true,
				},
				{
					Name: "B",
					Not:  false,
				},
			},
			TargetOutput: `@e[tag=!A,tag=B]`,
		},
	}

	for _, c := range cases {
		s.Run(c.Name, func() {
			s.Assert().Equal(
				c.TargetOutput,
				NewTargetSelector(AllEntities).WithDataTag(c.Tags...).String(),
			)
		})
	}
}

/*////////////////////////////
//////// All  Example ////////
////////////////////////////*/

func (s *TargetSelectorSuite) TestAll() {
	expected := TargetSelector{
		t: AllEntities,
		args: map[string]string{
			"x":          "1.00",
			"distance":   "2",
			"dz":         "3.00",
			"score":      "{score=4..5}",
			"team":       "!team",
			"limit":      "6",
			"sort":       "arbitrary",
			"level":      "7",
			"gamemode":   "survival",
			"y_rotation": "..8.00",
		},
		types: []string{"!minecraft:player"},
		tags:  []string{"tag"},
	}

	actual := NewTargetSelector(AllEntities).
		WithPositional(Positional{X, 1, false}).
		WithDistance(ExactDistance{2}).
		WithVolume(Volume{DZ, 3}).
		WithScore(RangeScore{"score", 4, 5}).
		WithTeam("team", true).
		WithLimit(6, Arbitrary).
		WithExperience(ExactExperience{7}).
		WithGameMode(Survival, false).
		WithRotation(RangeRotation{YRotation, nil, NewFloat64(8)}).
		WithType(TypeArg{PlayerEntity, true}).
		WithDataTag(Tag{"tag", false})

	s.Assert().EqualValues(expected, actual)
}
