package internal

import (
	"encoding/json"
	"unsafe"
)

type TopFunc func() unsafe.Pointer

type Node struct {
	Index int64
	Value any
	next  unsafe.Pointer // Pointer to *Node
	down  unsafe.Pointer // Pointer to down level *Node
}

func NewNode(index int64, next unsafe.Pointer, value any) *Node {
	return &Node{Index: index, next: next, Value: value}
}

func (node *Node) String() string {
	b, _ := json.Marshal(node)
	return string(b)
}

// Hi @ep-us
//
//I try to get access to github from cdjenkins to https://github.freewheel.tv/ops/ui-terraform.git
//
//PS
//jenkins-github, "svc-fw-jenkins-sh PAT", github-api-token, svc-fw-jenkins-sh:
//Password authentication is not available for Git operations;  You must use a personal access token or SSH key.
