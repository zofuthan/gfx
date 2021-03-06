// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gfx

import (
	"sync"

	"azul3d.org/lmath.v1"
)

// Boundable represents any object that can return it's axis-aligned bounding
// box.
type Boundable interface {
	// Bounds returns the axis-aligned bounding box of this boundable object.
	Bounds() lmath.Rect3
}

// Bounds is a simple datatype which implements the Boundable interface.
type Bounds lmath.Rect3

// Bounds implements the Boundable interface.
func (b Bounds) Bounds() lmath.Rect3 {
	return lmath.Rect3(b)
}

// TexCoordSet represents a single texture coordinate set for a mesh.
type TexCoordSet struct {
	// The slice of texture coordinates for the set.
	Slice []TexCoord

	// Weather or not the texture coordinates of this set have changed since
	// the last time the mesh was loaded. If set to true the renderer should
	// take note and re-upload the data slice to the graphics hardware.
	Changed bool
}

// VertexAttrib represents a per-vertex attribute.
type VertexAttrib struct {
	// The literal per-vertex data slice. It must be a slice whose length is
	// exactly the same as the mesh's Vertices slice (because it is literally
	// per-vertex data). The underlying type must be one of the following or
	// else the attribute may be ignored completely:
	//  []float32
	//  [][]float32
	//  []gfx.Vec3
	//  [][]gfx.Vec3
	//  []gfx.Mat4
	//  [][]gfx.Mat4
	//  []gfx.Vec4
	//  [][]gfx.Vec4
	Data interface{}

	// Weather or not the per-vertex data (see the Data field) has changed
	// since the last time the mesh was loaded. If set to true the renderer
	// should take note and re-upload the data slice to the graphics hardware.
	Changed bool
}

// Copy returns a new copy of this vertex attribute data set. It makes a deep
// copy of the underlying Data slice. Explicitly not copied is the Changed
// boolean.
func (a VertexAttrib) Copy() VertexAttrib {
	var cpy interface{}
	switch t := a.Data.(type) {
	case []float32:
		c := make([]float32, len(t))
		copy(c, t)
		cpy = c

	case []Vec3:
		c := make([]Vec3, len(t))
		copy(c, t)
		cpy = c

	case []Vec4:
		c := make([]Vec4, len(t))
		copy(c, t)
		cpy = c

	case []Mat4:
		c := make([]Mat4, len(t))
		copy(c, t)
		cpy = c

	case [][]float32:
		c := make([][]float32, len(t))
		for i, s := range t {
			c[i] = make([]float32, len(s))
			copy(c[i], t[i])
		}

	case [][]Vec3:
		c := make([][]Vec3, len(t))
		for i, s := range t {
			c[i] = make([]Vec3, len(s))
			copy(c[i], t[i])
		}

	case [][]Vec4:
		c := make([][]Vec4, len(t))
		for i, s := range t {
			c[i] = make([]Vec4, len(s))
			copy(c[i], t[i])
		}

	case [][]Mat4:
		c := make([][]Mat4, len(t))
		for i, s := range t {
			c[i] = make([]Mat4, len(s))
			copy(c[i], t[i])
		}

	default:
		return VertexAttrib{}
	}
	return VertexAttrib{Data: cpy}
}

// NativeMesh represents the native object of a mesh, typically only renderers
// create these.
type NativeMesh Destroyable

// Mesh represents a single mesh made up of several components. A mesh may or
// may not be made up of indexed vertices, etc, depending on whether or not
// len(m.Indices) == 0 holds true.
// In the event that a mesh is indexed, m.Indices holds the indices and it can
// be expected that each other slice (Vertices for instance) will hold at least
// enough elements (or be nil) such that the each index will not be out of
// bounds.
//
// Clients are responsible for utilizing the RWMutex of the mesh when using it
// or invoking methods.
type Mesh struct {
	sync.RWMutex

	// The native object of this mesh. Once loaded the renderer using this mesh
	// must assign a value to this field. Typically clients should not assign
	// values to this field at all.
	NativeMesh

	// Weather or not this mesh is currently loaded or not.
	Loaded bool

	// If true then when this mesh is loaded the sources of it will be kept
	// instead of being set to nil (which allows them to be garbage collected).
	KeepDataOnLoad bool

	// Weather or not the mesh will be dynamically updated. Only used as a hint
	// to increase performence of dynamically updated meshes, does not actually
	// control whether or not a mesh may be dynamically updated.
	Dynamic bool

	// AABB is the axis aligned bounding box of this mesh. There may not be one
	// if AABB.Empty() == true, but one can be calculate using the
	// CalculateBounds() method.
	AABB lmath.Rect3

	// A slice of indices, if non-nil then this slice contains indices into
	// each other slice (such as Vertices) and this is a indexed mesh.
	// The indices are uint32 (instead of int) for compatability with graphics
	// hardware.
	Indices []uint32

	// Weather or not the indices have changed since the last time the mesh
	// was loaded. If set to true the renderer should take note and
	// re-upload the data slice to the graphics hardware.
	IndicesChanged bool

	// The slice of vertices for the mesh.
	Vertices []Vec3

	// Weather or not the vertices have changed since the last time the
	// mesh was loaded. If set to true the renderer should take note and
	// re-upload the data slice to the graphics hardware.
	VerticesChanged bool

	// The slice of vertex colors for the mesh.
	Colors []Color

	// Weather or not the vertex colors have changed since the last time
	// the mesh was loaded. If set to true the renderer should take note
	// and re-upload the data slice to the graphics hardware.
	ColorsChanged bool

	// A slice of barycentric coordinates for the mesh.
	Bary []Vec3

	// Whether or not the barycentric coordinates have changed since the last
	// time the mesh was loaded. If set to true the renderer should take note
	// and re-upload the data slice to the graphics hardware.
	BaryChanged bool

	// A slice of texture coordinate sets for the mesh, there may be
	// multiple sets which directly relate to multiple textures on a
	// object.
	TexCoords []TexCoordSet

	// A map of custom per-vertex attributes for the mesh. It is analogous to
	// the Colors, Bary, and TexCoords fields. It allows you to submit a set of
	// named custom per-vertex data to shaders.
	//
	// For instance you could submit a set of per-vertex vec3's with:
	//  myData := make([]gfx.Vec3, len(mesh.Vertices))
	//  mesh.Attribs["MyName"] = gfx.VertexAttrib{
	//      Data: myData,
	//  }
	//
	// If changes to the data are made, the data set will have to be uploaded
	// to the graphics hardware again, so you must inform the renderer when you
	// change the data:
	//  ... modify myData ...
	//  mesh.Attribs["MyName"].Changed = true
	//
	// In GLSL you could access that per-vertex data by writing:
	//  attribute vec3 MyName;
	//
	// Arrays of data are available in GLSL by slice indice suffixes:
	//  // Data declared in Go:
	//  myData := make([][]gfx.Mat4, 2)
	//
	//  // And in GLSL:
	//  attribute mat4 MyName0; // Per-vertex data from myData[0].
	//  attribute mat4 MyName1; // Per-vertex data from myData[1].
	//
	// See the documentation on the VertexAttrib type for more information
	// regarding what data types may be used.
	Attribs map[string]VertexAttrib
}

// Copy returns a new copy of this Mesh. Depending on how large the mesh is
// this may be an expensive operation. Explicitly not copied over is the native
// mesh, the OnLoad slice, and the loaded and changed statuses (Loaded,
// IndicesChanged, VerticesChanged, etc).
//
// The mesh's read lock must be held for this method to operate safely.
func (m *Mesh) Copy() *Mesh {
	cpy := &Mesh{
		sync.RWMutex{},
		nil,   // Native mesh -- not copied.
		false, // Loaded status -- not copied.
		m.KeepDataOnLoad,
		m.Dynamic,
		m.AABB,
		make([]uint32, len(m.Indices)),
		false, // IndicesChanged -- not copied.
		make([]Vec3, len(m.Vertices)),
		false, // VerticesChanged -- not copied.
		make([]Color, len(m.Colors)),
		false, // ColorsChanged -- not copied.
		make([]Vec3, len(m.Bary)),
		false, // BaryChanged -- not copied.
		make([]TexCoordSet, len(m.TexCoords)),
		make(map[string]VertexAttrib, len(m.Attribs)),
	}

	copy(cpy.Indices, m.Indices)
	copy(cpy.Vertices, m.Vertices)
	copy(cpy.Colors, m.Colors)
	copy(cpy.Bary, m.Bary)
	for index, set := range m.TexCoords {
		setCpy := TexCoordSet{
			Slice: make([]TexCoord, len(set.Slice)),
		}
		copy(setCpy.Slice, set.Slice)
		cpy.TexCoords[index] = setCpy
	}
	for name, attrib := range m.Attribs {
		cpy.Attribs[name] = attrib.Copy()
	}
	return cpy
}

// Bounds implements the Boundable interface. It is thread-safe and performs
// locking automatically. If the AABB of this mesh is empty then the bounds are
// calculated.
func (m *Mesh) Bounds() lmath.Rect3 {
	m.Lock()
	if m.AABB.Empty() {
		m.CalculateBounds()
	}
	bounds := m.AABB
	m.Unlock()
	return bounds
}

// GenerateBary generates the barycentric coordinates for this mesh.
//
// The mesh's write lock must be held for this method to operate safely.
func (m *Mesh) GenerateBary() {
	var (
		bci = -1
		v   Vec3
	)
	for _ = range m.Vertices {
		// Add barycentric coordinates.
		bci++
		switch bci % 3 {
		case 0:
			v = Vec3{1, 0, 0}
		case 1:
			v = Vec3{0, 1, 0}
		case 2:
			v = Vec3{0, 0, 1}
		}
		m.Bary = append(m.Bary, v)
	}
}

// CalculateBounds calculates a new axis aligned bounding box for this mesh.
//
// The mesh's write lock must be held for this method to operate safely.
func (m *Mesh) CalculateBounds() {
	var bb lmath.Rect3
	if len(m.Vertices) > 0 {
		for _, v32 := range m.Vertices {
			v := v32.Vec3()
			bb.Min = bb.Min.Min(v)
			bb.Max = bb.Max.Max(v)
		}
	}
	m.AABB = bb
}

// HasChanged tells if any of the data slices of the mesh are marked as having
// changed.
//
// The mesh's read lock must be held for this method to operate safely.
func (m *Mesh) HasChanged() bool {
	if m.IndicesChanged || m.VerticesChanged || m.ColorsChanged || m.BaryChanged {
		return true
	}
	for _, texCoordSet := range m.TexCoords {
		if texCoordSet.Changed {
			return true
		}
	}
	for _, attrib := range m.Attribs {
		if attrib.Changed {
			return true
		}
	}
	return false
}

// ClearData sets the data slices of this mesh to nil if m.KeepDataOnLoad is
// set to false.
//
// The mesh's write lock must be held for this method to operate safely.
func (m *Mesh) ClearData() {
	if !m.KeepDataOnLoad {
		m.Indices = nil
		m.Vertices = nil
		m.Colors = nil
		m.Bary = nil
		m.TexCoords = nil
		m.Attribs = nil
	}
}

// Reset resets this mesh to it's default (NewMesh) state.
//
// The mesh's write lock must be held for this method to operate safely.
func (m *Mesh) Reset() {
	m.NativeMesh = nil
	m.Loaded = false
	m.KeepDataOnLoad = false
	m.Dynamic = false
	m.AABB = lmath.Rect3Zero
	m.Indices = m.Indices[:0]
	m.IndicesChanged = false
	m.Vertices = m.Vertices[:0]
	m.VerticesChanged = false
	m.Colors = m.Colors[:0]
	m.ColorsChanged = false
	m.Bary = m.Bary[:0]
	m.BaryChanged = false
	for _, tcs := range m.TexCoords {
		tcs.Slice = nil
		tcs.Changed = false
	}
	m.TexCoords = m.TexCoords[:0]
	m.Attribs = make(map[string]VertexAttrib)
}

// Destroy destroys this mesh for use by other callees to NewMesh. You must not
// use it after calling this method. This makes an implicit call to
// m.NativeMesh.Destroy.
//
// The mesh's write lock must be held for this method to operate safely.
func (m *Mesh) Destroy() {
	if m.NativeMesh != nil {
		m.NativeMesh.Destroy()
	}
	m.Reset()
	meshPool.Put(m)
}

var meshPool = sync.Pool{
	New: func() interface{} {
		return &Mesh{
			Attribs: make(map[string]VertexAttrib),
		}
	},
}

// NewMesh returns a new *Mesh, for effeciency it may be a re-used one (see the
// Destroy method) whose slices have zero-lengths.
func NewMesh() *Mesh {
	return meshPool.Get().(*Mesh)
}
