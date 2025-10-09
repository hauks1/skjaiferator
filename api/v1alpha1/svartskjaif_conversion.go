/*
Copyright 2025.

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

package v1alpha1

import (
	"sigs.k8s.io/controller-runtime/pkg/conversion"

	"github.com/hauks1/skjaiferator/api/v1beta1"
)

// ConvertTo converts this v1alpha1 SvartSkjaif to the Hub version (v1beta1).
func (src *SvartSkjaif) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1beta1.SvartSkjaif)

	// ObjectMeta
	dst.ObjectMeta = src.ObjectMeta

	// Spec - flatten the nested structure
	dst.Spec.Kopp = src.Spec.SvartSkjaifContainer.Kopp
	dst.Spec.Vann = src.Spec.SvartSkjaifContainer.Vann
	dst.Spec.Kaffe = src.Spec.SvartSkjaifContainer.Kaffe

	// Status
	// Add any status field conversions here

	return nil
}

// ConvertFrom converts from the Hub version (v1beta1) to this v1alpha1 SvartSkjaif.
func (dst *SvartSkjaif) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1beta1.SvartSkjaif)

	// ObjectMeta
	dst.ObjectMeta = src.ObjectMeta

	// Spec - nest the structure
	dst.Spec.SvartSkjaifContainer.Kopp = src.Spec.Kopp
	dst.Spec.SvartSkjaifContainer.Vann = src.Spec.Vann
	dst.Spec.SvartSkjaifContainer.Kaffe = src.Spec.Kaffe

	// Status
	// Add any status field conversions here

	return nil
}
