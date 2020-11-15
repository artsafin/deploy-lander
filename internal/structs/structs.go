package structs

import "regexp"

type ProjectMap map[string]Project

type Project struct {
	Alias         string
	Path          string
	TagsRe        string
	BranchesRe    string
	DeployCommand string
}

type BuildVersioner interface {
	Value() string
	IsTag(re string) bool
	IsBranch(re string) bool
}

func matches(re, value string) bool {
	if len(re) == 0 {
		return false
	}
	matched, err := regexp.Match(re, []byte(value))
	if err != nil {
		return false
	}
	return matched
}

type tag struct {
	value string
}

func (t tag) Value() string {
	return t.value
}

func (t tag) IsTag(re string) bool {
	return matches(re, t.value)
}

func (t tag) IsBranch(re string) bool {
	return false
}

type branch struct {
	value string
}

func (t branch) Value() string {
	return t.value
}

func (t branch) IsTag(re string) bool {
	return false
}

func (t branch) IsBranch(re string) bool {
	return matches(re, t.value)
}

func NewTag(value string) BuildVersioner {
	return tag{value}
}

func NewBranch(value string) BuildVersioner {
	return branch{value}
}
