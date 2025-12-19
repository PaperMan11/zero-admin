package mockdb

import (
	"context"
	"sync"

	"zero-admin/rpc/sys/db/mysql/model"
)

type MockDB struct {
	mu sync.RWMutex

	users         map[int64]*model.SysUser
	usernames     map[string]int64
	roles         map[int64]*model.SysRole
	roleNames     map[string]int64
	roleCodes     map[string]int64
	userRoles     map[int64][]string               // userID -> []roleCode
	roleScopes    map[string][]*model.SysRoleScope // roleCode -> scopes
	menus         map[int64]*model.SysMenu
	scopes        map[int64]*model.SysScope
	scopeCodes    map[string]int64
	loginLogs     []*model.SysLoginLog
	operationLogs []*model.SysOperateLog
}

func NewMockDB() (*MockDB, error) {
	return &MockDB{
		users:      make(map[int64]*model.SysUser),
		usernames:  make(map[string]int64),
		roles:      make(map[int64]*model.SysRole),
		roleNames:  make(map[string]int64),
		roleCodes:  make(map[string]int64),
		userRoles:  make(map[int64][]string),
		roleScopes: make(map[string][]*model.SysRoleScope),
		menus:      make(map[int64]*model.SysMenu),
		scopes:     make(map[int64]*model.SysScope),
		scopeCodes: make(map[string]int64),
	}, nil
}

// ---------------------用户 & 角色---------------------

func (m *MockDB) CreateUser(ctx context.Context, user model.SysUser) (int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if user.ID == 0 {
		// Generate a simple ID based on map size
		user.ID = int64(len(m.users) + 1)
	}

	m.users[user.ID] = &user
	m.usernames[user.Username] = user.ID

	return user.ID, nil
}

func (m *MockDB) GetUserByUsername(ctx context.Context, username string) (*model.SysUser, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if id, exists := m.usernames[username]; exists {
		if user, exists := m.users[id]; exists {
			return user, nil
		}
	}

	return nil, nil
}

func (m *MockDB) GetUserByID(ctx context.Context, userID int64) (*model.SysUser, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if user, exists := m.users[userID]; exists {
		return user, nil
	}

	return nil, nil
}

func (m *MockDB) UpdateUserByID(ctx context.Context, userID int64, updates interface{}) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if user, exists := m.users[userID]; exists {
		// In a real implementation, we would apply the updates
		// For mock, we'll just confirm the user exists
		_ = user
		return nil
	}

	return nil
}

func (m *MockDB) AddUserRolesTx(ctx context.Context, userID int64, roleCodes []string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.userRoles[userID] = roleCodes
	return nil
}

func (m *MockDB) DeleteUserTx(ctx context.Context, userID int64) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if user, exists := m.users[userID]; exists {
		delete(m.usernames, user.Username)
		delete(m.users, userID)
		delete(m.userRoles, userID)
	}

	return nil
}

func (m *MockDB) GetUsersPagination(ctx context.Context, status int32, page, pageSize int) ([]*model.SysUser, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []*model.SysUser
	count := 0
	skip := (page - 1) * pageSize

	for _, user := range m.users {
		if status == 0 || user.Status == status {
			if count >= skip && len(result) < pageSize {
				result = append(result, user)
			}
			count++
		}
	}

	return result, nil
}

func (m *MockDB) CountUsers(ctx context.Context, status int32) (int64, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	count := int64(0)
	for _, user := range m.users {
		if status == 0 || user.Status == status {
			count++
		}
	}

	return count, nil
}

func (m *MockDB) SaveUser(ctx context.Context, user model.SysUser) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if user.ID == 0 {
		user.ID = int64(len(m.users) + 1)
		m.users[user.ID] = &user
		m.usernames[user.Username] = user.ID
	} else {
		if oldUser, exists := m.users[user.ID]; exists {
			delete(m.usernames, oldUser.Username)
		}
		m.users[user.ID] = &user
		m.usernames[user.Username] = user.ID
	}

	return nil
}

func (m *MockDB) CreateRole(ctx context.Context, role model.SysRole) (int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if role.ID == 0 {
		role.ID = int64(len(m.roles) + 1)
	}

	m.roles[role.ID] = &role
	m.roleNames[role.RoleName] = role.ID
	m.roleCodes[role.RoleCode] = role.ID

	return role.ID, nil
}

func (m *MockDB) DeleteRoleTx(ctx context.Context, roleCode string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if id, exists := m.roleCodes[roleCode]; exists {
		if role, exists := m.roles[id]; exists {
			delete(m.roleNames, role.RoleName)
			delete(m.roleCodes, role.RoleCode)
			delete(m.roles, id)
			delete(m.roleScopes, roleCode)

			// Remove this role from all users
			for userID, roles := range m.userRoles {
				newRoles := []string{}
				for _, rc := range roles {
					if rc != roleCode {
						newRoles = append(newRoles, rc)
					}
				}
				m.userRoles[userID] = newRoles
			}
		}
	}

	return nil
}

func (m *MockDB) DeleteRoleScopes(ctx context.Context, roleCode string, scopeCodes []string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if scopes, exists := m.roleScopes[roleCode]; exists {
		var newScopes []*model.SysRoleScope
		for _, scope := range scopes {
			shouldDelete := false
			for _, code := range scopeCodes {
				if scope.ScopeCode == code {
					shouldDelete = true
					break
				}
			}
			if !shouldDelete {
				newScopes = append(newScopes, scope)
			}
		}
		m.roleScopes[roleCode] = newScopes
	}

	return nil
}

func (m *MockDB) GetRoleByID(ctx context.Context, roleID int64) (*model.SysRole, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if role, exists := m.roles[roleID]; exists {
		return role, nil
	}

	return nil, nil
}

func (m *MockDB) GetAllRoles(ctx context.Context) ([]*model.SysRole, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var result []*model.SysRole
	for _, role := range m.roles {
		result = append(result, role)
	}
	return result, nil
}

func (m *MockDB) GetRoleByIDs(ctx context.Context, roleIDs []int64) ([]*model.SysRole, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var result []*model.SysRole
	for _, id := range roleIDs {
		if role, exists := m.roles[id]; exists {
			result = append(result, role)
		}
	}
	return result, nil
}

func (m *MockDB) GetRoleByCode(ctx context.Context, roleCode string) (*model.SysRole, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.roles[m.roleCodes[roleCode]], nil
}

func (m *MockDB) GetRoleByCodes(ctx context.Context, roleCodes []string) ([]*model.SysRole, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var result []*model.SysRole
	for _, code := range roleCodes {
		if role, exists := m.roles[m.roleCodes[code]]; exists {
			result = append(result, role)
		}
	}
	return result, nil
}

func (m *MockDB) GetRoleByName(ctx context.Context, roleName string) (*model.SysRole, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if id, exists := m.roleNames[roleName]; exists {
		if role, exists := m.roles[id]; exists {
			return role, nil
		}
	}

	return nil, nil
}

func (m *MockDB) ExistsRoleByName(ctx context.Context, roleName string) (bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	_, exists := m.roleNames[roleName]
	return exists, nil
}

func (m *MockDB) ExistsRoleByCode(ctx context.Context, roleCode string) (bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	_, exists := m.roleCodes[roleCode]
	return exists, nil
}

func (m *MockDB) ExistsRoleByID(ctx context.Context, roleID int64) (bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	_, exists := m.roles[roleID]
	return exists, nil
}

func (m *MockDB) GetRolesByUserID(ctx context.Context, userID int64) ([]*model.SysRole, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []*model.SysRole
	if roleCodes, exists := m.userRoles[userID]; exists {
		for _, code := range roleCodes {
			if id, exists := m.roleCodes[code]; exists {
				if role, exists := m.roles[id]; exists {
					result = append(result, role)
				}
			}
		}
	}

	return result, nil
}

func (m *MockDB) GetRolesPagination(ctx context.Context, status int32, page, pageSize int) ([]*model.SysRole, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []*model.SysRole
	count := 0
	skip := (page - 1) * pageSize

	for _, role := range m.roles {
		if status == 0 || role.Status == status {
			if count >= skip && len(result) < pageSize {
				result = append(result, role)
			}
			count++
		}
	}

	return result, nil
}

func (m *MockDB) CountRoles(ctx context.Context, status int32) (int64, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	count := int64(0)
	for _, role := range m.roles {
		if status == 0 || role.Status == status {
			count++
		}
	}

	return count, nil
}

func (m *MockDB) CountUserRoles(ctx context.Context, roleCode string) (int64, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	count := int64(0)
	for _, roles := range m.userRoles {
		for _, code := range roles {
			if code == roleCode {
				count++
				break
			}
		}
	}

	return count, nil
}

func (m *MockDB) GetUserRoleCodes(ctx context.Context, userID int64) ([]string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if roleCodes, exists := m.userRoles[userID]; exists {
		return roleCodes, nil
	}

	return []string{}, nil
}

func (m *MockDB) SaveRole(ctx context.Context, role model.SysRole) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if role.ID == 0 {
		role.ID = int64(len(m.roles) + 1)
		m.roles[role.ID] = &role
		m.roleNames[role.RoleName] = role.ID
		m.roleCodes[role.RoleCode] = role.ID
	} else {
		if oldRole, exists := m.roles[role.ID]; exists {
			delete(m.roleNames, oldRole.RoleName)
			delete(m.roleCodes, oldRole.RoleCode)
		}
		m.roles[role.ID] = &role
		m.roleNames[role.RoleName] = role.ID
		m.roleCodes[role.RoleCode] = role.ID
	}

	return nil
}

func (m *MockDB) UpdateRoleScopesTx(ctx context.Context, roleCode string, roleScopes []model.SysRoleScope) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	scopes := make([]*model.SysRoleScope, len(roleScopes))
	for i, rs := range roleScopes {
		rsCopy := rs
		scopes[i] = &rsCopy
	}
	m.roleScopes[roleCode] = scopes

	return nil
}

func (m *MockDB) AddRoleScopes(ctx context.Context, roleScopes []*model.SysRoleScope) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, rs := range roleScopes {
		if scopes, exists := m.roleScopes[rs.RoleCode]; exists {
			m.roleScopes[rs.RoleCode] = append(scopes, rs)
		} else {
			m.roleScopes[rs.RoleCode] = []*model.SysRoleScope{rs}
		}
	}

	return nil
}

func (m *MockDB) ToggleRoleStatus(ctx context.Context, roleID int64, status int32, operator string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if role, exists := m.roles[roleID]; exists {
		role.Status = status
		role.Updater = operator
	}

	return nil
}

// ---------------------菜单 & 权限---------------------

func (m *MockDB) GetMenus(ctx context.Context, status int32, page, pageSize int) ([]*model.SysMenu, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []*model.SysMenu
	count := 0
	skip := (page - 1) * pageSize

	for _, menu := range m.menus {
		if status == 0 || menu.Status == status {
			if count >= skip && len(result) < pageSize {
				result = append(result, menu)
			}
			count++
		}
	}

	return result, nil
}

func (m *MockDB) GetUnassignedMenus(ctx context.Context) ([]*model.SysMenu, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	res := make([]*model.SysMenu, 0)
	for _, menu := range m.menus {
		if menu.ScopeID == 0 {
			res = append(res, menu)
		}
	}
	return res, nil
}

func (m *MockDB) GetAllMenus(ctx context.Context) ([]*model.SysMenu, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	menus := make([]*model.SysMenu, 0, len(m.menus))
	for _, menu := range m.menus {
		menus = append(menus, menu)
	}
	return menus, nil
}

func (m *MockDB) GetMenuByID(ctx context.Context, menuID int64) (*model.SysMenu, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if menu, exists := m.menus[menuID]; exists {
		return menu, nil
	}

	return nil, nil
}

func (m *MockDB) GetMenusByRoles(ctx context.Context, roleCodes []string) ([]*model.SysMenu, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// In a real implementation, this would filter menus by role permissions
	// For mock, we'll return all menus
	var result []*model.SysMenu
	for _, menu := range m.menus {
		result = append(result, menu)
	}

	return result, nil
}

func (m *MockDB) GetMenusByScopeIDs(ctx context.Context, scopeIDs []int64) ([]*model.SysMenu, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []*model.SysMenu
	for _, menu := range m.menus {
		for _, scopeID := range scopeIDs {
			if menu.ScopeID == scopeID {
				result = append(result, menu)
				break
			}
		}
	}

	return result, nil
}

func (m *MockDB) CreateMenu(ctx context.Context, menu model.SysMenu) (int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if menu.ID == 0 {
		menu.ID = int64(len(m.menus) + 1)
	}

	m.menus[menu.ID] = &menu

	return menu.ID, nil
}

func (m *MockDB) CreateMenus(ctx context.Context, menus []*model.SysMenu) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, menu := range menus {
		if menu.ID == 0 {
			menu.ID = int64(len(m.menus) + 1)
		}
		m.menus[menu.ID] = menu
	}

	return nil
}

func (m *MockDB) DeleteMenu(ctx context.Context, menuID int64) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.menus, menuID)
	return nil
}

func (m *MockDB) UpdateMenu(ctx context.Context, menuID int64, updates interface{}) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// In a real implementation, this would apply updates to the menu
	// For mock, we'll just confirm the menu exists
	_, exists := m.menus[menuID]
	if !exists {
		return nil
	}

	return nil
}

func (m *MockDB) SaveMenu(ctx context.Context, menu model.SysMenu) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if menu.ID == 0 {
		menu.ID = int64(len(m.menus) + 1)
		m.menus[menu.ID] = &menu
	} else {
		m.menus[menu.ID] = &menu
	}

	return nil
}

func (m *MockDB) GetMenusByScopeID(ctx context.Context, scopeID int64) ([]*model.SysMenu, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// For mock, return all menus
	var result []*model.SysMenu
	for _, menu := range m.menus {
		result = append(result, menu)
	}

	return result, nil
}

func (m *MockDB) ExistsMenuByName(ctx context.Context, menuName string) (bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, menu := range m.menus {
		if menu.MenuName == menuName {
			return true, nil
		}
	}

	return false, nil
}

func (m *MockDB) ExistsMenuByPath(ctx context.Context, menuPath string) (bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, menu := range m.menus {
		if menu.Path == menuPath {
			return true, nil
		}
	}

	return false, nil
}

func (m *MockDB) ExistsMenu(ctx context.Context, menuID int64) (bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	_, exists := m.menus[menuID]
	return exists, nil
}

func (m *MockDB) CreateScope(ctx context.Context, scope model.SysScope) (int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if scope.ID == 0 {
		scope.ID = int64(len(m.scopes) + 1)
	}

	m.scopes[scope.ID] = &scope
	m.scopeCodes[scope.ScopeCode] = scope.ID

	return scope.ID, nil
}

func (m *MockDB) CreateScopeTx(ctx context.Context, scope model.SysScope, menuIDs []int64) (int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if scope.ID == 0 {
		scope.ID = int64(len(m.scopes) + 1)
	}

	m.scopes[scope.ID] = &scope
	m.scopeCodes[scope.ScopeCode] = scope.ID

	// In a real implementation, we would associate the scope with menus
	// For mock, we'll just ignore the menuIDs

	return scope.ID, nil
}

func (m *MockDB) ExistsScope(ctx context.Context, scopeID int64) (bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	_, exists := m.scopes[scopeID]
	return exists, nil
}

func (m *MockDB) ExistsScopeByCode(ctx context.Context, scopeCode string) (bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	_, exists := m.scopeCodes[scopeCode]
	return exists, nil
}

func (m *MockDB) CountScopes(ctx context.Context) (int64, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return int64(len(m.scopes)), nil
}

func (m *MockDB) SaveScope(ctx context.Context, scope model.SysScope) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if scope.ID == 0 {
		scope.ID = int64(len(m.scopes) + 1)
		m.scopes[scope.ID] = &scope
		m.scopeCodes[scope.ScopeCode] = scope.ID
	} else {
		if oldScope, exists := m.scopes[scope.ID]; exists {
			delete(m.scopeCodes, oldScope.ScopeCode)
		}
		m.scopes[scope.ID] = &scope
		m.scopeCodes[scope.ScopeCode] = scope.ID
	}

	return nil
}

func (m *MockDB) GetScopeByID(ctx context.Context, scopeID int64) (*model.SysScope, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if scope, exists := m.scopes[scopeID]; exists {
		return scope, nil
	}

	return nil, nil
}

func (m *MockDB) GetScopeByCode(ctx context.Context, scopeCode string) (*model.SysScope, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.scopes[m.scopeCodes[scopeCode]], nil
}

func (m *MockDB) GetScopes(ctx context.Context, scopeIDs []int64) ([]*model.SysScope, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	scopes := make([]*model.SysScope, 0)

	for _, id := range scopeIDs {
		for _, scope := range m.scopes {
			if scope.ID == id {
				scopes = append(scopes, scope)
			}
		}
	}

	return scopes, nil
}

func (m *MockDB) GetAllScopes(ctx context.Context) ([]*model.SysScope, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var result []*model.SysScope
	for _, scope := range m.scopes {
		result = append(result, scope)
	}
	return result, nil
}

func (m *MockDB) GetScopesByCodes(ctx context.Context, scopeCodes []string) ([]*model.SysScope, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	scopes := make([]*model.SysScope, 0)
	for _, code := range scopeCodes {
		for _, scope := range m.scopes {
			if scope.ScopeCode == code {
				scopes = append(scopes, scope)
			}
		}
	}
	return scopes, nil
}

func (m *MockDB) UpsertRoleScopes(ctx context.Context, roleScope model.SysRoleScope) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	roleScopes, ok := m.roleScopes[roleScope.RoleCode]
	if !ok {
		roleScopes = append(roleScopes, &roleScope)
	} else {
		for _, v := range roleScopes {
			if v.ScopeCode == roleScope.ScopeCode {
				roleScope.Perm = v.Perm
			}
		}
	}
	return nil
}

func (m *MockDB) GetScopesPagination(ctx context.Context, status int32, page, pageSize int) ([]*model.SysScope, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []*model.SysScope
	count := 0
	skip := (page - 1) * pageSize

	for _, scope := range m.scopes {
		if count >= skip && len(result) < pageSize {
			result = append(result, scope)
		}
		count++
	}

	return result, nil
}

func (m *MockDB) ToggleScopeStatus(ctx context.Context, scopeID int64, status int32, operator string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if scope, exists := m.scopes[scopeID]; exists {
		scope.Status = status
		scope.Updater = operator
	}
	return nil
}

func (m *MockDB) AddScopeMenus(ctx context.Context, scopeID int64, menus []int64) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// In a real implementation, this would associate the scope with menus
	// For mock, we'll just confirm the scope exists
	_, exists := m.scopes[scopeID]
	if !exists {
		return nil
	}

	return nil
}

func (m *MockDB) DeleteScopeTx(ctx context.Context, scopeID int64) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.scopes[scopeID]; exists {
		delete(m.scopes, scopeID)
		for menuID, menu := range m.menus {
			if menu.ParentID == scopeID {
				delete(m.menus, menuID)
			}
		}
	}

	return nil
}

func (m *MockDB) DeleteScopeMenus(ctx context.Context, scopeID int64) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// In a real implementation, this would remove menu associations
	// For mock, we'll just confirm the scope exists
	_, exists := m.scopes[scopeID]
	if !exists {
		return nil
	}

	return nil
}

func (m *MockDB) GetRoleScopesPerm(ctx context.Context, roleCode string) ([]model.RoleScopeInfo, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Return empty slice for mock
	return []model.RoleScopeInfo{}, nil
}

func (m *MockDB) GetRolesScopesPerm(ctx context.Context, roleCodes []string) ([]model.RoleScopeInfo, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Return empty slice for mock
	return []model.RoleScopeInfo{}, nil
}

func (m *MockDB) GetRolesByScopeCode(ctx context.Context, scopeCode string) ([]*model.SysRole, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	res := make([]*model.SysRole, 0)
	for _, roleScope := range m.roleScopes[scopeCode] {
		if roleScope.ScopeCode == scopeCode {
			for _, role := range m.roles {
				if role.RoleCode == roleScope.RoleCode {
					res = append(res, role)
				}
			}
		}
	}
	return res, nil
}

func (m *MockDB) GetRolePermsByScopeCode(ctx context.Context, scopeCode string) ([]*model.SysRoleScope, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	res := make([]*model.SysRoleScope, 0)
	for _, roleScope := range m.roleScopes[scopeCode] {
		if roleScope.ScopeCode == scopeCode {
			res = append(res, roleScope)
		}
	}
	return res, nil
}

// 全量更新安全范围的菜单树
func (m *MockDB) UpdateScopeMenusTx(ctx context.Context, scopeID int64, menus []*model.SysMenu) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	return nil
}

// ---------------------登录日志 & 操作日志---------------------

func (m *MockDB) CreateLoginLog(ctx context.Context, log model.SysLoginLog) (int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	log.ID = int64(len(m.loginLogs) + 1)
	m.loginLogs = append(m.loginLogs, &log)

	return log.ID, nil
}

func (m *MockDB) CreateOperationLog(ctx context.Context, log model.SysOperateLog) (int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	log.ID = int64(len(m.operationLogs) + 1)
	m.operationLogs = append(m.operationLogs, &log)

	return log.ID, nil
}

func (m *MockDB) CreateOperationLogs(ctx context.Context, logs []*model.SysOperateLog) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, log := range logs {
		log.ID = int64(len(m.operationLogs) + 1)
		m.operationLogs = append(m.operationLogs, log)
	}

	return nil
}

// 获取操作日志详情
func (m *MockDB) GetOperateLog(ctx context.Context, logID int64) (*model.SysOperateLog, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.operationLogs[logID-1], nil
}

// 获取操作日志列表
func (m *MockDB) GetOperateLogs(ctx context.Context, filter model.OperateLogFilter, page int, pageSize int) ([]*model.SysOperateLog, int64, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.operationLogs, int64(len(m.operationLogs)), nil
}

// 删除操作日志
func (m *MockDB) DeleteOperateLogs(ctx context.Context, logIDs []int64) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, logID := range logIDs {
		for i, log := range m.operationLogs {
			if log.ID == logID {
				m.operationLogs = append(m.operationLogs[:i], m.operationLogs[i+1:]...)
				break
			}
		}
	}
	return nil
}

func (m *MockDB) DeleteOperateLog(ctx context.Context, logID int64) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for i, log := range m.operationLogs {
		if log.ID == logID {
			m.operationLogs = append(m.operationLogs[:i], m.operationLogs[i+1:]...)
		}
	}
	return nil
}
