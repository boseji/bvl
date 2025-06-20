// inventory_test.go - Part of Tests for the `inventory` Package
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
	"fmt"
	"testing"

	"github.com/boseji/bvl/inventory"
)

func setupInventoryDB(t *testing.T) *inventory.InventoryDB {
	inv := inventory.NewInventoryDB(":memory:")
	if inv == nil {
		t.Fatal("failed to create InventoryDB")
	}
	return inv
}

func TestInventoryDB_ResetSequence(t *testing.T) {
	inv := setupInventoryDB(t)
	defer inv.Close()

	err := inv.ResetSequence()
	if err != nil {
		t.Fatalf("ResetSequence failed: %v", err)
	}
}

func TestInventoryDB_AddAndGetItem(t *testing.T) {
	inv := setupInventoryDB(t)
	defer inv.Close()

	item := inventory.Item{
		Description: "UPS", Location: "Rack 1",
		Status: "Operational", Remarks: "installed",
	}
	err := inv.AddItem(item)
	if err != nil {
		t.Fatalf("AddItem failed: %v", err)
	}

	items, err := inv.ListAll()
	if err != nil || len(items) != 1 {
		t.Fatalf("ListAll failed: %v", err)
	}

	got, err := inv.GetItemByID(items[0].ID)
	if err != nil {
		t.Fatalf("GetItemByID failed: %v", err)
	}
	if got.Description != "UPS" {
		t.Errorf("unexpected Description: %s", got.Description)
	}
}

func TestInventoryDB_AppendItem(t *testing.T) {
	inv := setupInventoryDB(t)
	defer inv.Close()

	item := inventory.Item{
		ID:          1005,
		Description: "Battery", Location: "Shelf 1",
		Status: "New", Remarks: "received",
	}
	err := inv.AppendItem(item)
	if err != nil {
		t.Fatalf("AppendItem failed: %v", err)
	}

	got, err := inv.GetItemByID(1005)
	if err != nil {
		t.Fatalf("GetItemByID failed: %v", err)
	}
	if got.Description != "Battery" {
		t.Errorf("unexpected Description: %s", got.Description)
	}
}

func TestInventoryDB_EditItem(t *testing.T) {
	inv := setupInventoryDB(t)
	defer inv.Close()

	item := inventory.Item{
		Description: "Switch", Location: "Rack 2",
		Status: "Installed", Remarks: "new",
	}
	err := inv.AddItem(item)
	if err != nil {
		t.Fatalf("AddItem failed: %v", err)
	}

	items, _ := inv.ListAll()
	item = items[0]
	item.Status = "Active"
	item.Remarks = "activated"

	err = inv.EditItem(item)
	if err != nil {
		t.Fatalf("EditItem failed: %v", err)
	}

	got, _ := inv.GetItemByID(item.ID)
	if got.Status != "Active" {
		t.Errorf("unexpected Status: %s", got.Status)
	}
}

func TestInventoryDB_EditItem_NotFound(t *testing.T) {
	inv := setupInventoryDB(t)
	defer inv.Close()

	item := inventory.Item{
		ID: 9999, Description: "Ghost", Location: "Void",
		Status: "Lost", Remarks: "none",
	}
	err := inv.EditItem(item)
	if err != nil {
		t.Fatalf("EditItem on missing ID should not fail: %v", err)
	}
}

func TestInventoryDB_AppendRemarksEntry(t *testing.T) {
	inv := setupInventoryDB(t)
	defer inv.Close()

	item := inventory.Item{
		Description: "Router", Location: "Rack 3",
		Status: "Operational", Remarks: "setup",
	}
	_ = inv.AddItem(item)
	items, _ := inv.ListAll()

	err := inv.AppendRemarksEntry(items[0].ID, "tested")
	if err != nil {
		t.Fatalf("AppendRemarksEntry failed: %v", err)
	}
}

func TestInventoryDB_DeleteItem(t *testing.T) {
	inv := setupInventoryDB(t)
	defer inv.Close()

	item := inventory.Item{
		Description: "Server", Location: "Rack 4",
		Status: "Provisioning", Remarks: "image install",
	}
	_ = inv.AddItem(item)
	items, _ := inv.ListAll()

	err := inv.DeleteItem(items[0].ID)
	if err != nil {
		t.Fatalf("DeleteItem failed: %v", err)
	}

	_, err = inv.GetItemByID(items[0].ID)
	if err == nil {
		t.Errorf("expected error for deleted item")
	}
}

func TestInventoryDB_ListItemsPaged(t *testing.T) {
	inv := setupInventoryDB(t)
	defer inv.Close()

	for i := 0; i < 5; i++ {
		item := inventory.Item{
			Description: "Device", Location: "Loc",
			Status: "OK", Remarks: "added",
		}
		_ = inv.AddItem(item)
	}

	items, err := inv.ListItemsPaged(0, 3)
	if err != nil || len(items) != 3 {
		t.Fatalf("ListItemsPaged failed: %v", err)
	}

	empty, err := inv.ListItemsPaged(9999, 5)
	if err != nil {
		t.Fatalf("ListItemsPaged failed: %v", err)
	}
	if len(empty) != 0 {
		t.Errorf("expected empty slice for large afterID")
	}
}

func TestInventoryDB_ItemIterator(t *testing.T) {
	inv := setupInventoryDB(t)
	defer inv.Close()

	for i := 0; i < 3; i++ {
		item := inventory.Item{
			Description: "FW", Location: "Loc",
			Status: "Ready", Remarks: "fw added",
		}
		_ = inv.AddItem(item)
	}

	iter, err := inv.NewItemIterator("WHERE status = ?", "Ready")
	if err != nil {
		t.Fatalf("NewItemIterator failed: %v", err)
	}
	defer iter.Close()

	count := 0
	for {
		item, ok, err := iter.Next()
		if err != nil {
			t.Fatalf("Iterator Next failed: %v", err)
		}
		if !ok {
			break
		}
		if item.Status != "Ready" {
			t.Errorf("unexpected status: %s", item.Status)
		}
		count++
	}
	if count != 3 {
		t.Errorf("expected 3 items, got %d", count)
	}
}

func TestInventoryDB_ItemIterator_BadWhereClause(t *testing.T) {
	inv := setupInventoryDB(t)
	defer inv.Close()

	_, err := inv.NewItemIterator("WHERE no_such_field = 1")
	if err == nil {
		t.Errorf("expected error for bad WHERE clause")
	}
}

func TestInventoryDB_WithTransactionRollback(t *testing.T) {
	inv := setupInventoryDB(t)
	defer inv.Close()

	err := inv.WithTransaction(func(tx inventory.Execer) error {
		item := inventory.Item{
			Description: "Test", Location: "Loc",
			Status: "Temp", Remarks: "test",
		}
		if err := inventory.AddItem(tx, item); err != nil {
			return err
		}
		return fmt.Errorf("force rollback")
	})

	if err == nil {
		t.Fatalf("expected error from transaction rollback")
	}

	items, _ := inv.ListAll()
	if len(items) != 0 {
		t.Errorf("expected no items after rollback, found %d", len(items))
	}
}

func TestInventoryDB_WithTransaction_BeginFails(t *testing.T) {
	inv := setupInventoryDB(t)
	// close DB to simulate begin fail
	inv.Close()

	err := inv.WithTransaction(func(tx inventory.Execer) error {
		return nil
	})

	if err == nil {
		t.Fatalf("expected error from begin tx fail")
	}
}

func TestInventoryDB_WithTransaction_CommitFails(t *testing.T) {
	t.Skip("sqlite :memory: does not allow forcing commit fail reliably")
	inv := setupInventoryDB(t)
	defer inv.Close()

	// simulate commit fail by closing DB inside transaction
	err := inv.WithTransaction(func(tx inventory.Execer) error {
		// close DB mid-tx to force commit fail
		inv.Close()
		return nil
	})

	if err == nil {
		t.Fatalf("expected error from commit tx fail")
	}
}
