# Specification: Config Service

## 1. Overview

The `ConfigService` is responsible for managing system configuration parameters. It provides standard CRUD (Create, Read, Update, Delete) functionalities for configuration settings. Additionally, it incorporates a caching layer using Redis to optimize read operations for frequently accessed configurations, thereby reducing database load.

## 2. Data Structures

### Input DTOs

-   **`dto.SaveConfig`**: Used for creating and updating configurations.
    -   `ConfigId` (int, optional for create)
    -   `ConfigName` (string)
    -   `ConfigKey` (string)
    -   `ConfigValue` (string)
    -   `ConfigType` (string)
    -   `CreateBy` (string, for create)
    -   `UpdateBy` (string, for update)
    -   `Remark` (string)

-   **`dto.ConfigListRequest`**: Used for querying a list of configurations.
    -   `ConfigName` (string, optional filter)
    -   `ConfigKey` (string, optional filter)
    -   `ConfigType` (string, optional filter)
    -   `BeginTime` (string, optional filter)
    -   `EndTime` (string, optional filter)
    -   `PageNum` (int, for pagination)
    -   `PageSize` (int, for pagination)

### Output DTOs

-   **`dto.ConfigListResponse`**: Represents a single item in the configuration list.
    -   *(Fields inferred from `model.SysConfig` and context)*

-   **`dto.ConfigDetailResponse`**: Represents the detailed view of a single configuration.
    -   *(Fields inferred from `model.SysConfig` and context)*

### Database Model

-   **`model.SysConfig`**: The GORM model for the `sys_config` table.
    -   `ConfigId`
    -   `ConfigName`
    -   `ConfigKey`
    -   `ConfigValue`
    -   `ConfigType`
    -   `CreateBy`
    -   `UpdateBy`
    -   `Remark`
    -   `CreateTime`
    -   `UpdateTime`

## 3. Modules & Pseudocode

### `ConfigService` Struct

A container for configuration-related business logic.

```go
type ConfigService struct{}
```

---

### `CreateConfig`

Creates a new configuration parameter.

**Pseudocode:**

```plaintext
FUNCTION CreateConfig(param: SaveConfig) -> ERROR:
  VALIDATE input `param` for required fields (ConfigName, ConfigKey, etc.) // TDD Anchor
  INITIALIZE a new `SysConfig` model object from `param`.
  CALL database to CREATE the new record.
  IF database error:
    RETURN error.
  ENDIF
  RETURN nil.
END FUNCTION
```

**TDD Anchors:**
-   `TestCreateConfig_Success`: Verify a config is created with valid data.
-   `TestCreateConfig_MissingRequiredFields`: Ensure an error is returned for invalid input.
-   `TestCreateConfig_DuplicateKey`: Ensure database constraints (if any) for unique `ConfigKey` are handled.

---

### `UpdateConfig`

Updates an existing configuration parameter based on its ID.

**Pseudocode:**

```plaintext
FUNCTION UpdateConfig(param: SaveConfig) -> ERROR:
  VALIDATE input `param` for `ConfigId`. // TDD Anchor
  FIND `SysConfig` record by `param.ConfigId`.
  IF not found:
    RETURN "not found" error. // TDD Anchor
  ENDIF
  CREATE a map or `SysConfig` object with fields to update from `param`.
  CALL database to UPDATE the record WHERE `config_id` matches.
  IF database error:
    RETURN error.
  ENDIF
  RETURN nil.
END FUNCTION
```

**TDD Anchors:**
-   `TestUpdateConfig_Success`: Verify a config is updated with valid data.
-   `TestUpdateConfig_NotFound`: Ensure an error is returned if `ConfigId` does not exist.
-   `TestUpdateConfig_PartialUpdate`: Verify that only specified fields are updated.

---

### `DeleteConfig`

Deletes one or more configuration parameters by their IDs.

**Pseudocode:**

```plaintext
FUNCTION DeleteConfig(configIds: ARRAY<INT>) -> ERROR:
  VALIDATE `configIds` is not empty. // TDD Anchor
  CALL database to DELETE `SysConfig` records WHERE `config_id` is IN `configIds`.
  IF database error:
    RETURN error.
  ENDIF
  RETURN nil.
END FUNCTION
```

**TDD Anchors:**
-   `TestDeleteConfig_Success`: Verify configs are deleted for a list of valid IDs.
-   `TestDeleteConfig_EmptyList`: Verify no operation is performed for an empty ID list.
-   `TestDeleteConfig_InvalidId`: Verify behavior when one or more IDs do not exist.

---

### `GetConfigList`

Retrieves a paginated and filtered list of configurations.

**Pseudocode:**

```plaintext
FUNCTION GetConfigList(param: ConfigListRequest, isPaging: BOOLEAN) -> (ARRAY<ConfigListResponse>, INT):
  INITIALIZE database query on `SysConfig` model.
  ORDER results by `config_id`.

  IF `param.ConfigName` is not empty:
    APPEND `WHERE config_name LIKE ...` to query.
  ENDIF
  IF `param.ConfigKey` is not empty:
    APPEND `WHERE config_key LIKE ...` to query.
  ENDIF
  IF `param.ConfigType` is not empty:
    APPEND `WHERE config_type = ...` to query.
  ENDIF
  IF `param.BeginTime` and `param.EndTime` are not empty:
    APPEND `WHERE create_time BETWEEN ...` to query.
  ENDIF

  INITIALIZE `count` to 0.
  IF `isPaging` is TRUE:
    EXECUTE `COUNT(*)` on the query to get total `count`.
    APPLY `OFFSET` and `LIMIT` to the query based on `param.PageNum` and `param.PageSize`.
  ENDIF

  DECLARE `configs` as an empty array of `ConfigListResponse`.
  EXECUTE `FIND` on the query and populate `configs`.

  RETURN `configs`, `count`.
END FUNCTION
```

**TDD Anchors:**
-   `TestGetConfigList_NoFilters`: Verify it returns all configs.
-   `TestGetConfigList_WithPaging`: Verify `LIMIT` and `OFFSET` work correctly.
-   `TestGetConfigList_WithAllFilters`: Verify combined filtering works as expected.
-   `TestGetConfigList_DateRangeFilter`: Test the `create_time` filter specifically.

---

### `GetConfigByConfigId`

Retrieves a single configuration's details by its ID.

**Pseudocode:**

```plaintext
FUNCTION GetConfigByConfigId(configId: INT) -> ConfigDetailResponse:
  DECLARE `config` as `ConfigDetailResponse`.
  CALL database to FIND the last record from `SysConfig` WHERE `config_id` matches `configId`.
  POPULATE `config` with the result.
  RETURN `config`.
END FUNCTION
```

**TDD Anchors:**
-   `TestGetConfigByConfigId_Success`: Verify it returns the correct config for a valid ID.
-   `TestGetConfigByConfigId_NotFound`: Verify it returns an empty/zeroed struct for a non-existent ID.

---

### `GetConfigByConfigKey`

Retrieves a single configuration's details by its key.

**Pseudocode:**

```plaintext
FUNCTION GetConfigByConfigKey(configKey: STRING) -> ConfigDetailResponse:
  DECLARE `config` as `ConfigDetailResponse`.
  CALL database to FIND the last record from `SysConfig` WHERE `config_key` matches `configKey`.
  POPULATE `config` with the result.
  RETURN `config`.
END FUNCTION
```

**TDD Anchors:**
-   `TestGetConfigByConfigKey_Success`: Verify it returns the correct config for a valid key.
-   `TestGetConfigByConfigKey_NotFound`: Verify it returns an empty/zeroed struct for a non-existent key.

---

### `GetConfigCacheByConfigKey`

Retrieves a configuration by its key, utilizing a cache-aside strategy.

**Pseudocode:**

```plaintext
FUNCTION GetConfigCacheByConfigKey(configKey: STRING) -> ConfigDetailResponse:
  // 1. Attempt to fetch from cache
  GET `configCache` from Redis HASH `sys_config` with key `configKey`.

  IF `configCache` exists AND is not empty:
    TRY to UNMARSHAL `configCache` (JSON string) into a `ConfigDetailResponse` object.
    IF unmarshal is successful:
      RETURN the `config` object. // Cache Hit
    ENDIF
  ENDIF

  // 2. Cache Miss: Fetch from database
  CALL `GetConfigByConfigKey(configKey)` to get `config` from the database.

  // 3. Populate cache if found in DB
  IF `config.ConfigId` is valid (e.g., > 0):
    MARSHAL `config` object into a JSON string.
    SET the JSON string in Redis HASH `sys_config` with key `configKey`.
  ENDIF

  RETURN `config`.
END FUNCTION
```

**TDD Anchors:**
-   `TestGetConfigCache_CacheHit`: Pre-populate cache and verify the function returns data without hitting the DB.
-   `TestGetConfigCache_CacheMiss`: Ensure the cache is empty, verify the function hits the DB, returns data, and populates the cache.
-   `TestGetConfigCache_DBNotFound`: Verify that if a key is not in the DB, it is not added to the cache.
-   `TestGetConfigCache_CorruptedCache`: Test behavior when cache data is malformed and cannot be unmarshalled.
