package logic

import (
	"testing"
	"zero-admin/rpc/sys/db/mysql/model"
	"zero-admin/rpc/sys/sysclient"
)

func TestBuildMenuTree(t *testing.T) {
	// 测试用例定义
	tests := []struct {
		name      string            // 测试用例名称
		menus     []*model.SysMenu  // 输入的菜单列表
		parentID  int64             // 父级ID
		want      []*sysclient.Menu // 期望的结果
		wantCount int               // 期望的菜单数量(用于简单验证)
	}{
		{
			name:      "TC001-空菜单列表",
			menus:     []*model.SysMenu{},
			parentID:  0,
			want:      []*sysclient.Menu{},
			wantCount: 0,
		},
		{
			name: "TC002-无匹配父ID",
			menus: []*model.SysMenu{
				{ID: 1, ParentID: 100, MenuName: "Menu1", Status: 1},
				{ID: 2, ParentID: 100, MenuName: "Menu2", Status: 1},
			},
			parentID:  0,
			want:      []*sysclient.Menu{},
			wantCount: 0,
		},
		{
			name: "TC003-状态禁用过滤",
			menus: []*model.SysMenu{
				{ID: 1, ParentID: 0, MenuName: "DisabledMenu", Status: 0},
				{ID: 2, ParentID: 0, MenuName: "EnabledMenu", Status: 1},
			},
			parentID: 0,
			want: []*sysclient.Menu{
				{
					Id:       2,
					ParentId: 0,
					MenuName: "EnabledMenu",
					Children: []*sysclient.Menu{},
				},
			},
			wantCount: 1,
		},
		{
			name: "TC004-单层菜单",
			menus: []*model.SysMenu{
				{ID: 1, ParentID: 0, MenuName: "Root1", Status: 1},
				{ID: 2, ParentID: 0, MenuName: "Root2", Status: 1},
				{ID: 3, ParentID: 0, MenuName: "Root3", Status: 1},
			},
			parentID: 0,
			want: []*sysclient.Menu{
				{
					Id:       1,
					ParentId: 0,
					MenuName: "Root1",
					Children: []*sysclient.Menu{},
				},
				{
					Id:       2,
					ParentId: 0,
					MenuName: "Root2",
					Children: []*sysclient.Menu{},
				},
				{
					Id:       3,
					ParentId: 0,
					MenuName: "Root3",
					Children: []*sysclient.Menu{},
				},
			},
			wantCount: 3,
		},
		{
			name: "TC005-多层嵌套菜单",
			menus: []*model.SysMenu{
				{ID: 1, ParentID: 0, MenuName: "Root", Status: 1},
				{ID: 2, ParentID: 1, MenuName: "Child1", Status: 1},
				{ID: 3, ParentID: 1, MenuName: "Child2", Status: 1},
				{ID: 4, ParentID: 2, MenuName: "GrandChild", Status: 1},
			},
			parentID:  0,
			wantCount: 1, // 只有一个根节点
		},
		{
			name: "TC006-混合场景-包含禁用和不匹配项",
			menus: []*model.SysMenu{
				{ID: 1, ParentID: 0, MenuName: "ValidRoot", Status: 1},
				{ID: 2, ParentID: 0, MenuName: "DisabledRoot", Status: 0},
				{ID: 3, ParentID: 999, MenuName: "UnmatchedParent", Status: 1},
				{ID: 4, ParentID: 1, MenuName: "ValidChild", Status: 1},
			},
			parentID:  0,
			wantCount: 1, // 只有一个有效的根节点
		},
	}

	// 执行测试用例
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BuildMenuTree(tt.menus, tt.parentID)

			// 基本数量验证
			if len(got) != tt.wantCount {
				t.Errorf("%s: 菜单数量不匹配, got = %v, want %v", tt.name, len(got), tt.wantCount)
				return
			}

			// 对于简单的测试用例，进行详细结构验证
			if tt.name == "TC001-空菜单列表" || tt.name == "TC002-无匹配父ID" {
				if len(got) != 0 {
					t.Errorf("%s: 应该返回空切片, got = %v", tt.name, got)
				}
			} else if tt.name == "TC003-状态禁用过滤" || tt.name == "TC004-单层菜单" {
				// 验证基本属性
				validateMenuStructure(t, tt.name, got, tt.want)
			}
		})
	}
}

// validateMenuStructure 验证菜单结构的基本属性
func validateMenuStructure(t *testing.T, testName string, got, want []*sysclient.Menu) {
	if len(got) != len(want) {
		t.Errorf("%s: 菜单数量不匹配, got = %v, want %v", testName, len(got), len(want))
		return
	}

	for i, gotMenu := range got {
		if i >= len(want) {
			break
		}
		wantMenu := want[i]

		if gotMenu.Id != wantMenu.Id {
			t.Errorf("%s: 菜单ID不匹配, index %d, got = %v, want %v", testName, i, gotMenu.Id, wantMenu.Id)
		}
		if gotMenu.ParentId != wantMenu.ParentId {
			t.Errorf("%s: 父级ID不匹配, index %d, got = %v, want %v", testName, i, gotMenu.ParentId, wantMenu.ParentId)
		}
		if gotMenu.MenuName != wantMenu.MenuName {
			t.Errorf("%s: 菜单名称不匹配, index %d, got = %v, want %v", testName, i, gotMenu.MenuName, wantMenu.MenuName)
		}
		// 验证Children字段已初始化
		if gotMenu.Children == nil {
			t.Errorf("%s: Children字段未初始化, index %d", testName, i)
		}
	}
}

// TestBuildMenuTree_Recursive 测试递归构建功能
func TestBuildMenuTree_Recursive(t *testing.T) {
	// 构造一个多层级的测试数据
	menus := []*model.SysMenu{
		{ID: 1, ParentID: 0, MenuName: "Level1-Root", Status: 1},
		{ID: 2, ParentID: 1, MenuName: "Level2-Child1", Status: 1},
		{ID: 3, ParentID: 1, MenuName: "Level2-Child2", Status: 1},
		{ID: 4, ParentID: 2, MenuName: "Level3-GrandChild", Status: 1},
		{ID: 5, ParentID: 4, MenuName: "Level4-GreatGrandChild", Status: 1},
	}

	result := BuildMenuTree(menus, 0)

	// 验证根节点
	if len(result) != 1 {
		t.Fatalf("应该只有一个根节点, got = %v", len(result))
	}

	root := result[0]
	if root.Id != 1 || root.MenuName != "Level1-Root" {
		t.Errorf("根节点信息错误, got ID=%d, Name=%s", root.Id, root.MenuName)
	}

	// 验证第二层
	if len(root.Children) != 2 {
		t.Errorf("第二层应该有两个子节点, got = %v", len(root.Children))
	}

	// 验证递归结构
	child1 := root.Children[0]
	if child1.Id != 2 || len(child1.Children) != 1 {
		t.Errorf("第一个子节点结构错误")
	}

	grandChild := child1.Children[0]
	if grandChild.Id != 4 || len(grandChild.Children) != 1 {
		t.Errorf("孙节点结构错误")
	}

	greatGrandChild := grandChild.Children[0]
	if greatGrandChild.Id != 5 || len(greatGrandChild.Children) != 0 {
		t.Errorf("曾孙节点结构错误")
	}
}

// TestBuildMenuTree_EdgeCases 边界情况测试
func TestBuildMenuTree_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		menus    []*model.SysMenu
		parentID int64
	}{
		{
			name:     "负数ParentID",
			menus:    []*model.SysMenu{{ID: 1, ParentID: -1, MenuName: "NegativeParent", Status: 1}},
			parentID: -1,
		},
		{
			name:     "大数值ID",
			menus:    []*model.SysMenu{{ID: 9999999999, ParentID: 0, MenuName: "LargeID", Status: 1}},
			parentID: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("%s: 函数执行出现panic: %v", tt.name, r)
				}
			}()

			got := BuildMenuTree(tt.menus, tt.parentID)
			if len(got) == 0 && len(tt.menus) > 0 && tt.menus[0].Status == 1 {
				t.Errorf("%s: 应该返回非空结果", tt.name)
			}
		})
	}
}

// BenchmarkBuildMenuTree 性能基准测试
func BenchmarkBuildMenuTree(b *testing.B) {
	// 创建大量测试数据
	var menus []*model.SysMenu
	for i := int64(0); i < 1000; i++ {
		status := int32(1)
		if i%10 == 0 {
			status = 0 // 10%的菜单是禁用状态
		}
		menus = append(menus, &model.SysMenu{
			ID:       i + 1,
			ParentID: i / 10, // 创建一些父子关系
			MenuName: "Menu" + string(rune(i)),
			Status:   status,
		})
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		BuildMenuTree(menus, 0)
	}
}
