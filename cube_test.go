package main

import (
	"encoding/json"
	"testing"

	"github.com/getlantern/deepcopy"
	"github.com/qmuntal/gltf"
	"github.com/rakyll/statik/fs"
)

func BenchmarkOpenGltf(b *testing.B) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		statikFS, err := fs.New()
		if err != nil {
			b.Fatal(err)
		}

		g, err := statikFS.Open("/cube.gltf")
		if err != nil {
			b.Fatal(err)
		}

		doc := new(gltf.Document)
		if err := gltf.NewDecoder(g).Decode(doc); err != nil {
			b.Fatal(err)
		}

		if len(doc.Nodes) != 26 {
			b.Fatalf("nodes count must be 26, actual: %d", len(doc.Nodes))
		}
	}
}

func BenchmarkOpenGlb(b *testing.B) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		statikFS, err := fs.New()
		if err != nil {
			b.Fatal(err)
		}

		g, err := statikFS.Open("/cube.glb")
		if err != nil {
			b.Fatal(err)
		}

		doc := new(gltf.Document)
		if err := gltf.NewDecoder(g).Decode(doc); err != nil {
			b.Fatal(err)
		}

		if len(doc.Nodes) != 26 {
			b.Fatalf("nodes count must be 26, actual: %d", len(doc.Nodes))
		}
	}
}

func BenchmarkOpenGltfUnlit(b *testing.B) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		statikFS, err := fs.New()
		if err != nil {
			b.Fatal(err)
		}

		g, err := statikFS.Open("/cube_unlit.gltf")
		if err != nil {
			b.Fatal(err)
		}

		doc := new(gltf.Document)
		if err := gltf.NewDecoder(g).Decode(doc); err != nil {
			b.Fatal(err)
		}

		if len(doc.Nodes) != 26 {
			b.Fatalf("nodes count must be 26, actual: %d", len(doc.Nodes))
		}
	}
}

func BenchmarkOpenGlbUnlit(b *testing.B) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		statikFS, err := fs.New()
		if err != nil {
			b.Fatal(err)
		}

		g, err := statikFS.Open("/cube_unlit.glb")
		if err != nil {
			b.Fatal(err)
		}

		doc := new(gltf.Document)
		if err := gltf.NewDecoder(g).Decode(doc); err != nil {
			b.Fatal(err)
		}

		if len(doc.Nodes) != 26 {
			b.Fatalf("nodes count must be 26, actual: %d", len(doc.Nodes))
		}
	}
}

func BenchmarkOpenOnceAndDeepCopy(b *testing.B) {
	statikFS, err := fs.New()
	if err != nil {
		b.Fatal(err)
	}

	g, err := statikFS.Open("/cube.gltf")
	if err != nil {
		b.Fatal(err)
	}

	doc := new(gltf.Document)
	if err := gltf.NewDecoder(g).Decode(doc); err != nil {
		b.Fatal(err)
	}

	if len(doc.Nodes) != 26 {
		b.Fatalf("nodes count must be 26, actual: %d", len(doc.Nodes))
	}

	copiedDoc := new(gltf.Document)
	if err := deepcopy.Copy(copiedDoc, doc); err != nil {
		b.Fatal(err)
	}

	before, err := json.Marshal(doc)
	if err != nil {
		b.Fatal(err)
	}
	after, err := json.Marshal(copiedDoc)
	if err != nil {
		b.Fatal(err)
	}
	if string(before) != string(after) {
		b.Fatal("copied doc is not equal to original doc")
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		copied := new(gltf.Document)
		if err := deepcopy.Copy(copied, doc); err != nil {
			b.Fatal(err)
		}

		if len(copied.Nodes) != 26 {
			b.Fatalf("nodes count must be 26, actual: %d", len(doc.Nodes))
		}
	}
}

func BenchmarkNodeToDefinition(b *testing.B) {
	statikFS, err := fs.New()
	if err != nil {
		b.Fatal(err)
	}

	gltfDefault, err := statikFS.Open("/cube.gltf")
	if err != nil {
		b.Fatal(err)
	}

	doc := new(gltf.Document)
	if err := gltf.NewDecoder(gltfDefault).Decode(doc); err != nil {
		b.Fatal(err)
	}

	if len(doc.Nodes) != 26 {
		b.Fatalf("nodes count must be 26, actual: %d", len(doc.Nodes))
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := nodeToDefinition(doc.Nodes)
		if err != nil {
			b.Fatal(err)
		}
	}
}
