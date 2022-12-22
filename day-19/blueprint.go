package main

import (
	"fmt"
	"regexp"
	"strings"
)

type Resource int

func (res Resource) String() string {
	switch res {
	case OreResource:
		return "ore"
	case ClayResource:
		return "clay"
	case ObsidianResource:
		return "obsidian"
	case GeodeResource:
		return "geode"
	default:
		return fmt.Sprintf("unknown resource (%d)", res)
	}
}

const (
	NoneResource     Resource = -1
	OreResource      Resource = 0
	ClayResource     Resource = 1
	ObsidianResource Resource = 2
	GeodeResource    Resource = 3
)

type Resources [4]int

func (res Resources) String() string {
	var str strings.Builder

	for resource, count := range res {
		if str.Len() > 0 {
			str.WriteString(", ")
		}

		str.WriteString(fmt.Sprintf("%s: %d", Resource(resource), count))
	}

	return str.String()
}

func (res Resources) Sub(other Resources) Resources {
	var out Resources

	for resource, count := range res {
		out[resource] = count - other[resource]
	}

	return out
}

func (res Resources) Min(val int) Resources {
	var out Resources

	for resource, count := range res {
		if count < val {
			out[resource] = val
		} else {
			out[resource] = count
		}
	}

	return out
}

func (res Resources) AllZero() bool {
	for _, count := range res {
		if count != 0 {
			return false
		}
	}

	return true
}

func (res Resources) EqualValues(other Resources) []Resource {
	var out []Resource

	for resource, count := range res {
		if count > 0 && count == other[resource] {
			out = append(out, Resource(resource))
		}
	}

	return out
}

// var resources = []Resource{
// 	"ore",
// 	"clay",
// 	"obsidian",
// 	"geode",
// }

var resources = []Resource{
	OreResource,
	ClayResource,
	ObsidianResource,
	GeodeResource,
}

func ParseResource(str string) Resource {
	switch str {
	case "ore":
		return OreResource
	case "clay":
		return ClayResource
	case "obsidian":
		return ObsidianResource
	case "geode":
		return GeodeResource
	default:
		panic(fmt.Sprintf("unknown resource: %s", str))
	}
}

// type Cost map[Resource]int

// func (cost Cost) Sub(other map[Resource]int) Cost {
// 	res := make(Cost, len(resources))

// 	for resource, count := range cost {
// 		res[resource] = count - other[resource]
// 	}

// 	return res
// }

// func AsSlice(res map[Resource]int) []ResourceCount {
// 	slice := make([]ResourceCount, 0, len(res))

// 	for resource, count := range res {
// 		slice = append(slice, ResourceCount{resource, count})
// 	}

// 	// Sort largest first
// 	sort.Slice(slice, func(i, j int) bool {
// 		return slice[i].Count > slice[j].Count
// 	})

// 	return slice
// }

type ResourceCount struct {
	Resource Resource
	Count    int
}

type Blueprint struct {
	Number int
	Costs  map[Resource]Resources
	// OreRobot      RobotCost
	// ClayRobot     RobotCost
	// ObsidianRobot RobotCost
	// GeodeRobot    RobotCost
}

// I mean, it's simple and it works, so why not?
// var blueprintRe = regexp.MustCompile(`^Blueprint (\d+): Each ore robot costs (\d+) ore. Each clay robot costs (\d+) ore. Each obsidian robot costs (\d+) ore and (\d+) clay. Each geode robot costs (\d+) ore and (\d+) obsidian.$`)
var blueprintRe = regexp.MustCompile(`^Blueprint (\d+):`)

var (
	robotCostRe = regexp.MustCompile(`Each (\w+) robot costs`)
	resourceRe  = regexp.MustCompile(`(\d+)\s+(\w+)`)
)

func ParseBlueprint(str string) Blueprint {
	blueprint := Blueprint{
		Costs: make(map[Resource]Resources),
	}

	match := blueprintRe.FindStringSubmatch(str)
	if match == nil {
		panic(fmt.Sprintf("invalid input: '%s'", str))
	}

	blueprint.Number = MustAtoi(match[1])

	for _, robotCostStr := range strings.Split(str[len(match[0]):], ".") {
		robotCostStr = strings.TrimSpace(robotCostStr)
		if robotCostStr == "" {
			continue
		}

		robotMatch := robotCostRe.FindStringSubmatch(robotCostStr)
		if robotMatch == nil {
			panic(fmt.Sprintf("invalid input: '%s'", robotCostStr))
		}

		var robotCosts Resources

		for _, resource := range resourceRe.FindAllStringSubmatch(robotCostStr[len(robotMatch[0]):], -1) {
			robotCosts[ParseResource(resource[2])] = MustAtoi(resource[1])
		}

		blueprint.Costs[ParseResource(robotMatch[1])] = robotCosts
	}

	return blueprint
}

type RobotCost struct {
	Ore      int
	Clay     int
	Obsidian int
}

func (rc RobotCost) EnoughResources(ore, clay, obsidian int) bool {
	return rc.Ore <= ore && rc.Clay <= clay && rc.Obsidian <= obsidian
}
