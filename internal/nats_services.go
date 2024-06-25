package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"

	"github.com/friend0/transformations"
	"github.com/nats-io/nats.go"
	"github.com/vmihailenco/msgpack/v5"
	"gonum.org/v1/gonum/mat"
)

func (s *Server) NATSSubscriptions() error {
	// todo: track subscriptions so we can cleanup
	_, err := s.urlSubscription()
	if err != nil {
		return err
	}

	// todo: not able to write objects to the browser directly via websocket,
	// right now we tell the browser to fetch the file from the server,
	// and for some reason that works. Suspect issue with serialization/deserialization pipeline from nats -> msgpack -> ws
	// somewhat lower priority for now, given that we're still able to load objects. Will be much more flexible if we can load the objects directly.

	// Manage requests to add mesh objects
	_, err = s.setObjectSubscription()
	if err != nil {
		return err
	}

	// Add stock geometry objects, like boxes, spheres, etc.
	_, err = s.setGeometrySubscription()
	if err != nil {
		return err
	}

	_, err = s.setTransformationSubscription()
	if err != nil {
		return err
	}

	_, err = s.missionSubscription()
	if err != nil {
		return err
	}

	_, err = s.delete()
	if err != nil {
		return err
	}

	s.NATS.Flush()
	log.Printf("Listening on [%s]", "meshcat")
	return nil
}

var MeshcatCommands = map[string]bool{
	"set_transform": true,
	"set_object":    true,
	"delete":        true,
	"set_property":  true,
	"set_animation": true,
}

func (s *Server) urlSubscription() (*nats.Subscription, error) {
	sub, err := s.NATS.QueueSubscribe("meshcat.url", "MESHCAT_URL_Q", func(msg *nats.Msg) {
		b, err := msgpack.Marshal(&msg)
		if err != nil {
			s.Logger.Error(fmt.Sprintf("error encoding message: %v", err))
		}
		err = s.Hub.Write(b)
		if err != nil {
			s.Logger.Error(fmt.Sprintf("error writing to web socket %v", err))
		}
	})
	if err != nil {
		s.Logger.Error(fmt.Sprintf("error creating NATS subscription: %v", err))
	}
	return sub, err
}

type SetFromServerMetadata struct {
	ResourceName string  `msgpack:"resource_name"`
	Path         string  `msgpack:"path"`
	PositionX    float64 `msgpack:"x"`
	PositionY    float64 `msgpack:"y"`
	PositionZ    float64 `msgpack:"z"`
}

type SetFromServer struct {
	Command `msgpack:"command"`
	Object  SetFromServerMetadata `msgpack:"object"`
}

// SetObjectSubscription handler
func (s *Server) setObjectSubscription() (*nats.Subscription, error) {
	sub, err := s.NATS.Subscribe("meshcat.objects", func(msg *nats.Msg) {
		s.Logger.Info(fmt.Sprintf("Received meshcat message from NATS `%s` on subject `%s`", string(msg.Data), strings.Split(msg.Subject, ".")[2:]))
		path := strings.Join(strings.Split(string(msg.Subject), ".")[2:], "/")
		log.Printf("Received meshcat message from NATS: %s", string(msg.Data))

		// let's say for now the message has the form "object_name path positionx positiony positionz"
		cmd := strings.Split(string(msg.Data), " ")
		object_name, path, x, y, z := cmd[0], cmd[1], cmd[2], cmd[3], cmd[4]
		fx, fy, fz, err := ParseFloats(x, y, z)
		if err != nil {
			s.Logger.Info(fmt.Sprintf("error processing position input in the command message %s", msg.Data))
			return
		}

		var buf bytes.Buffer
		enc := msgpack.NewEncoder(&buf)
		err = enc.Encode(SetFromServer{
			Object: SetFromServerMetadata{
				ResourceName: object_name,
				Path:         path,
				PositionX:    fx,
				PositionY:    fy,
				PositionZ:    fz,
			},
			Command: Command{
				Type: "set_object_from_server",
				Path: path,
			},
		})
		if err != nil {
			s.Logger.Error(fmt.Sprintf("unable to build `SetFromServer` object: %v", err))
			return
		}

		// Forward the message to the WebSocket server
		err = s.Hub.Write(buf.Bytes())
		if err != nil {
			s.Logger.Error(fmt.Sprintf("error writing to web socket %v", err))
		}
		buf.Reset()
	})
	if err != nil {
		log.Fatalf("Error subscribing to NATS subject: %v", err)
	}
	return sub, err
}

type AddObject struct {
	Geometry Geometry  `msgpack:"geometry"`
	Path     string    `msgpack:"path"`
	Position []float64 `msgpack:"position"`
}

// SetObject handler
func (s *Server) setGeometrySubscription() (*nats.Subscription, error) {
	sub, err := s.NATS.Subscribe("meshcat.geometries", func(msg *nats.Msg) {
		// shape := strings.Split(string(msg.Subject), ".")[2]
		shape := ""
		// todo: check here to see if the shape is available
		var buf bytes.Buffer
		enc := msgpack.NewEncoder(&buf)
		log.Printf("Received meshcat message from NATS: %s: shape %s", string(msg.Data), shape)

		// let's say for now the message has the form "object_name.path.positionx.positiony.positionz"
		// todo: command parsing for Geometry message data
		// based on the message type at the given path, parse the atrtibutes for the particular geometry

		// let's say for now the message has the form "object_name path positionx positiony positionz"
		if shape == "box" {
			box := Box{}
			err := json.Unmarshal(msg.Data, &box)
			if err != nil {
				s.Logger.Info(fmt.Sprintf("error processing add object request %v", err))
				return
			}
			box.init_element()
			obj := Objectify(&box)
			err = enc.Encode(SetObject{
				Object: obj,
				Command: Command{
					Type: "set_object",
					Path: "environment/box_geometries",
				},
			})
			if err != nil {
				s.Logger.Info(fmt.Sprintf("error processing add object request %v", err))
			}
		} else if shape == "sphere" {
			sphere := Sphere{}
			err := json.Unmarshal(msg.Data, &sphere)
			if err != nil {
				s.Logger.Info(fmt.Sprintf("error processing add object request %v", err))
				return
			}
			sphere.init_element()
			s.Logger.Info(fmt.Sprintf("sphere: %#v", sphere))
			obj := Objectify(&sphere)
			err = enc.Encode(SetObject{
				Object: obj,
				Command: Command{
					Type: "set_object",
					Path: "environment/sphere_geometries",
				},
			})
			if err != nil {
				s.Logger.Info(fmt.Sprintf("error processing add object request %v", err))
			}
		} else {
			var geom GenericGeom
			err := json.Unmarshal(msg.Data, &geom)
			if err != nil {
				s.Logger.Info(fmt.Sprintf("error processing add object request %v", err))
			}
			geom.init_element()
			obj := Objectify(geom)
			err = enc.Encode(SetObject{
				Object: obj,
				Command: Command{
					Type: "set_object",
					Path: fmt.Sprintf("environment/%v", "geometries"),
				},
			})
			if err != nil {
				s.Logger.Info(fmt.Sprintf("error encoding generaic add object request %v", err))
			}

		}
		// Forward the message to the WebSocket server
		err := s.Hub.Write(buf.Bytes())
		if err != nil {
			s.Logger.Error(fmt.Sprintf("error writing to web socket %v", err))
		}
		buf.Reset()
	})
	if err != nil {
		log.Fatalf("Error subscribing to NATS subject: %v", err)
	}
	return sub, err
}

type TransformationCommand struct {
	Matrix4     []float64 `msgpack:"matrix"`
	Translation []float64 `msgpack:"translation"`
	Rotation    []float64 `msgpack:"rotation"`
	Scale       []float64 `msgpack:"scale"`
}

type SetTransformationCommand struct {
	Command
	Object TransformationCommand `msgpack:"object"`
}

// eulerToRotationMatrix converts Euler angles to a rotation matrix.
// The Euler angles are represented as an array of 3 float64 values: [roll, pitch, yaw].
// Converts to the aerospace sequence of rotations: ZYX, applied in order from right to left.
func eulerToRotationMatrix(e [3]float64) *mat.Dense {
	c1, s1 := math.Cos(e[0]), math.Sin(e[0])
	c2, s2 := math.Cos(e[1]), math.Sin(e[1])
	c3, s3 := math.Cos(e[2]), math.Sin(e[2])

	return mat.NewDense(3, 3, []float64{
		c2 * c3, c2 * s3, -s2,
		c3*s1*s2 - c1*s3, c1*c3 + s1*s2*s3, c2 * s1,
		c3*c1*s2 + s1*s3, c1*s2*s3 - c3*s1, c2 * c1,
	})
}

func scalingMatrix(s [3]float64) *mat.Dense {
	return mat.NewDense(4, 4, []float64{
		s[0], 0, 0, 0,
		0, s[1], 0, 0,
		0, 0, s[2], 0,
		0, 0, 0, 1,
	})
}

func translationMatrix(t [3]float64) *mat.Dense {
	return mat.NewDense(4, 4, []float64{
		1, 0, 0, t[0],
		0, 1, 0, t[1],
		0, 0, 1, t[2],
		0, 0, 0, 1,
	})
}

func NewTransformation(data []byte) (transformation_matrix TransformationCommand, err error) {
	err = json.Unmarshal(data, &transformation_matrix)
	if err != nil {
		return transformation_matrix, fmt.Errorf("unable to unmarshal transformation matrix: %v", err)
	}

	// check if Matrix4 is already specified, in which case, just return the result
	if transformation_matrix.Matrix4 != nil && len(transformation_matrix.Matrix4) == 16 {
		for _, v := range transformation_matrix.Matrix4 {
			if v != 0.0 {
				break
			}
		}
		return transformation_matrix, nil
	}

	// determine rotation type, normalize to quaternion
	rotation := transformation_matrix.Rotation
	if rotation == nil {
		rotation = []float64{0, 0, 0, 1}
	} else if len(rotation) == 3 {
		// euler to quaternion
		rotation, err = transformations.EulerToQuaternion(([3]float64)(rotation))
		if err != nil {
			return transformation_matrix, fmt.Errorf("unable to convert euler angles to quaternion: %v", err)
		}
	} else if len(rotation) == 4 {
		// determine if the quaternion is valid, and normalized
		rotation = transformations.NewQuaternionFromSlice(rotation)
	}
	transformation_matrix.Rotation = rotation

	// todo: handle scaling matrix. Not a high priority for now
	return transformation_matrix, nil
}

// SetObject handler
func (s *Server) setTransformationSubscription() (*nats.Subscription, error) {
	sub, err := s.NATS.Subscribe("meshcat.transformations.>", func(msg *nats.Msg) {
		s.Logger.Info(fmt.Sprintf("Received meshcat message from NATS `%s` on subject `%s`", string(msg.Data), strings.Split(msg.Subject, ".")[2:]))
		path := strings.Join(strings.Split(string(msg.Subject), ".")[2:], "/")

		var buf bytes.Buffer
		enc := msgpack.NewEncoder(&buf)

		transformation_matrix, err := NewTransformation(msg.Data)
		if err != nil {
			s.Logger.Error(fmt.Sprintf("unable to build `TransformationCommand` object: %v", err))
			return
		}

		s.Logger.Info(fmt.Sprintf("transformation matrix: %v", transformation_matrix))
		err = enc.Encode(SetTransformationCommand{
			Object: transformation_matrix,
			Command: Command{
				Type: "set_transform",
				Path: path,
			},
		})
		if err != nil {
			log.Printf("error sending msg: %v", err)
		}

		// Forward the message to the WebSocket server
		err = s.Hub.Write(buf.Bytes())
		if err != nil {
			s.Logger.Error(fmt.Sprintf("error writing to web socket %v", err))
		}
		buf.Reset()
	})
	if err != nil {
		log.Fatalf("Error subscribing to NATS subject: %v", err)
	}
	return sub, err
}

func (s *Server) missionSubscription() (*nats.Subscription, error) {
	sub, err := s.NATS.Subscribe("meshcat.mission.>", func(msg *nats.Msg) {
		path := []string{"meshcat.transformations"}
		path = append(path, strings.Join(strings.Split(string(msg.Subject), ".")[2:], "."))
		full_path := strings.Join(path, ".")

		s.Logger.Info(fmt.Sprintf("Received meshcat message from NATS: %s on path %s", string(msg.Data), path))
		s.Q.Add(MissionWork{Conn: s.NATS, Path: full_path, Type: "orbit", Radius: 1, Omega: 1})
	})
	if err != nil {
		log.Fatalf("Error subscribing to NATS subject: %v", err)
	}
	return sub, err
}

func ParseFloats(x, y, z string) (fx, fy, fz float64, err error) {
	fx, err = strconv.ParseFloat(x, 64)
	if err != nil {
		return fx, fy, fz, err
	}
	fy, err = strconv.ParseFloat(y, 64)
	if err != nil {
		return fx, fy, fz, err
	}
	fz, err = strconv.ParseFloat(z, 64)
	if err != nil {
		return fx, fy, fz, err
	}
	return fx, fy, fz, err
}

func (s *Server) delete() (*nats.Subscription, error) {
	return nil, nil
}
