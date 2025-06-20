// model_test.go - Part of Tests for the `inventory` Package
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

package inventory_test

import (
	"testing"

	"github.com/boseji/bvl/inventory"
)

func TestItem_FormatRemarks(t *testing.T) {
	item := inventory.Item{
		Description: "UPS",
		Location:    "Rack 1",
		Status:      "Operational",
		Remarks:     "installed new battery",
	}

	formatted := item.FormatRemarks()

	if len(formatted) == 0 {
		t.Fatalf("FormatRemarks returned empty string")
	}

	if formatted == item.Remarks {
		t.Errorf("FormatRemarks did not change original remarks")
	}

	if formatted[0] != '[' {
		t.Errorf("FormatRemarks missing timestamp prefix")
	}

	t.Logf("Formatted remarks: %q", formatted)
}

func TestItem_FormatRemarks_AlreadyFormatted(t *testing.T) {
	oldRemark := "[2025-06-20 12:00] initial install"

	item := inventory.Item{
		Description: "UPS",
		Location:    "Rack 1",
		Status:      "Operational",
		Remarks:     oldRemark,
	}

	formatted := item.FormatRemarks()

	if formatted != oldRemark {
		t.Errorf("FormatRemarks changed already formatted remarks")
	}

	t.Logf("Already formatted remarks: %q", formatted)
}

func TestItem_FormatRemarks_Blank(t *testing.T) {
	item := inventory.Item{
		Description: "UPS",
		Location:    "Rack 1",
		Status:      "Operational",
		Remarks:     "",
	}

	formatted := item.FormatRemarks()

	if len(formatted) == 0 {
		t.Fatalf("FormatRemarks returned empty string for blank input")
	}

	if formatted[0] != '[' {
		t.Errorf("FormatRemarks missing timestamp prefix on blank input")
	}

	t.Logf("Formatted blank remarks: %q", formatted)
}
