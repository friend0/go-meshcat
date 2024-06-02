package internal

type GeometryMeta struct {
	item_size int
	_type     string
	array     []float32
}

type GeometryAttributes struct {
	position GeometryMeta
	normal   GeometryMeta
	uv       GeometryMeta
}

type BoundingSphere struct {
	center []int
	radius float32
}

type GeometryData struct {
	attributes      GeometryAttributes
	bounding_sphere BoundingSphere
}

type Geometry struct {
	*GeometryMeta
	data GeometryData
}

type Box struct {
}
