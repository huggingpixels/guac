//
// Copyright 2023 The GUAC Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package backend

import (
	"crypto/sha256"
	"fmt"
	"strings"

	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
	"github.com/guacsec/guac/pkg/assembler/graphql/model"
)

const (
	// guacIDSplit is used as a separator to concatenate the type and namespace to create an ID
	guacIDSplit = "guac-split-@@"
)

type globalID struct {
	nodeType string
	id       string
}

func toGlobalID(nodeType string, id string) string {
	return strings.Join([]string{nodeType, id}, ":")
}

func toGlobalIDs(nodeType string, ids []string) []string {
	var globalIDs []string
	for _, id := range ids {
		globalIDs = append(globalIDs, strings.Join([]string{nodeType, id}, ":"))
	}
	return globalIDs
}

func fromGlobalID(gID string) *globalID {
	idSplit := strings.Split(string(gID), ":")
	if len(idSplit) == 2 {
		return &globalID{
			nodeType: idSplit[0],
			id:       idSplit[1],
		}
	} else {
		return &globalID{
			id:       idSplit[0],
			nodeType: "",
		}
	}
}

func IDEQ(id string) func(*sql.Selector) {
	filterGlobalID := fromGlobalID(id)
	return sql.FieldEQ("id", filterGlobalID.id)
}

func NoOpSelector() func(*sql.Selector) {
	return func(s *sql.Selector) {}
}

type Predicate interface {
	~func(*sql.Selector)
}

func optionalPredicate[P Predicate, T any](value *T, fn func(s T) P) P {
	if value == nil {
		return NoOpSelector()
	}

	return fn(*value)
}

func ptrWithDefault[T any](value *T, defaultValue T) T {
	if value == nil {
		return defaultValue
	}

	return *value
}

func toPtrSlice[T any](slice []T) []*T {
	ptrs := make([]*T, len(slice))
	for i := range slice {
		ptrs[i] = &slice[i]
	}
	return ptrs
}

// func fromPtrSlice[T any](slice []*T) []T {
// 	ptrs := make([]T, len(slice))
// 	for i := range slice {
// 		if slice[i] == nil {
// 			continue
// 		}
// 		ptrs[i] = *slice[i]
// 	}
// 	return ptrs
// }

func toLowerPtr(s *string) *string {
	if s == nil {
		return nil
	}
	lower := strings.ToLower(*s)
	return &lower
}

func chunk[T any](collection []T, size int) [][]T {
	if size <= 0 {
		panic("Second parameter must be greater than 0")
	}

	chunksNum := len(collection) / size
	if len(collection)%size != 0 {
		chunksNum += 1
	}

	result := make([][]T, 0, chunksNum)

	for i := 0; i < chunksNum; i++ {
		last := (i + 1) * size
		if last > len(collection) {
			last = len(collection)
		}
		result = append(result, collection[i*size:last])
	}

	return result
}

// generateUUIDKey is used to generate the ID based on the sha256 hash of the content of the inputSpec that is passed in.
// For example, for artifact it would be
// artifactID := uuid.NewHash(sha256.New(), uuid.NameSpaceDNS, []byte(helpers.GetKey[*model.ArtifactInputSpec, string](artInput.ArtifactInput, helpers.ArtifactServerKey)), 5)
// where the data is generated by converting the artifactInputSpec into a canonicalized key
func generateUUIDKey(data []byte) uuid.UUID {
	return uuid.NewHash(sha256.New(), uuid.NameSpaceDNS, data, 5)
}

func getIDfromNode(node model.Node) (string, error) {
	switch v := node.(type) {
	case *model.Package:
		if v != nil && len(v.Namespaces) > 0 && len(v.Namespaces[0].Names) > 0 && len(v.Namespaces[0].Names[0].Versions) > 0 {
			return v.Namespaces[0].Names[0].Versions[0].ID, nil
		} else if v != nil && len(v.Namespaces) > 0 && len(v.Namespaces[0].Names) > 0 {
			return v.Namespaces[0].Names[0].ID, nil
		} else if v != nil && len(v.Namespaces) > 0 {
			return v.Namespaces[0].ID, nil
		} else {
			return v.ID, nil
		}
	case *model.Artifact:
		return v.ID, nil
	case *model.Builder:
		return v.ID, nil
	case *model.Source:
		if v != nil && len(v.Namespaces) > 0 && len(v.Namespaces[0].Names) > 0 {
			return v.Namespaces[0].Names[0].ID, nil
		} else if v != nil && len(v.Namespaces) > 0 {
			return v.Namespaces[0].ID, nil
		} else {
			return v.ID, nil
		}
	case *model.Vulnerability:
		if len(v.VulnerabilityIDs) > 0 {
			return v.VulnerabilityIDs[0].ID, nil
		} else {
			return v.ID, nil
		}
	case *model.License:
		return v.ID, nil
	case *model.CertifyBad:
		return v.ID, nil
	case *model.CertifyGood:
		return v.ID, nil
	case *model.CertifyLegal:
		return v.ID, nil
	case *model.CertifyScorecard:
		return v.ID, nil
	case *model.CertifyVEXStatement:
		return v.ID, nil
	case *model.CertifyVuln:
		return v.ID, nil
	case *model.HashEqual:
		return v.ID, nil
	case *model.HasMetadata:
		return v.ID, nil
	case *model.HasSbom:
		return v.ID, nil
	case *model.HasSlsa:
		return v.ID, nil
	case *model.HasSourceAt:
		return v.ID, nil
	case *model.IsDependency:
		return v.ID, nil
	case *model.IsOccurrence:
		return v.ID, nil
	case *model.PkgEqual:
		return v.ID, nil
	case *model.PointOfContact:
		return v.ID, nil
	case *model.VulnEqual:
		return v.ID, nil
	case *model.VulnerabilityMetadata:
		return v.ID, nil
	default:
		return "", fmt.Errorf("unknown type: %v", v)
	}
}
