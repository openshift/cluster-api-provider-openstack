/*
Copyright 2019 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package yaml

import (
	"bufio"
	"bytes"
	"io"
	"os"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer/streaming"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes/scheme"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha2"
)

type ParseInput struct {
	File string
}

type ParseOutput struct {
	Clusters            []*clusterv1.Cluster
	Machines            []*clusterv1.Machine
	MachineSets         []*clusterv1.MachineSet
	MachineDeployments  []*clusterv1.MachineDeployment
	UnstructuredObjects []*unstructured.Unstructured
}

// Parse extracts runtime objects from a file.
func Parse(input ParseInput) (*ParseOutput, error) {
	output := &ParseOutput{}

	// Open the input file.
	reader, err := os.Open(input.File)
	if err != nil {
		return nil, err
	}

	// Create a new decoder.
	decoder := NewYAMLDecoder(reader)
	defer decoder.Close()

	for {
		u := &unstructured.Unstructured{}
		_, gvk, err := decoder.Decode(nil, u)
		if err == io.EOF {
			break
		}
		if runtime.IsNotRegisteredError(err) {
			continue
		}
		if err != nil {
			return nil, err
		}

		switch gvk.Kind {
		case "Cluster":
			obj := &clusterv1.Cluster{}
			if err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.Object, obj); err != nil {
				return nil, errors.Wrapf(err, "cannot convert object to %s", gvk.Kind)
			}
			output.Clusters = append(output.Clusters, obj)
		case "Machine":
			obj := &clusterv1.Machine{}
			if err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.Object, obj); err != nil {
				return nil, errors.Wrapf(err, "cannot convert object to %s", gvk.Kind)
			}
			output.Machines = append(output.Machines, obj)
		case "MachineSet":
			obj := &clusterv1.MachineSet{}
			if err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.Object, obj); err != nil {
				return nil, errors.Wrapf(err, "cannot convert object to %s", gvk.Kind)
			}
			output.MachineSets = append(output.MachineSets, obj)
		case "MachineDeployment":
			obj := &clusterv1.MachineDeployment{}
			if err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.Object, obj); err != nil {
				return nil, errors.Wrapf(err, "cannot convert object to %s", gvk.Kind)
			}
			output.MachineDeployments = append(output.MachineDeployments, obj)
		default:
			output.UnstructuredObjects = append(output.UnstructuredObjects, u)
		}

	}

	return output, nil
}

type yamlDecoder struct {
	reader  *yaml.YAMLReader
	decoder runtime.Decoder
	close   func() error
}

func (d *yamlDecoder) Decode(defaults *schema.GroupVersionKind, into runtime.Object) (runtime.Object, *schema.GroupVersionKind, error) {
	for {
		doc, err := d.reader.Read()
		if err != nil {
			return nil, nil, err
		}

		//  Skip over empty documents, i.e. a leading `---`
		if len(bytes.TrimSpace(doc)) == 0 {
			continue
		}

		return d.decoder.Decode(doc, defaults, into)
	}

}

func (d *yamlDecoder) Close() error {
	return d.close()
}

func NewYAMLDecoder(r io.ReadCloser) streaming.Decoder {
	return &yamlDecoder{
		reader:  yaml.NewYAMLReader(bufio.NewReader(r)),
		decoder: scheme.Codecs.UniversalDeserializer(),
		close:   r.Close,
	}
}
