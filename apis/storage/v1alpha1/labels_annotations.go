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
const StoragePluginClassLabelKey = "storage.katanomi.dev/storagePluginClass" // NOSONAR // ignore: "Key" detected here, make sure this is not a hard-coded credential

// LastModifiedAnnotation for recording last modified time
const LastModifiedAnnotation = "storage.katanomi.dev/lastModified"

// FileTypeAnnotation for recording business file type of file object
const FileTypeAnnotation = "storage.katanomi.dev/fileType"

// StorageEntryAnnotation for recording file entry of directory
const StorageEntryAnnotation = "storage.katanomi.dev/entry"

// StorageAnnotationPrefix is prefix of user-defined annotations
const StorageAnnotationPrefix = "storage.katanomi.dev/annotation."
