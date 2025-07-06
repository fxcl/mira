# Spec: Role Validator (`app/validator/role_validator.go`)

This document specifies the validation logic for role-related operations.

## Module: `validator`

### Function: `CreateRoleValidator`

Validates the request payload for creating a new role.

**Signature:**
`func CreateRoleValidator(param dto.CreateRoleRequest) error`

**Pseudocode:**
```
FUNCTION CreateRoleValidator(request):
  IF request.RoleName is empty THEN
    RETURN error "please enter the role name"
  ENDIF

  IF request.RoleKey is empty THEN
    RETURN error "please enter the permission string"
  ENDIF

  RETURN nil
ENDFUNCTION
```

**TDD Anchors:**
- `TestCreateRole_Failure_EmptyRoleName`: Should fail when `RoleName` is an empty string.
- `TestCreateRole_Failure_EmptyRoleKey`: Should fail when `RoleKey` is an empty string.
- `TestCreateRole_Success`: Should pass when `RoleName` and `RoleKey` are provided.

---

### Function: `UpdateRoleValidator`

Validates the request payload for updating an existing role.

**Signature:**
`func UpdateRoleValidator(param dto.UpdateRoleRequest) error`

**Pseudocode:**
```
FUNCTION UpdateRoleValidator(request):
  IF request.RoleId <= 0 THEN
    RETURN error "parameter error"
  ENDIF

  IF request.RoleName is empty THEN
    RETURN error "please enter the role name"
  ENDIF

  IF request.RoleKey is empty THEN
    RETURN error "please enter the permission string"
  ENDIF

  RETURN nil
ENDFUNCTION
```

**TDD Anchors:**
- `TestUpdateRole_Failure_InvalidRoleId`: Should fail when `RoleId` is zero or negative.
- `TestUpdateRole_Failure_EmptyRoleName`: Should fail when `RoleName` is an empty string.
- `TestUpdateRole_Failure_EmptyRoleKey`: Should fail when `RoleKey` is an empty string.
- `TestUpdateRole_Success`: Should pass with valid `RoleId`, `RoleName`, and `RoleKey`.

---

### Function: `RemoveRoleValidator`

Validates the request to remove one or more roles.

**Signature:**
`func RemoveRoleValidator(roleIds []int, roleId int, roleName string) error`

**Pseudocode:**
```
FUNCTION RemoveRoleValidator(ids_to_remove, current_user_role_id, current_user_role_name):
  IF ids_to_remove contains 1 (Super Admin ID) THEN
    RETURN error "the super administrator cannot be deleted"
  ENDIF

  IF ids_to_remove contains current_user_role_id THEN
    RETURN error "the {current_user_role_name} role cannot be deleted"
  ENDIF

  RETURN nil
ENDFUNCTION
```

**TDD Anchors:**
- `TestRemoveRole_Failure_AttemptToDeleteSuperAdmin`: Should fail when the list of IDs to remove contains `1`.
- `TestRemoveRole_Failure_AttemptToDeleteOwnRole`: Should fail when the list of IDs to remove contains the current user's own `roleId`.
- `TestRemoveRole_Success`: Should pass when attempting to delete other valid roles.

---

### Function: `ChangeRoleStatusValidator`

Validates the request to change a role's status.

**Signature:**
`func ChangeRoleStatusValidator(param dto.UpdateRoleRequest) error`

**Pseudocode:**
```
FUNCTION ChangeRoleStatusValidator(request):
  IF request.RoleId <= 0 THEN
    RETURN error "parameter error"
  ENDIF

  IF request.Status is empty THEN
    RETURN error "please select a status"
  ENDIF

  RETURN nil
ENDFUNCTION
```

**TDD Anchors:**
- `TestChangeStatus_Failure_InvalidRoleId`: Should fail when `RoleId` is zero or negative.
- `TestChangeStatus_Failure_EmptyStatus`: Should fail when `Status` is an empty string.
- `TestChangeStatus_Success`: Should pass with a valid `RoleId` and `Status`.
