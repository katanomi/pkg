/*
Copyright 2023 The Katanomi Authors.

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

// StoragePluginClassLabelKey for labeling owner StoragePluginClass
const StoragePluginClassLabelKey = "storage.katanomi.dev/storagePluginClass"

// LastModifiedAnnotation for recording last modified time
const LastModifiedAnnotation = "storage.katanomi.dev/lastModified"

// FileTypeAnnotation for recording business file type of file object
const FileTypeAnnotation = "storage.katanomi.dev/fileType"

// EntryAnnotation for recording file entry of directory
const EntryAnnotation = "storage.katanomi.dev/entry"

// AnnotationPrefix is prefix of user-defined annotations
const AnnotationPrefix = "storage.katanomi.dev/annotation"
