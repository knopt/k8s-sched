package main

import (
	v1 "k8s.io/api/core/v1"
	schedulerapi "k8s.io/kubernetes/pkg/scheduler/api"
)

type Filter struct {
	Name string
	Func func(node v1.Node, pod *v1.Pod) bool
}

func (f Filter) Handler(args schedulerapi.ExtenderArgs) *schedulerapi.ExtenderFilterResult {
	filteredNodes := make([]v1.Node, 0, len(args.Nodes.Items))
	filteredNodesNames := make([]string, 0, len(args.Nodes.Items))

	for _, node := range args.Nodes.Items {
		if f.Func(node, args.Pod) {
			filteredNodes = append(filteredNodes, node)
			filteredNodesNames = append(filteredNodesNames, node.Name)
		}
	}

	return &schedulerapi.ExtenderFilterResult{
		Nodes: &v1.NodeList{
			Items: filteredNodes,
		},
		NodeNames:   &filteredNodesNames,
		FailedNodes: nil,
		Error:       "",
	}
}
