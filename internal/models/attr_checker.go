package models

import "slices"

// AttrChecker This checker help checking if this attribute is validate for specific endpoint

type MemberMethod int

const (
	MemberUpdate MemberMethod = iota
	MemberQuery  MemberMethod = iota
)

type AttrChecker struct {
	memberMap map[MemberMethod]([]string)

	// TODO: Post Method
}

type hintResp struct {
	Hint          string   `json:"hint"`
	AvailableAttr []string `json:"available_attr"`
}

func NewAttrChecker() AttrChecker {
	memberUpdateCandidates := []string{
		"nick",
		"email",
		"power",
		"passwd",
		"is_delete",
	}
	memberQueryCandidates := []string{
		"id",
		"nick",
		"email",
		"power",
		"is_delete",
	}

	memberMap := map[MemberMethod]([]string){
		MemberUpdate: memberUpdateCandidates,
		MemberQuery:  memberQueryCandidates,
	}
	return AttrChecker{
		memberMap: memberMap,
	}
}

// MemberCheck Check if this attr in the candidates of this method
// return true if validate
// return false if not contains
func (ck AttrChecker) MemberCheck(method MemberMethod, attr string) (bool, hintResp) {
	hint := hintResp{
		Hint:          "the attribute must in available_attr list",
		AvailableAttr: ck.memberMap[method],
	}
	return slices.Contains(ck.memberMap[method], attr), hint
}

// MemberCandidates This function returns available attributes slices for specific method
func (ck AttrChecker) MemberCandidates(method MemberMethod) []string {
	return ck.memberMap[method]
}
