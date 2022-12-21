package main

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

type Resource string

var resources = []Resource{
	"ore",
	"clay",
	"obsidian",
	"geode",
}

// const (
// 	OreResource      Resource = "ore"
// 	ClayResource     Resource = "clay"
// 	ObsidianResource Resource = "obsidian"
// 	GeodeResource    Resource = "geode"
// )

// func ParseResource(str string) Resource {
// 	switch str {
// 	case "ore":
// 		return OreResource
// 	case "clay":
// 		return ClayResource
// 	case "obsidian":
// 		return ObsidianResource
// 	case "geode":
// 		return GeodeResource
// 	default:
// 		panic(fmt.Sprintf("unknown resource: %s", str))
// 	}
// }

type Cost map[Resource]int

func (cost Cost) Sub(other map[Resource]int) Cost {
	res := make(Cost, len(resources))

	for resource, count := range cost {
		res[resource] = count - other[resource]
	}

	return res
}

func AsSlice(res map[Resource]int) []ResourceCount {
	slice := make([]ResourceCount, 0, len(res))

	for resource, count := range res {
		slice = append(slice, ResourceCount{resource, count})
	}

	// Sort largest first
	sort.Slice(slice, func(i, j int) bool {
		return slice[i].Count > slice[j].Count
	})

	return slice
}

func ZeroValues(res map[Resource]int) []Resource {
	var out []Resource

	for resource, count := range res {
		if count == 0 {
			out = append(out, resource)
		}
	}

	return out
}

type ResourceCount struct {
	Resource Resource
	Count int
}

type Blueprint struct {
	Number int
	Costs  map[Resource]Cost
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
		Costs: make(map[Resource]Cost),
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

		robotCosts := make(map[Resource]int)

		for _, resource := range resourceRe.FindAllStringSubmatch(robotCostStr[len(robotMatch[0]):], -1) {
			robotCosts[Resource(resource[2])] = MustAtoi(resource[1])
		}

		blueprint.Costs[Resource(robotMatch[1])] = robotCosts
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
