package main

import (
	"bytes"
	"fmt"
	"math"

	"github.com/getlantern/deepcopy"
	"github.com/qmuntal/gltf"
	"github.com/rakyll/statik/fs"
	"github.com/westphae/quaternion"
)

type Definition struct {
	UBL gltf.Node `json:"UBL"`
	UBM gltf.Node `json:"UBM"`
	UBR gltf.Node `json:"UBR"`
	UML gltf.Node `json:"UML"`
	UMM gltf.Node `json:"UMM"`
	UMR gltf.Node `json:"UMR"`
	UFL gltf.Node `json:"UFL"`
	UFM gltf.Node `json:"UFM"`
	UFR gltf.Node `json:"UFR"`
	MBL gltf.Node `json:"MBL"`
	MBM gltf.Node `json:"MBM"`
	MBR gltf.Node `json:"MBR"`
	MML gltf.Node `json:"MML"`
	MMR gltf.Node `json:"MMR"`
	MFL gltf.Node `json:"MFL"`
	MFM gltf.Node `json:"MFM"`
	MFR gltf.Node `json:"MFR"`
	DBL gltf.Node `json:"DBL"`
	DBM gltf.Node `json:"DBM"`
	DBR gltf.Node `json:"DBR"`
	DML gltf.Node `json:"DML"`
	DMM gltf.Node `json:"DMM"`
	DMR gltf.Node `json:"DMR"`
	DFL gltf.Node `json:"DFL"`
	DFM gltf.Node `json:"DFM"`
	DFR gltf.Node `json:"DFR"`
}

type Degree int

type Rotation [4]float64

const (
	rotateRightU Degree = iota
	rotateRightD
	rotateRightF
	rotateRightB
	rotateRightL
	rotateRightR
	rotateLeftU
	rotateLeftD
	rotateLeftF
	rotateLeftB
	rotateLeftL
	rotateLeftR
	rotateU2
	rotateD2
	rotateF2
	rotateB2
	rotateL2
	rotateR2
)

const (
	degree90  = math.Pi / 2
	degree180 = math.Pi
)

var (
	quaternionZero = quaternion.FromEuler(0, 0, 0)
	quaternionU    = quaternion.FromEuler(0, degree90, 0)
	quaternionD    = quaternion.FromEuler(0, -degree90, 0)
	quaternionF    = quaternion.FromEuler(-degree90, 0, 0)
	quaternionB    = quaternion.FromEuler(degree90, 0, 0)
	quaternionL    = quaternion.FromEuler(0, 0, -degree90)
	quaternionR    = quaternion.FromEuler(0, 0, degree90)
	quaternionU2   = quaternion.FromEuler(0, degree180, 0)
	quaternionF2   = quaternion.FromEuler(-degree180, 0, 0)
	quaternionL2   = quaternion.FromEuler(0, 0, -degree180)
)

var gltfDoc = new(gltf.Document)

func initCube() error {
	sFS, err := fs.New()
	if err != nil {
		return err
	}

	gltfFile, err := sFS.Open("/cube.gltf")
	if err != nil {
		return err
	}
	if err := gltf.NewDecoder(gltfFile).Decode(gltfDoc); err != nil {
		return err
	}

	return nil
}

func generateCube(algorithm []string) ([]byte, error) {
	var (
		doc gltf.Document
		err error
		def *Definition
	)

	if err = deepcopy.Copy(&doc, gltfDoc); err != nil {
		return nil, err
	}

	def, err = nodeToDefinition(doc.Nodes)
	if err != nil {
		return nil, err
	}

	degrees, err := parseAlg(algorithm)
	if err != nil {
		return nil, err
	}

	for _, d := range degrees {
		rotate(def, d)
	}

	doc.Nodes = []*gltf.Node{
		&def.UBL, &def.UBM, &def.UBR,
		&def.UML, &def.UMM, &def.UMR,
		&def.UFL, &def.UFM, &def.UFR,
		&def.MBL, &def.MBM, &def.MBR,
		&def.MML, &def.MMR,
		&def.MFL, &def.MFM, &def.MFR,
		&def.DBL, &def.DBM, &def.DBR,
		&def.DML, &def.DMM, &def.DMR,
		&def.DFL, &def.DFM, &def.DFR,
	}

	buffer := new(bytes.Buffer)
	e := gltf.NewEncoder(buffer)
	if err := e.Encode(&doc); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

// nodeToDefinition
func nodeToDefinition(nodes []*gltf.Node) (*Definition, error) {
	def := new(Definition)

	if len(nodes) != 26 {
		return nil, fmt.Errorf("insufficient number of nodes")
	}

	for _, node := range nodes {
		switch node.Name {
		case "UBL":
			def.UBL = *node
		case "UBM":
			def.UBM = *node
		case "UBR":
			def.UBR = *node
		case "UML":
			def.UML = *node
		case "UMM":
			def.UMM = *node
		case "UMR":
			def.UMR = *node
		case "UFL":
			def.UFL = *node
		case "UFM":
			def.UFM = *node
		case "UFR":
			def.UFR = *node
		case "MBL":
			def.MBL = *node
		case "MBM":
			def.MBM = *node
		case "MBR":
			def.MBR = *node
		case "MML":
			def.MML = *node
		case "MMR":
			def.MMR = *node
		case "MFL":
			def.MFL = *node
		case "MFM":
			def.MFM = *node
		case "MFR":
			def.MFR = *node
		case "DBL":
			def.DBL = *node
		case "DBM":
			def.DBM = *node
		case "DBR":
			def.DBR = *node
		case "DML":
			def.DML = *node
		case "DMM":
			def.DMM = *node
		case "DMR":
			def.DMR = *node
		case "DFL":
			def.DFL = *node
		case "DFM":
			def.DFM = *node
		case "DFR":
			def.DFR = *node
		default:
			return nil, fmt.Errorf("unexpected node name: %s", node.Name)
		}
	}
	return def, nil
}

// parseAlg 回転記号の文字列を Degree のスライスに変換する
func parseAlg(alg []string) ([]Degree, error) {
	var degree []Degree
	for _, a := range alg {
		switch a {
		case "U":
			degree = append(degree, rotateRightU)
		case "D":
			degree = append(degree, rotateRightD)
		case "F":
			degree = append(degree, rotateRightF)
		case "B":
			degree = append(degree, rotateRightB)
		case "L":
			degree = append(degree, rotateRightL)
		case "R":
			degree = append(degree, rotateRightR)
		case "U2":
			degree = append(degree, rotateU2)
		case "D2":
			degree = append(degree, rotateD2)
		case "F2":
			degree = append(degree, rotateF2)
		case "B2":
			degree = append(degree, rotateB2)
		case "L2":
			degree = append(degree, rotateL2)
		case "R2":
			degree = append(degree, rotateR2)
		case "U'":
			degree = append(degree, rotateLeftU)
		case "D'":
			degree = append(degree, rotateLeftD)
		case "F'":
			degree = append(degree, rotateLeftF)
		case "B'":
			degree = append(degree, rotateLeftB)
		case "L'":
			degree = append(degree, rotateLeftL)
		case "R'":
			degree = append(degree, rotateLeftR)
		default:
			return nil, fmt.Errorf("unknown alg: %s", a)
		}
	}

	return degree, nil
}

// rotate 回転する面に合わせて、対象の gltf.Node に回転処理を行う
func rotate(def *Definition, degree Degree) {
	switch degree {
	case rotateRightU, rotateU2, rotateLeftU:
		def.UBL, def.UBM, def.UBR,
			def.UML, def.UMM, def.UMR,
			def.UFL, def.UFM, def.UFR =
			move(degree,
				def.UBL, def.UBM, def.UBR,
				def.UML, def.UMM, def.UMR,
				def.UFL, def.UFM, def.UFR,
			)

	case rotateRightD, rotateD2, rotateLeftD:
		def.DFL, def.DFM, def.DFR,
			def.DML, def.DMM, def.DMR,
			def.DBL, def.DBM, def.DBR =
			move(degree,
				def.DFL, def.DFM, def.DFR,
				def.DML, def.DMM, def.DMR,
				def.DBL, def.DBM, def.DBR,
			)

	case rotateRightF, rotateF2, rotateLeftF:
		def.UFL, def.UFM, def.UFR,
			def.MFL, def.MFM, def.MFR,
			def.DFL, def.DFM, def.DFR =
			move(degree,
				def.UFL, def.UFM, def.UFR,
				def.MFL, def.MFM, def.MFR,
				def.DFL, def.DFM, def.DFR,
			)

	case rotateRightB, rotateB2, rotateLeftB:
		def.UBR, def.UBM, def.UBL,
			def.MBR, def.MBM, def.MBL,
			def.DBR, def.DBM, def.DBL =
			move(degree,
				def.UBR, def.UBM, def.UBL,
				def.MBR, def.MBM, def.MBL,
				def.DBR, def.DBM, def.DBL,
			)

	case rotateRightL, rotateL2, rotateLeftL:
		def.UBL, def.UML, def.UFL,
			def.MBL, def.MML, def.MFL,
			def.DBL, def.DML, def.DFL =
			move(degree,
				def.UBL, def.UML, def.UFL,
				def.MBL, def.MML, def.MFL,
				def.DBL, def.DML, def.DFL,
			)

	case rotateRightR, rotateR2, rotateLeftR:
		def.UFR, def.UMR, def.UBR,
			def.MFR, def.MMR, def.MBR,
			def.DFR, def.DMR, def.DBR =
			move(degree,
				def.UFR, def.UMR, def.UBR,
				def.MFR, def.MMR, def.MBR,
				def.DFR, def.DMR, def.DBR,
			)
	}
}

// move 回転処理を行い、回転後の Rotation と gltf.Node の位置を返却する
func move(degree Degree,
	n1 gltf.Node, n2 gltf.Node, n3 gltf.Node,
	n4 gltf.Node, n5 gltf.Node, n6 gltf.Node,
	n7 gltf.Node, n8 gltf.Node, n9 gltf.Node,
) (
	gltf.Node, gltf.Node, gltf.Node,
	gltf.Node, gltf.Node, gltf.Node,
	gltf.Node, gltf.Node, gltf.Node,
) {
	n1.Rotation = prod(degree, n1.RotationOrDefault())
	n2.Rotation = prod(degree, n2.RotationOrDefault())
	n3.Rotation = prod(degree, n3.RotationOrDefault())
	n4.Rotation = prod(degree, n4.RotationOrDefault())
	n5.Rotation = prod(degree, n5.RotationOrDefault())
	n6.Rotation = prod(degree, n6.RotationOrDefault())
	n7.Rotation = prod(degree, n7.RotationOrDefault())
	n8.Rotation = prod(degree, n8.RotationOrDefault())
	n9.Rotation = prod(degree, n9.RotationOrDefault())

	switch degree {
	case rotateRightU, rotateRightD, rotateRightF, rotateRightB, rotateRightL, rotateRightR:
		return n7, n4, n1,
			n8, n5, n2,
			n9, n6, n3

	case rotateU2, rotateD2, rotateF2, rotateB2, rotateL2, rotateR2:
		return n9, n8, n7,
			n6, n5, n4,
			n3, n2, n1

	case rotateLeftU, rotateLeftD, rotateLeftF, rotateLeftB, rotateLeftL, rotateLeftR:
		return n3, n6, n9,
			n2, n5, n8,
			n1, n4, n7

	default:
		return n1, n2, n3,
			n4, n5, n6,
			n7, n8, n9
	}
}

// prod Rotation と回転の積を求め Rotation として取得する
func prod(degree Degree, rotation Rotation) Rotation {
	return quaternionToRotation(
		quaternion.Prod(
			rotationToQuaternion(rotation),
			getRotateQuaternion(degree),
		),
	)
}

// rotationToQuaternion glTF 内の Rotation を quaternion.Quaternion に変換する
func rotationToQuaternion(rotation Rotation) quaternion.Quaternion {
	return quaternion.New(
		rotation[0], // W
		rotation[1], // X
		rotation[2], // Y
		rotation[3], // Z
	)
}

// quaternionToRotation quaternion.Quaternion を glTF 内の Rotation に変換する
func quaternionToRotation(q quaternion.Quaternion) Rotation {
	return Rotation{q.W, q.X, q.Y, q.Z}
}

// getRotateQuaternion Degree で表す回転方向・回転量を quaternion.Quaternion に変換して取得する
func getRotateQuaternion(degree Degree) quaternion.Quaternion {
	switch degree {
	case rotateRightU, rotateLeftD:
		return quaternionU
	case rotateRightD, rotateLeftU:
		return quaternionD
	case rotateRightF, rotateLeftB:
		return quaternionF
	case rotateRightB, rotateLeftF:
		return quaternionB
	case rotateRightL, rotateLeftR:
		return quaternionL
	case rotateRightR, rotateLeftL:
		return quaternionR
	case rotateU2, rotateD2:
		return quaternionU2
	case rotateF2, rotateB2:
		return quaternionF2
	case rotateL2, rotateR2:
		return quaternionL2
	default:
		return quaternionZero
	}
}
