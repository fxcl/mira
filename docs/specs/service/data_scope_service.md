# Specification: Data Scope Service Module

## 1. Overview

The Data Scope Service provides a centralized function, `GetDataScope`, for applying row-level data permissions to database queries. It is designed to be used as a GORM Scope, dynamically constructing SQL `WHERE` clauses based on the roles and permissions of the currently authenticated user. This ensures that users can only access data they are authorized to see, such as data from their own department, their department and its sub-departments, or only their own personal data.

The core principle is to inspect the user's assigned roles, check the `DataScope` setting for each role, and build a compound `OR` condition that combines the permissions from all applicable roles.

## 2. Dependencies

-   **`gorm.io/gorm`**: The GORM library for database interaction.
-   **`app/service/UserService`**: To fetch details of the current user (e.g., their department ID).
-   **`app/service/RoleService`**: To fetch the list of roles assigned to the current user.
-   **`common/types/constant`**: For system-wide constants like `NORMAL_STATUS`.

## 3. Function: `GetDataScope`

This is the primary and only function in the module. It returns a `func(*gorm.DB) *gorm.DB`, which is a GORM Scope.

### 3.1. Parameters

-   **`deptAlias` (string)**: The SQL alias for the department table (or the table containing `dept_id`) in the target query. Defaults to `"sys_dept"` if an empty string is provided. This is crucial for constructing correct SQL joins and conditions.
-   **`userId` (int)**: The ID of the user for whom the data scope is being calculated.
-   **`userAlias` (string)**: The SQL alias for the user table in the target query. This is only required when applying the "Personal data only" scope (`DataScope == "5"`).

### 3.2. Return Value

-   **`func(*gorm.DB) *gorm.DB`**: A function that GORM can use in its `.Scopes()` method to apply the calculated data permission filters to a query.

## 4. Core Logic Flow

**Pseudocode:**

```
FUNCTION GetDataScope(deptAlias: string, userId: integer, userAlias: string):
  // TDD: Test case for super administrator (userId=1), should apply no filters.
  // TDD: Test case for a user with a role that has "All" data scope, should apply no filters.
  // TDD: Test case for a user with one role with "Department" scope.
  // TDD: Test case for a user with one role with "Department and Sub-department" scope.
  // TDD: Test case for a user with one role with "Custom" scope.
  // TDD: Test case for a user with one role with "Personal" scope (with and without userAlias).
  // TDD: Test case for a user with multiple roles combining different scopes (e.g., "Department" OR "Custom").
  // TDD: Test case where a user has no roles with data scopes, should not apply filters.

  // 1. Handle Super Administrator
  IF userId is 1 (Super Admin):
    RETURN a function that does nothing to the query (returns the DB object as is).
  END IF

  // 2. Set default alias
  IF deptAlias is empty:
    SET deptAlias = "sys_dept".
  END IF

  // 3. Fetch User and Role Data
  FETCH user details using UserService for the given userId.
  FETCH the list of roles for the user using RoleService.

  // 4. Prepare for Custom Scope
  INITIALIZE an empty list `customRoleIds`.
  LOOP through the user's roles:
    IF a role has DataScope "2" (Custom) and is active:
      ADD its RoleId to `customRoleIds`.
    END IF
  END LOOP

  // 5. Return the GORM Scope Function
  RETURN a function that takes a GORM DB object:
    INITIALIZE empty list `sqlCondition` for WHERE clauses.
    INITIALIZE empty list `sqlArg` for query parameters.

    LOOP through each of the user's roles:
      SWITCH on the role's DataScope:

        // CASE "1": All Permissions
        // If any role has this, immediately return the DB object without any filters.
        CASE "1":
          RETURN the DB object as is.

        // CASE "2": Custom Permissions
        // Adds a condition to allow access to departments linked to the user's custom-scoped roles.
        CASE "2":
          IF `customRoleIds` is not empty:
            ADD condition: `deptAlias.dept_id IN (SELECT dept_id FROM sys_role_dept WHERE role_id IN (?))`
            ADD `customRoleIds` to `sqlArg`.
          ELSE:
            // Fallback for a single role, though the pre-calculated list is more efficient.
            ADD condition: `deptAlias.dept_id IN (SELECT dept_id FROM sys_role_dept WHERE role_id = ?)`
            ADD current `role.RoleId` to `sqlArg`.
          END IF

        // CASE "3": Department Only
        // Restricts access to the user's own department.
        CASE "3":
          ADD condition: `deptAlias.dept_id = ?`
          ADD `user.DeptId` to `sqlArg`.

        // CASE "4": Department and Sub-departments
        // Restricts access to the user's department and all its children.
        // This relies on an `ancestors` column in the departments table.
        CASE "4":
          ADD condition: `deptAlias.dept_id IN (SELECT dept_id FROM sys_dept WHERE dept_id = ? OR find_in_set(?, ancestors))`
          ADD `user.DeptId` to `sqlArg` twice.

        // CASE "5": Personal Data Only
        // Restricts access to records created by or assigned to the user.
        CASE "5":
          IF `userAlias` is provided:
            ADD condition: `userAlias.user_id = ?`
            ADD `user.UserId` to `sqlArg`.
          ELSE:
            // If no user alias is provided, this scope cannot be applied.
            // Add a condition that will always be false to prevent data leakage.
            ADD condition: `deptAlias.dept_id = ?`
            ADD `0` to `sqlArg`.
          END IF
      END SWITCH
    END LOOP

    // 6. Apply the final condition
    IF `sqlCondition` list has items:
      JOIN all conditions with " OR ".
      APPLY the combined WHERE clause to the DB object.
    END IF

    RETURN the modified DB object.
  END FUNCTION
```

## 5. Usage Example

To use the data scope, you apply it to a GORM query chain using the `.Scopes()` method.

```go
// Example: Fetching a list of users with data permissions applied.
// Here, the user table is aliased as "u" and the department table as "d".
func GetScopedUserList(db *gorm.DB, currentUserID int) {
    var users []model.User

    db.Model(&model.User{}).
       Alias("u").
       Joins("LEFT JOIN sys_dept d ON u.dept_id = d.dept_id").
       Scopes(GetDataScope("d", currentUserID, "u")). // Apply the scope
       Find(&users)

    // The 'users' slice will now only contain users that the 'currentUserID' is allowed to see.
}
