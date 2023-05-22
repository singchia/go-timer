package scheduler

import (
	"fmt"
	"testing"
)

func ForEachFunction(node *CircularNode) error {
	fmt.Printf("%p, %v\n", node, node.data)
	return nil
}

func Test_AppendNode(t *testing.T) {
	circularList := NewCircularList()

	circularList.ForEach(ForEachFunction)
	fmt.Printf("**********************\n")

	node := circularList.AppendNode(1)
	if node == nil {
		t.Error("empty")
		return
	}
	circularList.ForEach(ForEachFunction)
	fmt.Printf("**********************\n")

	node = circularList.AppendNode(2)
	if node == nil {
		t.Error("empty")
		return
	}
	circularList.ForEach(ForEachFunction)
	fmt.Printf("**********************\n")

	node = circularList.AppendNode(3)
	if node == nil {
		t.Error("empty")
		return
	}
	circularList.ForEach(ForEachFunction)
	fmt.Printf("**********************\n")
}

func Test_DeleteNode_Empty(t *testing.T) {
	circularList := NewCircularList()
	ok := circularList.DeleteNode(&CircularNode{})
	if !ok {
		t.Log("node not exists")
		return
	}
	t.Error("delete empty node abnormal")
}

func Test_DeleteNode_Unempty(t *testing.T) {
	circularList := NewCircularList()

	circularList.ForEach(ForEachFunction)
	fmt.Printf("**********************\n")

	//add node1
	node1 := circularList.AppendNode(1)
	if node1 == nil {
		t.Error("empty")
		return
	}
	circularList.ForEach(ForEachFunction)
	fmt.Printf("**********************\n")

	//add node2
	node2 := circularList.AppendNode(2)
	if node2 == nil {
		t.Error("empty")
		return
	}
	circularList.ForEach(ForEachFunction)
	fmt.Printf("**********************\n")

	//delete node2
	ok := circularList.DeleteNode(node2)
	if !ok {
		t.Error("node not exists")
		return
	}
	circularList.ForEach(ForEachFunction)
	fmt.Printf("**********************\n")

	//add node3
	node3 := circularList.AppendNode(3)
	if node3 == nil {
		t.Error("empty")
		return
	}
	circularList.ForEach(ForEachFunction)
	fmt.Printf("**********************\n")

	//delete node3
	ok = circularList.DeleteNode(node3)
	if !ok {
		t.Error("node not exists")
		return
	}
	circularList.ForEach(ForEachFunction)
	fmt.Printf("**********************\n")

	//delete a unexist node
	ok = circularList.DeleteNode(&CircularNode{})
	if ok {
		t.Error("node should not exists")
		return
	}
	circularList.ForEach(ForEachFunction)
	fmt.Printf("**********************\n")

	//delete node1
	ok = circularList.DeleteNode(node1)
	if !ok {
		t.Error("node not exists")
		return
	}
	circularList.ForEach(ForEachFunction)
	fmt.Printf("**********************\n")
}

func Test_GetCurNode(t *testing.T) {
	circularList := NewCircularList()
	//get current node when list empty
	node := circularList.GetCurNodeWithNoCopied()
	if node != nil {
		t.Error("cur node abnormal exists")
		return
	}
	circularList.ForEach(ForEachFunction)
	fmt.Printf("**********************\n")

	//add node1
	node1 := circularList.AppendNode(1)
	if node1 == nil {
		t.Error("empty")
		return
	}
	circularList.ForEach(ForEachFunction)

	//get current node
	node = circularList.GetCurNodeWithNoCopied()
	if node == nil {
		t.Error("cur node abnormal exists")
		return
	}
	fmt.Printf("current node: %p, %v\n", node, node.data)
	fmt.Printf("**********************\n")

	//add node2
	node2 := circularList.AppendNode(2)
	if node2 == nil {
		t.Error("empty")
		return
	}
	circularList.ForEach(ForEachFunction)
	//get current node
	node = circularList.GetCurNodeWithNoCopied()
	if node == nil {
		t.Error("cur node not exists")
		return
	}
	fmt.Printf("current node: %p, %v\n", node, node.data)
	fmt.Printf("**********************\n")

	//delete node2
	ok := circularList.DeleteNode(node2)
	if !ok {
		t.Error("node not exists")
		return
	}
	circularList.ForEach(ForEachFunction)

	//get current node
	node = circularList.GetCurNodeWithNoCopied()
	if node == nil {
		t.Error("cur node not exists")
		return
	}
	fmt.Printf("current node: %p, %v\n", node, node.data)
	fmt.Printf("**********************\n")

	//delete node1
	ok = circularList.DeleteNode(node1)
	if !ok {
		t.Error("node not exists")
		return
	}
	circularList.ForEach(ForEachFunction)

	//get current node
	node = circularList.GetCurNode()
	if node != nil {
		t.Error("cur node abnormal exists")
		fmt.Printf("current node: %p, %v", node, node.data)
		return
	}
	fmt.Printf("**********************\n")
}

func Test_RightShiftCurPointer(t *testing.T) {
	circularList := NewCircularList()

	//add node1
	node1 := circularList.AppendNode(1)
	if node1 == nil {
		t.Error("empty")
		return
	}
	circularList.ForEach(ForEachFunction)
	fmt.Printf("**********************\n")

	//get current node
	node := circularList.GetCurNodeWithNoCopied()
	if node == nil {
		t.Error("cur node abnormal exists")
		return
	}
	fmt.Printf("current node: %p, %v\n", node, node.data)
	fmt.Printf("**********************\n")

	//add node2
	node2 := circularList.AppendNode(2)
	if node2 == nil {
		t.Error("empty")
		return
	}
	circularList.ForEach(ForEachFunction)
	//get current node
	node = circularList.GetCurNodeWithNoCopied()
	if node == nil {
		t.Error("cur node not exists")
		return
	}
	fmt.Printf("current node: %p, %v\n", node, node.data)
	fmt.Printf("**********************\n")

	//right shift the cur pointer
	ok := circularList.RightShiftCurPointer()
	if !ok {
		t.Error("rigth shift error")
		return
	}
	//get current node
	node = circularList.GetCurNodeWithNoCopied()
	if node == nil {
		t.Error("cur node not exists")
		return
	}
	fmt.Printf("current node: %p, %v\n", node, node.data)
	fmt.Printf("**********************\n")

	//right shift the cur pointer
	ok = circularList.RightShiftCurPointer()
	if !ok {
		t.Error("rigth shift error")
		return
	}
	//get current node
	node = circularList.GetCurNodeWithNoCopied()
	if node == nil {
		t.Error("cur node not exists")
		return
	}
	fmt.Printf("current node: %p, %v\n", node, node.data)
	fmt.Printf("**********************\n")
}

func Test_RightShiftCurPointerAndUpdate(t *testing.T) {
	circularList := NewCircularList()

	//add node1
	node1 := circularList.AppendNode(1)
	if node1 == nil {
		t.Error("empty")
		return
	}
	circularList.ForEach(ForEachFunction)
	fmt.Printf("**********************\n")

	//add node2
	node2 := circularList.AppendNode(2)
	if node2 == nil {
		t.Error("empty")
		return
	}
	circularList.ForEach(ForEachFunction)
	fmt.Printf("**********************\n")

	//right shift the cur pointer and update the data
	ok := circularList.RightShiftCurPointerAndUpdate(3)
	if !ok {
		t.Error("rigth shift error")
		return
	}
	//get current node
	node := circularList.GetCurNodeWithNoCopied()
	if node == nil {
		t.Error("cur node not exists")
		return
	}
	fmt.Printf("current node: %p, %v\n", node, node.data)
	fmt.Printf("**********************\n")
}

func Test_RightShiftCurPointerToCertainNode(t *testing.T) {
	circularList := NewCircularList()

	//add node1
	node1 := circularList.AppendNode(1)
	if node1 == nil {
		t.Error("empty")
		return
	}
	circularList.ForEach(ForEachFunction)
	fmt.Printf("**********************\n")

	//add node2
	node2 := circularList.AppendNode(2)
	if node2 == nil {
		t.Error("empty")
		return
	}
	circularList.ForEach(ForEachFunction)
	fmt.Printf("**********************\n")

	//add node3
	node3 := circularList.AppendNode(3)
	if node3 == nil {
		t.Error("empty")
		return
	}
	circularList.ForEach(ForEachFunction)
	fmt.Printf("**********************\n")

	//right shift the cur pointer
	ok := circularList.RightShiftCurPointerToCertainNode(node3)
	if !ok {
		t.Error("rigth shift error")
		return
	}
	//get current node
	node := circularList.GetCurNodeWithNoCopied()
	if node == nil {
		t.Error("cur node not exists")
		return
	}
	fmt.Printf("current node: %p, %v\n", node, node.data)
	fmt.Printf("**********************\n")
}
