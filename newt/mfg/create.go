/**
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package mfg

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"mynewt.apache.org/newt/newt/builder"
	"mynewt.apache.org/newt/newt/flash"
	"mynewt.apache.org/newt/newt/pkg"
	"mynewt.apache.org/newt/newt/target"
	"mynewt.apache.org/newt/util"
)

type mfgManifest struct {
	BuildTime string `json:"build_time"`
	MfgHash   string `json:"mfg_hash"`
}

func insertPartIntoBlob(blob []byte, part mfgPart) {
	partEnd := part.offset + len(part.data)

	if len(blob) < partEnd {
		panic("internal error; mfg blob too small")
	}

	copy(blob[part.offset:partEnd], part.data)
}

func (mi *MfgImage) partFromImage(
	imgPath string, flashAreaName string) (mfgPart, error) {

	part := mfgPart{}

	area, ok := mi.bsp.FlashMap.Areas[flashAreaName]
	if !ok {
		return part, util.FmtNewtError(
			"Image at \"%s\" requires undefined flash area \"%s\"",
			imgPath, flashAreaName)
	}

	part.name = fmt.Sprintf("%s (%s)", flashAreaName, filepath.Base(imgPath))
	part.offset = area.Offset

	var err error

	part.data, err = ioutil.ReadFile(imgPath)
	if err != nil {
		return part, util.ChildNewtError(err)
	}

	overflow := len(part.data) - area.Size
	if overflow > 0 {
		return part, util.FmtNewtError(
			"Image \"%s\" is too large to fit in flash area \"%s\"; "+
				"image-size=%d flash-area-size=%d overflow=%d",
			imgPath, flashAreaName, len(part.data), area.Size, overflow)
	}

	return part, nil
}

func (mi *MfgImage) section0Size() int {
	greatest := 0

	bootArea := mi.bsp.FlashMap.Areas[flash.FLASH_AREA_NAME_BOOTLOADER]
	image0Area := mi.bsp.FlashMap.Areas[flash.FLASH_AREA_NAME_IMAGE_0]
	image1Area := mi.bsp.FlashMap.Areas[flash.FLASH_AREA_NAME_IMAGE_1]

	if mi.boot != nil {
		greatest = util.IntMax(greatest, bootArea.Offset+bootArea.Size)
	}
	if len(mi.images) >= 1 {
		greatest = util.IntMax(greatest, image0Area.Offset+image0Area.Size)
	}
	if len(mi.images) >= 2 {
		greatest = util.IntMax(greatest, image1Area.Offset+image1Area.Size)
	}

	for _, entry := range mi.rawEntries {
		greatest = util.IntMax(greatest, entry.offset+len(entry.data))
	}

	return greatest
}

// @return						section-0-blob, hash-offset, error
func (mi *MfgImage) section0Data(parts []mfgPart) ([]byte, int, error) {
	blobSize := mi.section0Size()
	blob := make([]byte, blobSize)

	// Initialize section 0's data as unwritten flash (0xff).
	for i, _ := range blob {
		blob[i] = 0xff
	}

	for _, part := range parts {
		insertPartIntoBlob(blob, part)
	}

	hashOffset, err := insertMeta(blob, mi.bsp.FlashMap)
	if err != nil {
		return nil, 0, err
	}

	return blob, hashOffset, nil

}

// @return						[section0blob, section1blob,...], hash, err
func (mi *MfgImage) createSections(parts []mfgPart) (
	[][]byte, []byte, error) {

	section0Data, hashOff, err := mi.section0Data(parts)
	if err != nil {
		return nil, nil, err
	}

	// XXX: Append additional sections.

	// Calculate manufacturing has.
	sections := [][]byte{section0Data}
	hash := calcMetaHash(sections)

	// Write hash to meta region in section 0.
	copy(section0Data[hashOff:hashOff+META_HASH_SZ], hash)

	return sections, hash, nil
}

func areaNameFromImgIdx(imgIdx int) (string, error) {
	switch imgIdx {
	case 0:
		return flash.FLASH_AREA_NAME_IMAGE_0, nil
	case 1:
		return flash.FLASH_AREA_NAME_IMAGE_1, nil
	default:
		return "", util.FmtNewtError("invalid image index: %d", imgIdx)
	}
}

func (mi *MfgImage) rawEntryParts() []mfgPart {
	parts := make([]mfgPart, len(mi.rawEntries))
	for i, entry := range mi.rawEntries {
		parts[i].name = fmt.Sprintf("entry-%d (%s)", i, entry.filename)
		parts[i].offset = entry.offset
		parts[i].data = entry.data
	}

	return parts
}

func bootLoaderFromPaths(t *target.Target) []string {
	return []string{
		/* boot.elf */
		builder.AppElfPath(t.Name(), builder.BUILD_NAME_APP, t.App().Name()),

		/* boot.elf.bin */
		builder.AppBinPath(t.Name(), builder.BUILD_NAME_APP, t.App().Name()),

		/* manifest.json */
		builder.ManifestPath(t.Name(), builder.BUILD_NAME_APP, t.App().Name()),
	}
}

func loaderFromPaths(t *target.Target) []string {
	if t.LoaderName == "" {
		return nil
	}

	return []string{
		/* <loader>.elf */
		builder.AppElfPath(t.Name(), builder.BUILD_NAME_LOADER,
			t.Loader().Name()),

		/* <app>.img */
		builder.AppImgPath(t.Name(), builder.BUILD_NAME_LOADER,
			t.Loader().Name()),
	}
}

func appFromPaths(t *target.Target) []string {
	return []string{
		/* <app>.elf */
		builder.AppElfPath(t.Name(), builder.BUILD_NAME_APP, t.App().Name()),

		/* <app>.img */
		builder.AppImgPath(t.Name(), builder.BUILD_NAME_APP, t.App().Name()),

		/* manifest.json */
		builder.ManifestPath(t.Name(), builder.BUILD_NAME_APP, t.App().Name()),
	}
}

func imageFromPaths(t *target.Target) []string {
	paths := loaderFromPaths(t)
	paths = append(paths, appFromPaths(t)...)
	return paths
}

func (mi *MfgImage) copyBinFile(srcPath string, dstDir string) error {
	dstPath := dstDir + "/" + filepath.Base(srcPath)

	util.StatusMessage(util.VERBOSITY_VERBOSE, "copying file %s --> %s\n",
		srcPath, dstPath)

	if err := util.CopyFile(srcPath, dstPath); err != nil {
		return err
	}

	return nil
}

func (mi *MfgImage) copyBinFiles() error {
	dstPath := MfgBinDir(mi.basePkg.Name())
	if err := os.MkdirAll(filepath.Dir(dstPath), 0755); err != nil {
		return util.ChildNewtError(err)
	}

	bootPaths := bootLoaderFromPaths(mi.boot)
	for _, path := range bootPaths {
		dstDir := MfgBootDir(mi.basePkg.Name())
		if err := mi.copyBinFile(path, dstDir); err != nil {
			return err
		}
	}

	for i, imgTarget := range mi.images {
		imgPaths := imageFromPaths(imgTarget)
		dstDir := MfgImageBinDir(mi.basePkg.Name(), i)
		for _, path := range imgPaths {
			if err := mi.copyBinFile(path, dstDir); err != nil {
				return err
			}
		}
	}

	return nil
}

func (mi *MfgImage) dstBootBinPath() string {
	if mi.boot == nil {
		return ""
	}

	return fmt.Sprintf("%s/%s.elf.bin",
		MfgBootDir(mi.basePkg.Name()),
		pkg.ShortName(mi.boot.App()))
}

func (mi *MfgImage) dstImgPath(slotIdx int) string {
	var pack *pkg.LocalPackage
	var imgIdx int

	if len(mi.images) >= 1 {
		switch slotIdx {
		case 0:
			if mi.images[0].LoaderName != "" {
				pack = mi.images[0].Loader()
			} else {
				pack = mi.images[0].App()
			}
			imgIdx = 0

		case 1:
			if mi.images[0].LoaderName != "" {
				pack = mi.images[0].App()
				imgIdx = 0
			} else {
				if len(mi.images) >= 2 {
					pack = mi.images[1].App()
				}
				imgIdx = 1
			}

		default:
			panic(fmt.Sprintf("invalid image index: %d", imgIdx))
		}
	}

	if pack == nil {
		return ""
	}

	return fmt.Sprintf("%s/%s.img",
		MfgImageBinDir(mi.basePkg.Name(), imgIdx), pkg.ShortName(pack))
}

func (mi *MfgImage) targetParts() ([]mfgPart, error) {
	parts := []mfgPart{}

	bootPath := mi.dstBootBinPath()
	if bootPath != "" {
		bootPart, err := mi.partFromImage(
			bootPath, flash.FLASH_AREA_NAME_BOOTLOADER)
		if err != nil {
			return nil, err
		}

		parts = append(parts, bootPart)
	}

	for i := 0; i < 2; i++ {
		imgPath := mi.dstImgPath(i)
		if imgPath != "" {
			areaName, err := areaNameFromImgIdx(i)
			if err != nil {
				return nil, err
			}

			part, err := mi.partFromImage(imgPath, areaName)
			if err != nil {
				return nil, err
			}
			parts = append(parts, part)
		}
	}

	return parts, nil
}

// Returns a slice containing the path of each file required to build the
// manufacturing image.
func (mi *MfgImage) FromPaths() []string {
	paths := []string{}

	if mi.boot != nil {
		paths = append(paths, bootLoaderFromPaths(mi.boot)...)
	}
	if len(mi.images) >= 1 {
		paths = append(paths, imageFromPaths(mi.images[0])...)
	}
	if len(mi.images) >= 2 {
		paths = append(paths, imageFromPaths(mi.images[1])...)
	}

	for _, raw := range mi.rawEntries {
		paths = append(paths, raw.filename)
	}

	return paths
}

// @return						[section0blob, section1blob,...], hash, err
func (mi *MfgImage) build() ([][]byte, []byte, error) {
	if err := mi.copyBinFiles(); err != nil {
		return nil, nil, err
	}

	targetParts, err := mi.targetParts()
	if err != nil {
		return nil, nil, err
	}

	rawParts := mi.rawEntryParts()

	parts := append(targetParts, rawParts...)
	sortParts(parts)

	sections, hash, err := mi.createSections(parts)
	if err != nil {
		return nil, nil, err
	}

	return sections, hash, nil
}

func (mi *MfgImage) createManifest(hash []byte) ([]byte, error) {
	manifest := mfgManifest{
		BuildTime: time.Now().Format(time.RFC3339),
		MfgHash:   fmt.Sprintf("%x", hash),
	}
	buffer, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return nil, util.FmtNewtError("Failed to encode mfg manifest: %s",
			err.Error())
	}

	return buffer, nil
}

func appendNonEmptyStr(dst []string, src string) []string {
	if src != "" {
		dst = append(dst, src)
	}

	return dst
}

func (mi *MfgImage) ToPaths() []string {
	paths := []string{}

	paths = appendNonEmptyStr(paths, mi.BootBinPath())
	paths = appendNonEmptyStr(paths, mi.BootElfPath())
	paths = appendNonEmptyStr(paths, mi.BootManifestPath())

	for i := 0; i < len(mi.images); i++ {
		paths = appendNonEmptyStr(paths, mi.LoaderImgPath(i))
		paths = appendNonEmptyStr(paths, mi.LoaderElfPath(i))
		paths = appendNonEmptyStr(paths, mi.AppImgPath(i))
		paths = appendNonEmptyStr(paths, mi.AppElfPath(i))
		paths = appendNonEmptyStr(paths, mi.ImageManifestPath(i))
	}

	paths = append(paths, mi.SectionBinPaths()...)
	paths = append(paths, mi.ManifestPath())

	return paths
}

// @return                      [paths-of-artifacts], error
func (mi *MfgImage) CreateMfgImage() ([]string, error) {
	sections, hash, err := mi.build()
	if err != nil {
		return nil, err
	}

	sectionDir := MfgSectionBinDir(mi.basePkg.Name())
	if err := os.MkdirAll(sectionDir, 0755); err != nil {
		return nil, util.ChildNewtError(err)
	}

	for i, section := range sections {
		sectionPath := MfgSectionBinPath(mi.basePkg.Name(), i)
		if err := ioutil.WriteFile(sectionPath, section, 0644); err != nil {
			return nil, util.ChildNewtError(err)
		}
	}

	manifest, err := mi.createManifest(hash)
	if err != nil {
		return nil, err
	}

	manifestPath := mi.ManifestPath()
	if err := ioutil.WriteFile(manifestPath, manifest, 0644); err != nil {
		return nil, util.FmtNewtError("Failed to write mfg manifest file: %s",
			err.Error())
	}

	return mi.ToPaths(), nil
}
