package main

import (
	"testing"

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
