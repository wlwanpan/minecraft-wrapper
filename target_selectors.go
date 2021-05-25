package wrapper

import (
	"fmt"
	"strings"
)

// TargetSelectorType - an enum type for the 5 different target selectors
type TargetSelectorType string

const (
	// AllPlayers - targets every player (alive or dead) by default
	AllPlayers TargetSelectorType = "a"
	// AllEntities - targets all alive entities in loaded chunks (includes players)
	AllEntities TargetSelectorType = "e"
	// NearestPlayer - targets the nearest player. When run by the console,
	// the origin of selection is (0,0,0). If there are multiple nearest
	// players, caused by them being precisely the same distance away, the
	// payer who most recently joined the server is selected
	NearestPlayer TargetSelectorType = "p"
	// RandomPlayer - targets a random player
	RandomPlayer TargetSelectorType = "r"
	// ExecutingEntity - targets the entity (alive or dead) that executed
	// the command. It does not target anyhing if the command was run by a
	// command block or server console
	ExecutingEntity TargetSelectorType = "s"
)

// TargetSelector - defines and creates a TargetSelector,
// complete with any arguments defined through argument functions
type TargetSelector struct {
	t     TargetSelectorType
	args  map[string]string
	types []string
	tags  []string
}

// NewTargetSelector - creates a new TargetSelector of the specified type
func NewTargetSelector(t TargetSelectorType) TargetSelector {
	return TargetSelector{
		t:    t,
		args: make(map[string]string),
	}
}

// String - returns, in the expected minecraft console format, the string
// version of the TargetSelector
func (s TargetSelector) String() string {
	var allArgs []string = make([]string, 0, len(s.args)+len(s.types)+len(s.tags))

	// args
	for key, value := range s.args {
		allArgs = append(allArgs, fmt.Sprintf("%s=%s", key, value))
	}

	// types
	for _, t := range s.types {
		allArgs = append(allArgs, "type="+t)
	}

	// tags
	for _, tag := range s.tags {
		allArgs = append(allArgs, "tag="+tag)
	}

	return fmt.Sprintf("@%s[%s]", s.t, strings.Join(allArgs, ","))
}

func (s TargetSelector) copyArgs() {
	var newMap = make(map[string]string)
	for key, val := range s.args {
		newMap[key] = val
	}
	s.args = newMap
}

func (s TargetSelector) copyStringSlice(slice []string) []string {
	return append([]string{}, slice...)
}

type positionalType string

// Positional argument types
const (
	X positionalType = "x"
	Y positionalType = "y"
	Z positionalType = "z"
)

// Positional - defines a positional argument for TargetSelectors
type Positional struct {
	Type     positionalType
	Value    float64
	Relative bool
}

// WithPositional - adds positional argument(s) to the TargetSelector.
// Returns a new TargetSelector with the positional argument(s) added;
// original is left unchanged. Allows for method chaining
func (s TargetSelector) WithPositional(args ...Positional) TargetSelector {
	s.copyArgs()
	for _, arg := range args {
		if arg.Relative {
			s.args[string(arg.Type)] = fmt.Sprintf("~%.2f", arg.Value)
		} else {
			s.args[string(arg.Type)] = fmt.Sprintf("%.2f", arg.Value)
		}
	}
	return s
}

// ExactDistance - specifes the exact distance targets must be from the point of command origin
type ExactDistance struct {
	Value uint
}

func (s ExactDistance) distance() string {
	return fmt.Sprintf("%d", s.Value)
}

// RangeDistance - specifes a range of distance targets can be from the point of command origin
type RangeDistance struct {
	Min, Max uint
}

func (s RangeDistance) distance() string {
	return fmt.Sprintf("%d..%d", s.Min, s.Max)
}

// Distance - accepts either Exacct or Range Distance types
type Distance interface {
	distance() string
}

// WithDistance - adds a diustance argument to the TargetSelector.
// Returns a new TargetSelector with the distance argument added;
// original is left unchanged. Allows for method chaining
func (s TargetSelector) WithDistance(distance Distance) TargetSelector {
	s.copyArgs()
	s.args["distance"] = distance.distance()
	return s
}

type volumeType string

// Volume argument types
const (
	DX volumeType = "dx"
	DY volumeType = "dy"
	DZ volumeType = "dz"
)

// Volume - defines a volume argument for TargetSelectors
type Volume struct {
	Type  volumeType
	Value float64
}

// WithVolume - adds volume argument(s) to the TargetSelector.
// Returns a new TargetSelector with the volume argument(s) added;
// original is left unchanged. Allows for method chaining
func (s TargetSelector) WithVolume(args ...Volume) TargetSelector {
	s.copyArgs()
	for _, arg := range args {
		s.args[string(arg.Type)] = fmt.Sprintf("%.2f", arg.Value)
	}
	return s
}

// ExactScore - specifes the exact score(s) targets must have
type ExactScore struct {
	Objective string
	Value     int
}

func (s ExactScore) score() string {
	return fmt.Sprintf("%s=%d", s.Objective, s.Value)
}

// RangeScore - specifes a range of score(s) targets can have
type RangeScore struct {
	Objective string
	Min, Max  int
}

func (s RangeScore) score() string {
	return fmt.Sprintf("%s=%d..%d", s.Objective, s.Min, s.Max)
}

// Score - accepts either an Exact or Range Score types
type Score interface {
	score() string
}

// WithScore - adds score argument(s) to the TargetSelector.
// Returns a new TargetSelector with the score argument(s) added;
// original is left unchanged. Allows for method chaining
func (s TargetSelector) WithScore(scores ...Score) TargetSelector {
	s.copyArgs()
	var scoreStrings []string
	for _, score := range scores {
		scoreStrings = append(scoreStrings, score.score())
	}
	scoresJoined := strings.Join(scoreStrings, ",")
	s.args["score"] = fmt.Sprintf("{%s}", scoresJoined)
	return s
}

// WithTeam - adds a team argument to the TargetSelector.
// Empty teamName disginates those not on a team.
// Returns a new TargetSelector with the team argument added;
// original is left unchanged. Allows for method chaining
func (s TargetSelector) WithTeam(teamName string, not bool) TargetSelector {
	s.copyArgs()
	if not {
		s.args["team"] = "!" + teamName
	} else {
		s.args["team"] = teamName
	}
	return s
}

type sortType string

// Sort argument types
const (
	NoSort    sortType = ""
	Nearest   sortType = "nearest"
	Furthest  sortType = "furthest"
	Random    sortType = "random"
	Arbitrary sortType = "arbitrary"
)

// WithLimit - adds limit and optionally sort arguments to the TargetSelector.
// Returns a new TargetSelector with the argument(s) added;
// original is left unchanged. Allows for method chaining
func (s TargetSelector) WithLimit(value uint, sort sortType) TargetSelector {
	s.copyArgs()
	s.args["limit"] = fmt.Sprintf("%d", value)
	if sort != NoSort {
		s.args["sort"] = string(sort)
	} else {
		delete(s.args, "sort")
	}
	return s
}

// ExactExperience - specifies the exact level players must be
type ExactExperience struct {
	Value uint
}

func (e ExactExperience) experience() string {
	return fmt.Sprintf("%d", e.Value)
}

// RangeExperience - specifes a range of levels players may be
type RangeExperience struct {
	Min, Max uint
}

func (e RangeExperience) experience() string {
	return fmt.Sprintf("%d..%d", e.Min, e.Max)
}

// Experience - accepts either an Exact or Range Experience types
type Experience interface {
	experience() string
}

// WithExperience - adds limit and optionally sort arguments to the TargetSelector.
// Returns a new TargetSelector with the argument(s) added;
// original is left unchanged. Allows for method chaining
func (s TargetSelector) WithExperience(exp Experience) TargetSelector {
	s.copyArgs()
	s.args["level"] = exp.experience()
	return s
}

// WithGameMode - adds a gamemode argument to the TargetSelector.
// Returns a new TargetSelector with the argument added;
// original is left unchanged. Allows for method chaining
func (s TargetSelector) WithGameMode(mode GameMode, not bool) TargetSelector {
	s.copyArgs()
	if not {
		s.args["gamemode"] = "!" + string(mode)
	} else {
		s.args["gamemode"] = string(mode)
	}
	return s
}

// WithName - adds a name argument to the TargetSelector.
// Returns a new TargetSelector with the argument added;
// original is left unchanged. Allows for method chaining
func (s TargetSelector) WithName(name string, not bool) TargetSelector {
	s.copyArgs()
	if not {
		s.args["name"] = "!" + name
	} else {
		s.args["name"] = name
	}
	return s
}

type rotationType string

// Rotation argument types
const (
	XRotation rotationType = "x_rotation"
	YRotation rotationType = "y_rotation"
)

// ExactRotation - specifies the exact rotation entities must be facing
type ExactRotation struct {
	Type  rotationType
	Value float64
}

func (r ExactRotation) rotationType() rotationType {
	return r.Type
}

func (r ExactRotation) rotation() string {
	return fmt.Sprintf("%.2f", r.Value)
}

// NewFloat64 - returns a pointer to the float value passed as an argument
func NewFloat64(f float64) *float64 { return &f }

// RangeRotation - specifies the range of rotation entities may be facing
type RangeRotation struct {
	Type     rotationType
	Min, Max *float64
}

func (r RangeRotation) rotationType() rotationType {
	return r.Type
}

func (r RangeRotation) rotation() string {
	var rotStr string
	if r.Min != nil {
		rotStr += fmt.Sprintf("%.2f", *r.Min)
	}
	rotStr += ".."
	if r.Max != nil {
		rotStr += fmt.Sprintf("%.2f", *r.Max)
	}
	return rotStr
}

// Rotation - accepts either Exact or Range Rotation types
type Rotation interface {
	rotationType() rotationType
	rotation() string
}

// WithRotation - adds rotation argument(s) to the TargetSelector.
// Returns a new TargetSelector with the argument(s) added;
// original is left unchanged. Allows for method chaining
func (s TargetSelector) WithRotation(rotations ...Rotation) TargetSelector {
	s.copyArgs()
	for _, rotation := range rotations {
		key := rotation.rotationType()
		val := rotation.rotation()
		s.args[string(key)] = val
	}
	return s
}

// TargetSelectorEntityType - an interface for entities, which consist of a namespace and name
// as well as a string format for use in commands. Public so users can satisfy the interface
// for mods and data packs which introduce their own entities. Built-in Minecraft Entity IDS and
// Minecraft Entity tags have been defined
type TargetSelectorEntityType interface {
	Namespace() string
	Name() string
	String() string
}

// TypeArg - defines a type argument for TargetSelectors
type TypeArg struct {
	Type TargetSelectorEntityType
	Not  bool
}

// WithType - Appends the type argument(s) to the current
// list of type arguments in the TargetSelector. Only for use with @e.
// Returns a new TargetSelector with the argument(s) added;
// original is left unchanged. Allows for method chaining
func (s TargetSelector) WithType(args ...TypeArg) TargetSelector {
	s.types = s.copyStringSlice(s.types)
	for _, arg := range args {
		if arg.Not {
			s.types = append(s.types, "!"+arg.Type.String())
		} else {
			s.types = append(s.types, arg.Type.String())
		}
	}
	return s
}

// Tag - defines a tag argument for TargetSelectors.
// Empty Name corresponds to all entities with exactly 0 tags
type Tag struct {
	Name string
	Not  bool
}

// WithDataTag - Appends the tag argument(s) to the current
// list of tag arguments in the TargetSelector.
// Returns a new TargetSelector with the argument(s) added;
// original is left unchanged. Allows for method chaining
func (s TargetSelector) WithDataTag(tags ...Tag) TargetSelector {
	s.tags = s.copyStringSlice(s.tags)
	for _, tag := range tags {
		if tag.Not {
			s.tags = append(s.tags, "!"+tag.Name)
		} else {
			s.tags = append(s.tags, tag.Name)
		}
	}
	return s
}

// TODO: WithNBT - requires NBT implementation

// TODO: WithAdvancements - requires achievement namespace id definitions

// TODO: WithPredicates - requires predicate implementation
