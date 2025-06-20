// model.go - Part of the `inventory` Package
//
//     ॐ भूर्भुवः स्वः
//     तत्स॑वि॒तुर्वरे॑ण्यं॒
//    भर्गो॑ दे॒वस्य॑ धीमहि।
//   धियो॒ यो नः॑ प्रचो॒दया॑त्॥
//
//
//  बी.वी.एल - बोसजी के द्वारा रचित भंडार लेखांकन हेतु तन्त्राक्ष्।
// =============================================
//
// एक सुगम एवं उपयोगी भंडार संचालन हेतु तन्त्राक्ष्।
//
// एक रचनात्मक भारतीय उत्पाद ।
//
// bvl - Boseji's Inventory Management Program
//
// Easy to use and useful stock, goods and materials handling software.
//
// Sources
// -------
// https://github.com/boseji/bvl
//
// License
// -------
//
//   bvl - Boseji's Inventory Management Program.
//   Copyright (C) 2025 by Abhijit Bose (aka. Boseji)
//
//   This program is free software: you can redistribute it and/or modify
//   it under the terms of the GNU General Public License version 2 only
//   as published by the Free Software Foundation.
//
//   This program is distributed in the hope that it will be useful,
//   but WITHOUT ANY WARRANTY; without even the implied warranty of
//   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
//
//   You should have received a copy of the GNU General Public License
//   along with this program. If not, see <https://www.gnu.org/licenses/>.
//
//  SPDX-License-Identifier: GPL-2.0-only
//  Full Name: GNU General Public License v2.0 only
//  Please visit <https://spdx.org/licenses/GPL-2.0-only.html> for details.
//

//
// Model definition for Inventory CLI
//

package inventory

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/boseji/bsg/gen"
)

// Item represents an inventory record.
//
// Fields:
//
//	ID          - auto-increment primary key
//	Description - free text
//	Location    - free text
//	Status      - free text
//	Remarks     - audit log, may contain timestamped entries
//
// The Remarks field is typically maintained using FormatRemarks()
// to ensure consistent timestamp format.
//
// Example:
//
//	[2025-06-21 14:30] installed new battery
//
// The Item struct is used across all DB, CSV, and JSON functions.
type Item struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
	Location    string `json:"location"`
	Status      string `json:"status"`
	Remarks     string `json:"remarks"`
}

var reLogPrefix = regexp.MustCompile(`^\[\d{4}-\d{2}-\d{2} \d{2}:\d{2}\]`)

// FormatRemarks returns the Remarks field formatted as:
//
//	[YYYY-MM-DD HH:MM] <remarks>
//
// If Remarks is already formatted (starts with timestamp),
// it returns Remarks unchanged.
//
// If Remarks is blank, returns only timestamp prefix.
//
// Usage:
//
//	formatted := item.FormatRemarks()
//
// This function is used by AddItem(), AppendItem(), EditItem()
// to ensure Remarks field is consistent.
//
// Example output:
//
//	"[2025-06-21 15:00] installed UPS"
func (item *Item) FormatRemarks() string {
	ts := gen.BST().Format("2006-01-02 15:04")
	r := strings.TrimSpace(item.Remarks)

	if r == "" {
		return fmt.Sprintf("[%s] ", ts)
	}

	if reLogPrefix.MatchString(r) {
		return item.Remarks
	}

	return fmt.Sprintf("[%s] %s", ts, r)
}
