package main

import (
	"os"
	"bufio"
	"fmt"
	"runtime"
)

// at a point on the maze, can we go north, east, south, or west
type node struct {
	id			int
	n, e, s, w	*node
}

// the entire maze is a bunch of nodes glued together
type maze struct {
	// we need to keep hold of the start and finish in particular
	start, finish	*node

	// the complete array of maze nodes for gluing things together
	nodes		[][]*node

	// maze width, length, and distance along current line, for building
	w, l, p		int

	// number of nodes
	count		int
}

type walker struct {
	id		int
	visited	map[*node]bool
	path	string
}
var WID int = 0

func main() {
	runtime.GOMAXPROCS(2)
	if len(os.Args) != 2 {
		fmt.Println("mazes <file>, fool")
		return
	}
	mazes, err := mazeReader(os.Args[1])
	if err != nil {
		fmt.Printf("maze error: %s\n", err)
		return
	}
	for m := range mazes {
//		fmt.Printf("width: %d; length: %d; nodes: %d\n", m.w, m.l, m.count)

		w := newWalker()
		ch := make(chan *walker)
		go w.walk(m, m.start, ch)
		<-ch
//		for res := range ch {
//			fmt.Printf("walker %d found path '%s'\n", res.id, res.path)
//		}
	}
}

func mazeReader(path string) (chan *maze, os.Error) {
	file, err := os.Open(path, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}

	in := make(chan byte, 1)
	out := make(chan *maze, 1)
	br := bufio.NewReader(file)

	go func() {
		for {
			s, err := br.ReadString('\n')
			if err != nil {
				close(in)
				break
			}
			if s == "\n" {
				in <- '\xff'
				continue
			}
			for i := 0; i < len(s); i++ {
				in <- s[i]
			}
		}
		file.Close()
	}()
	go readMazes(in, out)
	return out, nil
}

func readMazes(in chan byte, out chan *maze) {
	for {
		m := newMaze()

		// read in first line to find width of maze; it will be all walls.
		for b := <-in; b == '#'; b = <-in {
			m.w++
		}

		// now we know the width, fill in the first line with "walls"
		// if *node == nil, it's a wall, conveniently
		m.addLine()
		// and add a second line to start filling with nodes, as the \n
		// of the first line has already been pulled from the channel
		// and thus will not trigger the switch...
		m.addLine()

		for b := <-in; b != '\xff' && !closed(in); b = <-in {
			switch b {
			case '\n': m.addLine()
			case '#':  m.addWall()
			case '.':  m.addNode()
			case 's':  m.start  = m.addNode()
			case 'f':  m.finish = m.addNode()
			}
		}
		if closed(in) {
			close(out)
			break
		}
		out <- m
	}
}

func newMaze() *maze {
	m := new(maze)
	m.nodes = make([][]*node, 0, 32)
	return m
}

func (m *maze) addLine() {
	if len(m.nodes) == cap(m.nodes) {
		// shouldn't hit a 512-line maze, so be lazy ;-)
		fmt.Println("OH GOD NOES!")
		os.Exit(1)
	}
	m.l = len(m.nodes)
	m.nodes = m.nodes[0:m.l+1]
	m.nodes[m.l] = make([]*node, m.w)
	m.p = 0
}

func (m *maze) addWall() {
	m.p++
}

func (m *maze) addNode() *node {
	m.count++
	n := &node{id: m.count}
	if v := m.nodes[m.l][m.p-1]; v != nil {
		n.w = v
		v.e = n
	}
	if v := m.nodes[m.l-1][m.p]; v != nil {
		n.n = v
		v.s = n
	}
	m.nodes[m.l][m.p] = n
	m.p++
	return n
}

func newWalker() *walker {
	w := new(walker)
	WID++
	w.id = WID
	w.visited = make(map[*node]bool)
	return w
}

func (w *walker) clone() *walker {
	c := newWalker();
	// uncomment this and remove close(ch) from walk() to get
	// the full set of possible paths through the maze
	// warning, it's a bit memory intensive ;p
//	for k,v := range w.visited {
//		c.visited[k] = v
//	}
	c.visited = w.visited
	c.path = w.path
	return c
}

func (w *walker) walk(m *maze, n *node, ch chan *walker) {
	if closed(ch) {
		return
	}
	if n == m.finish {
		ch <- w
		close(ch)
		return
	}
	w.visited[n] = true

	if _, ok := w.visited[n.n]; n.n != nil && !ok {
		c := w.clone()
		c.path += "n"
		go c.walk(m, n.n, ch)
	}
	if _, ok := w.visited[n.e]; n.e != nil && !ok {
		c := w.clone()
		c.path += "e"
		go c.walk(m, n.e, ch)
	}
	if _, ok := w.visited[n.s]; n.s != nil && !ok {
		c := w.clone()
		c.path += "s"
		go c.walk(m, n.s, ch)
	}
	if _, ok := w.visited[n.w]; n.w != nil && !ok {
		c := w.clone()
		c.path += "w"
		go c.walk(m, n.w, ch)
	}
}
