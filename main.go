package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"
)

type Node struct {
	value int
	left  *Node
	right *Node
}

type Tree struct {
	node *Node
}

func (t *Tree) insert(i int) *Tree {
	if t.node == nil {
		t.node = &Node{value: i}
	} else {
		t.node.insert(i)
	}
	return t
}

func (n *Node) insert(i int) {
	if n.value >= i {
		if n.left == nil {
			n.left = &Node{value: i}
		} else {
			n.left.insert(i)
		}
	} else {
		if n.right == nil {
			n.right = &Node{value: i}
		} else {
			n.right.insert(i)
		}
	}
}

func printNode(n *Node) {
	if n == nil {
		return
	}

	fmt.Println("value", n.value)
	printNode(n.left)
	printNode(n.right)
}

func (n *Node) find(i, c int) (bool, int) {
	c += 1

	if n == nil {
		return false, c
	}

	if n.value == i {
		return true, c
	}

	if n.value >= i {
		return n.left.find(i, c)
	} else {
		return n.right.find(i, c)
	}
}

func writeChild(graph *cgraph.Graph, n *Node, nn *cgraph.Node) {
	if n.left != nil {
		ln, err := graph.CreateNode(fmt.Sprint(n.left.value))
		if err != nil {
			log.Fatal(err)
		}
		_, err = graph.CreateEdge("", nn, ln)
		if err != nil {
			log.Fatal(err)
		}

		writeChild(graph, n.left, ln)
	}

	if n.right != nil {
		rn, err := graph.CreateNode(fmt.Sprint(n.right.value))
		if err != nil {
			log.Fatal(err)
		}
		_, err = graph.CreateEdge("", nn, rn)
		if err != nil {
			log.Fatal(err)
		}

		writeChild(graph, n.right, rn)
	}

}

func writeGraph(t *Tree) {
	g := graphviz.New()
	graph, err := g.Graph()
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := graph.Close(); err != nil {
			log.Fatal(err)
		}
		g.Close()
	}()

	n := t.node

	nn, err := graph.CreateNode(fmt.Sprint(n.value))
	if err != nil {
		log.Fatal(err)
	}

	writeChild(graph, n, nn)

	var buf bytes.Buffer
	if err := g.Render(graph, graphviz.PNG, &buf); err != nil {
		log.Fatal(err)
	}

	// 2. get as image.Image instance
	_, err = g.RenderImage(graph)
	if err != nil {
		log.Fatal(err)
	}

	// 3. write to file directly
	if err := g.RenderFilename(graph, graphviz.PNG, "graph.png"); err != nil {
		log.Fatal(err)
	}

}

func main() {
	rand.Seed(time.Now().UnixNano())
	t := &Tree{}

	for i := 0; i < 100; i++ {
		j := rand.Intn(10000)
		// fmt.Printf("insert: %v\n", j)

		t.insert(j)
	}

	writeGraph(t)

	for {
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Printf("> ")
		scanner.Scan()

		s := scanner.Text()
		if s == "" {
			continue
		}

		i, err := strconv.Atoi(scanner.Text())
		if err != nil {
			fmt.Println("invalid number")
			continue
		}

		c := 0
		fmt.Println(t.node.find(i, c))
	}
}
