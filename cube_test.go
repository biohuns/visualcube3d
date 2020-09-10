package main

import (
	"testing"

	"github.com/getlantern/deepcopy"
	"github.com/qmuntal/gltf"
	"github.com/rakyll/statik/fs"
)

func BenchmarkOpenGltf(b *testing.B) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		sFS, err := fs.New()
		if err != nil {
			b.Fatal(err)
		}

		g, err := sFS.Open("/cube.gltf")
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
	sFS, err := fs.New()
	if err != nil {
		b.Fatal(err)
	}

	g, err := sFS.Open("/cube.gltf")
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
	sFS, err := fs.New()
	if err != nil {
		b.Fatal(err)
	}

	gltfDefault, err := sFS.Open("/cube.gltf")
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
