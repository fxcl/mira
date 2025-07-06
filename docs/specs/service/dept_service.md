# Specification: DeptService Module

## 1. Overview

The `DeptService` module is responsible for managing department-related operations within the system. This includes creating, updating, deleting, and retrieving department information. It also provides utility functions for handling department hierarchies, such as building department trees and checking for child departments. The service interacts directly with the data access layer to perform database operations and applies data scope restrictions to ensure users can only access data they are permitted to see.

## 2. Dependencies

-   **`anima/dal`**: The Data Access Layer for database interactions (GORM).
-   **`app/dto`**: Data Transfer Objects for request and response structures.
-   **`app/model`**: Database entity models (e.g., `SysDept`, `SysRoleDept`).
-   **`common/types/constant`**: For system-wide constants like `NORMAL_STATUS`.
-   **`service/DataScope`**: A function (not shown in this file) used to apply role-based data access restrictions to queries.

## 3. Data Structures (Key DTOs)

-   **`dto.SaveDept`**: Used for both creating and updating departments. Contains all department fields.
-   **`dto.DeptListRequest`**: Used for querying a list of departments. Contains filter fields like `DeptName` and `Status`.
-   **`dto.DeptListResponse`**: The structure for each department in the returned list.
-   **`dto.DeptDetailResponse`**: The structure for a single, detailed department response.
-   **`dto.DeptTreeResponse`**: A simplified structure (`id`, `label`, `parent_id`) used for building department trees.
-   **`dto.SeleteTree`**: A structure (`id`, `label`, `parent_id`, `children`) used for building selectable dropdown trees.

## 4. Functions and Logic

### 4.1. `CreateDept`

Creates a new department record.

**Pseudocode:**

```
FUNCTION CreateDept(param: SaveDept_DTO):
  // TDD: Test case for successful creation with valid data.
  // TDD: Test case for creation failure due to database constraint (e.g., duplicate name if unique).
  // TDD: Test case for creation failure due to invalid ParentId.

  INSTANTIATE a new SysDept model from the `param` DTO.
  EXECUTE database INSERT operation with the new model.
  IF database operation returns an error:
    RETURN the error.
  END IF
  RETURN nil.
END FUNCTION
```

### 4.2. `UpdateDept`

Updates an existing department record based on its ID.

**Pseudocode:**

```
FUNCTION UpdateDept(param: SaveDept_DTO):
  // TDD: Test case for successful update of an existing department.
  // TDD: Test case for updating a non-existent department (should not error but affect 0 rows).
  // TDD: Test case for update failure due to database error.

  LOCATE SysDept record WHERE dept_id matches `param.DeptId`.
  EXECUTE database UPDATE operation on the located record with data from `param`.
  IF database operation returns an error:
    RETURN the error.
  END IF
  RETURN nil.
END FUNCTION
```

### 4.3. `DeleteDept`

Deletes a department record by its ID.

**Pseudocode:**

```
FUNCTION DeleteDept(deptId: integer):
  // TDD: Test case for successful deletion of an existing department.
  // TDD: Test case for deleting a department with child departments (should be handled by validation layer).
  // TDD: Test case for deleting a non-existent department.

  LOCATE SysDept record WHERE dept_id matches `deptId`.
  EXECUTE database DELETE operation.
  IF database operation returns an error:
    RETURN the error.
  END IF
  RETURN nil.
END FUNCTION
```

### 4.4. `GetDeptList`

Retrieves a filtered and permission-scoped list of departments.

**Pseudocode:**

```
FUNCTION GetDeptList(param: DeptListRequest_DTO, userId: integer):
  // TDD: Test case with no filters, should return all scoped departments.
  // TDD: Test case with DeptName filter.
  // TDD: Test case with Status filter.
  // TDD: Test case combining both filters.
  // TDD: Test case for a user with restricted data scope.

  INITIALIZE an empty list `depts` of type DeptListResponse.
  START a database query on the SysDept model.
  ORDER the query by `order_num`, then `dept_id`.
  APPLY data scope filter to the query using `GetDataScope("sys_dept", userId)`.

  IF `param.DeptName` is not empty:
    ADD a `WHERE dept_name LIKE '%...%'` condition to the query.
  END IF

  IF `param.Status` is not empty:
    ADD a `WHERE status = ?` condition to the query.
  END IF

  EXECUTE the query and populate the `depts` list.
  RETURN `depts`.
END FUNCTION
```

### 4.5. `GetDeptByDeptId`

Retrieves detailed information for a single department by its ID.

**Pseudocode:**

```
FUNCTION GetDeptByDeptId(deptId: integer):
  // TDD: Test case for a valid, existing department ID.
  // TDD: Test case for a non-existent department ID, should return an empty object.

  INITIALIZE an empty `dept` object of type DeptDetailResponse.
  QUERY the SysDept model for a record WHERE `dept_id` matches `deptId`.
  POPULATE the `dept` object with the first result found.
  RETURN `dept`.
END FUNCTION
```

### 4.6. `GetDeptByDeptName`

Retrieves detailed information for a single department by its name.

**Pseudocode:**

```
FUNCTION GetDeptByDeptName(deptName: string):
  // TDD: Test case for a valid, existing department name.
  // TDD: Test case for a non-existent department name.

  INITIALIZE an empty `dept` object of type DeptDetailResponse.
  QUERY the SysDept model for a record WHERE `dept_name` matches `deptName`.
  POPULATE the `dept` object with the first result found.
  RETURN `dept`.
END FUNCTION
```

### 4.7. `GetUserDeptTree`

Retrieves a permission-scoped, flat list of departments formatted for tree construction.

**Pseudocode:**

```
FUNCTION GetUserDeptTree(userId: integer):
  // TDD: Test case for a user that can see all departments.
  // TDD: Test case for a user with restricted department visibility.

  INITIALIZE an empty list `depts` of type DeptTreeResponse.
  QUERY the SysDept model.
  SELECT `dept_id` as `id`, `dept_name` as `label`, and `parent_id`.
  FILTER records WHERE `status` is `NORMAL_STATUS`.
  APPLY data scope filter using `GetDataScope("sys_dept", userId)`.
  ORDER the results by `order_num`, then `dept_id`.
  EXECUTE the query and populate the `depts` list.
  RETURN `depts`.
END FUNCTION
```

### 4.8. `GetDeptIdsByRoleId`

Gets a list of department IDs associated with a specific role.

**Pseudocode:**

```
FUNCTION GetDeptIdsByRoleId(roleId: integer):
  // TDD: Test case for a role with associated departments.
  // TDD: Test case for a role with no associated departments.
  // TDD: Test case where associated departments are disabled (status != NORMAL).

  INITIALIZE an empty integer slice `deptIds`.
  QUERY the SysRoleDept model.
  JOIN with SysDept on `sys_dept.dept_id = sys_role_dept.dept_id`.
  FILTER records WHERE `sys_dept.status` is `NORMAL_STATUS` AND `sys_role_dept.role_id` matches `roleId`.
  SELECT only the `sys_dept.dept_id` column into the `deptIds` slice.
  RETURN `deptIds`.
END FUNCTION
```

### 4.9. `DeptSelect` & `DeptSeleteToTree`

These two functions work together to provide a hierarchical department tree for UI dropdowns.

**`DeptSelect` Pseudocode:**

```
FUNCTION DeptSelect():
  // TDD: Test case to ensure only active departments are returned.
  // TDD: Test case to verify the sorting order.

  INITIALIZE an empty list `depts` of type SeleteTree.
  QUERY the SysDept model.
  SELECT `dept_id` as `id`, `dept_name` as `label`, and `parent_id`.
  FILTER records WHERE `status` is `NORMAL_STATUS`.
  ORDER the results by `order_num`, then `dept_id`.
  EXECUTE the query and populate the `depts` list.
  RETURN `depts`. // This is a flat list.
END FUNCTION
```

**`DeptSeleteToTree` Pseudocode:**

This is a recursive helper function to build a tree from a flat list.

```
FUNCTION DeptSeleteToTree(depts: list of SeleteTree, parentId: integer):
  // TDD: Test case for a simple two-level hierarchy.
  // TDD: Test case for a multi-level, complex hierarchy.
  // TDD: Test case with an empty input list.
  // TDD: Test case with a list that forms no valid tree from the given parentId.

  INITIALIZE an empty list `tree` of type SeleteTree.
  FOR each `dept` in the `depts` list:
    IF `dept.ParentId` matches the `parentId` parameter:
      // This dept is a direct child of the current node being processed.
      // Find all children of this current dept.
      children = RECURSIVE_CALL DeptSeleteToTree(depts, dept.Id).
      CREATE a new SeleteTree node using the current `dept`'s data and the `children` list.
      APPEND the new node to the `tree`.
    END IF
  END FOR
  RETURN `tree`.
END FUNCTION
```

### 4.10. `DeptHasChildren`

Checks if a given department has any child departments.

**Pseudocode:**

```
FUNCTION DeptHasChildren(deptId: integer):
  // TDD: Test case for a department with children, should return true.
  // TDD: Test case for a department with no children, should return false.

  QUERY the SysDept model.
  COUNT records WHERE `parent_id` matches `deptId`.
  RETURN `true` if the count is greater than 0, otherwise `false`.
END FUNCTION
